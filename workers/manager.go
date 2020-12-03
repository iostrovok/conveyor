/*
	Internal package. Package realizes the IManager interface.
*/
package workers

import (
	"context"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"

	"github.com/iostrovok/conveyor/faces"
	"github.com/iostrovok/conveyor/protobuf/go/nodes"
)

/*
	....
*/

type Manager struct {
	sync.RWMutex

	typ           faces.ManagerType
	name          faces.Name
	workerCounter int

	in, out       faces.IChan
	errCh         faces.IChan
	handler       faces.GiveBirth
	stopCh        chan struct{}
	isRun         bool
	lengthChannel int

	ctx                context.Context
	workers            []faces.IWorker
	minCount, maxCount int

	next     faces.IManager
	previous faces.IManager
	isLast   bool
	wgGlobal *sync.WaitGroup
	wgLocal  *sync.WaitGroup

	activeWorkers *int32

	metricPeriodDuration time.Duration
	tracer               faces.ITrace

	workersCounter faces.IWorkersCounter
}

//WorkerManagerType ManagerType = "worker"
//ErrorManagerType  ManagerType = "error"
//FinalManagerType  ManagerType = "final"

func NewManager(name faces.Name, typ faces.ManagerType, lengthChannel, minCount, maxCount int, tr faces.ITrace) faces.IManager {
	return &Manager{
		typ:                  typ,
		name:                 name,
		lengthChannel:        lengthChannel,
		minCount:             minCount,
		maxCount:             maxCount,
		workers:              make([]faces.IWorker, 0),
		wgLocal:              &sync.WaitGroup{},
		stopCh:               make(chan struct{}, 5),
		metricPeriodDuration: 10 * time.Second,
		tracer:               tr,
		activeWorkers:        new(int32),
	}
}

func (m *Manager) Statistic() *nodes.ManagerData {
	m.RLock()
	defer m.RUnlock()

	out := &nodes.ManagerData{
		Name:    string(m.name),
		Created: &timestamp.Timestamp{Seconds: time.Now().Unix()},
		Workers: &nodes.WorkersData{
			Min:    uint32(m.minCount),
			Max:    uint32(m.maxCount),
			Number: uint32(len(m.workers)),
			Active: uint32(atomic.LoadInt32(m.activeWorkers)),
		},
		ChanBefore: []*nodes.ChanData{},
		ChanAfter:  []*nodes.ChanData{},
	}

	if m.in != nil {
		out.ChanBefore = []*nodes.ChanData{m.in.Info()}
	}

	if m.out != nil {
		out.ChanAfter = []*nodes.ChanData{m.out.Info()}
	}

	return out
}

func (m *Manager) Name() faces.Name {
	return m.name
}

func (m *Manager) SetWorkersCounter(wc faces.IWorkersCounter) faces.IManager {
	m.workersCounter = wc
	return m
}

func (m *Manager) SetWaitGroup(wg *sync.WaitGroup) faces.IManager {
	wg.Add(1)
	m.wgGlobal = wg
	return m
}
func (m *Manager) GetNextManager() faces.IManager {
	return m.next
}

func (m *Manager) SetNextManager(next faces.IManager) faces.IManager {
	m.next = next
	return m
}

func (m *Manager) GetPrevManager() faces.IManager {
	return m.previous
}

func (m *Manager) SetPrevManager(previous faces.IManager) faces.IManager {
	m.previous = previous
	return m
}

func (m *Manager) setDataToWorkers() {
	for _, w := range m.workers {
		w.SetBorderCond(m.typ, m.isLast)
	}
}

func (m *Manager) SetIsLast(isLast bool) faces.IManager {
	m.Lock()
	m.isLast = isLast
	m.Unlock()

	m.setDataToWorkers()

	return m
}

func (m *Manager) IsLast() bool {
	return m.isLast
}

// The period between metric evaluations.
// By default 10 second
func (m *Manager) MetricPeriod(duration time.Duration) faces.IManager {
	m.metricPeriodDuration = duration
	return m
}

func (m *Manager) SetHandler(handler faces.GiveBirth) faces.IManager {
	m.handler = handler
	return m
}

func (m *Manager) SetChanIn(in faces.IChan) faces.IManager {
	m.in = in
	return m
}

func (m *Manager) SetChanOut(out faces.IChan) faces.IManager {
	m.out = out
	return m
}

func (m *Manager) SetChanErr(errCh faces.IChan) faces.IManager {
	m.errCh = errCh
	return m
}

func (m *Manager) Stop() {

	m.Lock()
	defer m.Unlock()

	m.isRun = false
	m.stopCh <- struct{}{}
	for _, w := range m.workers {
		w.Stop()
	}

	m.workers = make([]faces.IWorker, 0)
}

func (m *Manager) checkRun(checks ...bool) bool {
	m.Lock()
	defer m.Unlock()

	checks = append(checks, true)

	if m.isRun == checks[0] {
		return true
	}
	if checks[0] {
		m.isRun = true
	}

	return false
}

func (m *Manager) Start(ctx context.Context) error {
	if m.checkRun() {
		return nil
	}

	m.logTrace("[%s] is started", m.name)

	m.ctx = ctx

	if len(m.workers) >= m.minCount {
		return nil
	}

	// init workers
	for i := len(m.workers); i < m.minCount; i++ {
		if err := m.addOneWorker(); err != nil {
			return err
		}
	}

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case <-m.stopCh:
				return
			case <-time.After(m.metricPeriodDuration):
				if err := m.checkCountWorkers(); err != nil {
					m.logTrace("[%] error: %s", m.name, err.Error())
				}
			}
		}
	}(ctx)

	go func() {
		// if we don't have more workers - close out and stop next manager
		m.wgLocal.Wait()
		defer m.wgGlobal.Done()

		m.logTrace("[%s] all workers stopped", m.name)

		if m.typ == faces.WorkerManagerType && !m.isLast {
			if m.out != nil {
				m.out.Close()
			}
			return
		}

		if m.typ == faces.WorkerManagerType && m.isLast {
			if m.errCh != nil {
				m.errCh.Close()
			}
			return
		}

		if m.typ == faces.ErrorManagerType || m.typ == faces.FinalManagerType {
			if m.out != nil {
				m.out.Close()
			}
			return
		}
	}()

	return nil
}

func (m *Manager) checkCountWorkers() error {
	if m.checkRun(false) || !m.in.IsActive() {
		return nil
	}

	workersWantTo, err := m.workersCounter.Check(m.Statistic())
	if err != nil {
		return err
	}

	for i := 0; i < int(workersWantTo.Delta); i++ {
		switch workersWantTo.Action {
		case nodes.Action_UP:
			if err := m.addOneWorker(); err != nil {
				return err
			}
		case nodes.Action_DOWN:
			m.stopOneWorker()
		}
	}

	return nil
}

func (m *Manager) stopOneWorker() {
	if m.checkRun(false) || !m.in.IsActive() {
		return
	}

	if len(m.workers) == 0 {
		return
	}

	m.Lock()
	defer m.Unlock()

	m.logTrace("stopOneWorker: %s", m.workers[0].ID())

	// always stop first. ¯\_(ツ)_/¯ WHY.
	m.workers[0].Stop()
	m.workers = m.workers[1:]
}

func (m *Manager) addOneWorker() error {
	if m.checkRun(false) || !m.in.IsActive() {
		return nil
	}

	m.Lock()
	defer m.Unlock()

	m.workerCounter++
	workerName := string(m.name) + "-" + strconv.Itoa(m.workerCounter)

	m.logTrace("addOneWorker: %s", workerName)

	w, err := NewWorker(workerName, m.name, m.in, m.out, m.errCh, m.handler, m.wgLocal, m.tracer, m.activeWorkers)
	if err != nil {
		return err
	}

	w.SetBorderCond(m.typ, m.isLast)
	m.workers = append(m.workers, w)
	return w.Start(m.ctx)
}

func (m *Manager) logTrace(format string, a ...interface{}) {
	if m.tracer != nil {
		m.tracer.LazyPrintf(string(m.name)+":: "+format, a...)
	}
}
