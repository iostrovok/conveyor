/*
Package workers is an internal package. Package realizes the IManager interface.
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
	"github.com/iostrovok/conveyor/testobject"
)

const (
	stopChLength                = 5
	defaultMetricPeriodInSecond = 10 * time.Second
)

// Manager is an implementation of faces.IManager Interface .
type Manager struct {
	sync.RWMutex

	// global storage for items
	workBench faces.IWorkBench

	minCount, maxCount int
	workerCounter      int
	lengthChannel      int

	isRun  bool
	isLast bool

	activeWorkers *int32

	typ  faces.ManagerType
	name faces.Name

	in, out faces.IChan
	errCh   faces.IChan
	handler faces.GiveBirth
	stopCh  chan struct{}

	ctx     context.Context
	workers []faces.IWorker

	next     faces.IManager
	previous faces.IManager
	wgGlobal *sync.WaitGroup
	wgLocal  *sync.WaitGroup

	metricPeriodDuration time.Duration
	tracer               faces.ITrace

	workersCounter faces.IWorkersCounter

	// need to use in test mode
	testObject faces.ITestObject
}

/*
	WorkerManagerType ManagerType = "worker"
	ErrorManagerType  ManagerType = "error"
	FinalManagerType  ManagerType = "final"
*/

// NewManager is a constructor.
func NewManager(name faces.Name, typ faces.ManagerType, wb faces.IWorkBench, lengthCh, minC, maxC int, tr faces.ITrace) faces.IManager {
	return &Manager{
		typ:                  typ,
		name:                 name,
		lengthChannel:        lengthCh,
		minCount:             minC,
		maxCount:             maxC,
		workers:              make([]faces.IWorker, 0),
		wgLocal:              &sync.WaitGroup{},
		stopCh:               make(chan struct{}, stopChLength),
		metricPeriodDuration: defaultMetricPeriodInSecond,
		tracer:               tr,
		activeWorkers:        new(int32),
		workBench:            wb,
	}
}

// Statistic returns data about manager condition.
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

// Name is a simple getter. It returns Manager name.
func (m *Manager) Name() faces.Name {
	return m.name
}

// SetTestMode is a simple setter. It attaches the testObject.
func (m *Manager) SetTestMode(testObject faces.ITestObject) faces.IManager {
	if testObject == nil {
		testObject = testobject.Empty()
	}

	m.testObject = testObject

	return m
}

// SetWorkersCounter is a simple setter.
func (m *Manager) SetWorkersCounter(wc faces.IWorkersCounter) faces.IManager {
	m.workersCounter = wc

	return m
}

// SetWaitGroup is a simple setter.
func (m *Manager) SetWaitGroup(wg *sync.WaitGroup) faces.IManager {
	wg.Add(1)
	m.wgGlobal = wg

	return m
}

// GetNextManager is a simple getter. It returns next Manager or nil.
func (m *Manager) GetNextManager() faces.IManager {
	return m.next
}

// SetNextManager is a simple setter.
func (m *Manager) SetNextManager(next faces.IManager) faces.IManager {
	m.next = next

	return m
}

// GetPrevManager is a simple getter. It returns previous Manager or nil.
func (m *Manager) GetPrevManager() faces.IManager {
	return m.previous
}

// SetPrevManager is a simple setter.
func (m *Manager) SetPrevManager(previous faces.IManager) faces.IManager {
	m.previous = previous

	return m
}

func (m *Manager) setDataToWorkers() {
	nextManagerName := faces.UnknownName
	if m.next != nil {
		nextManagerName = m.next.Name()
	}

	for _, w := range m.workers {
		w.SetBorderCond(m.typ, m.isLast, nextManagerName)
	}
}

// SetIsLast is a setter. It set up isLast flag to manager and all it's workers.
func (m *Manager) SetIsLast(isLast bool) faces.IManager {
	m.Lock()
	m.isLast = isLast
	m.Unlock()

	m.setDataToWorkers()

	return m
}

// IsLast is a simple getter. It returns flag "is Manager last".
func (m *Manager) IsLast() bool {
	return m.isLast
}

// MetricPeriod is a simple setter.
// It sets up the period between metric evaluations. By default 10 second.
func (m *Manager) MetricPeriod(duration time.Duration) faces.IManager {
	m.metricPeriodDuration = duration

	return m
}

// SetHandler is a simple setter.
func (m *Manager) SetHandler(handler faces.GiveBirth) faces.IManager {
	m.handler = handler

	return m
}

// SetChanIn is a simple setter.
func (m *Manager) SetChanIn(in faces.IChan) faces.IManager {
	m.in = in

	return m
}

// SetChanOut is a simple setter.
func (m *Manager) SetChanOut(out faces.IChan) faces.IManager {
	m.out = out

	return m
}

// SetChanErr is a simple setter.
func (m *Manager) SetChanErr(errCh faces.IChan) faces.IManager {
	m.errCh = errCh

	return m
}

// Stop stops all workers.
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

func (m *Manager) metricPrint(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-m.stopCh:
			return
		case <-time.After(m.metricPeriodDuration):
			if err := m.checkCountWorkers(); err != nil {
				m.logf("[%] error: %s", m.name, err.Error())
			}
		}
	}
}

// is a kind of destructor.
func (m *Manager) waitingGlobalStopAndClose() {
	// if we don't have more workers - close out and stop next manager
	m.wgLocal.Wait()
	defer m.wgGlobal.Done()

	m.logf("[%s] all workers stopped", m.name)

	if m.typ == faces.WorkerManagerType && !m.isLast {
		// it's not last manages - need to close next one
		if m.out != nil {
			m.out.Close()
		}

		return
	}

	if m.typ == faces.WorkerManagerType && m.isLast {
		// it's the last manages - need to close err channel
		if m.errCh != nil {
			m.errCh.Close()
		}

		return
	}

	if m.typ == faces.ErrorManagerType || m.typ == faces.FinalManagerType {
		// close out channel if it's necessary
		if m.out != nil {
			m.out.Close()
		}

		return
	}
}

// Start runs workers and metrics.
func (m *Manager) Start(ctx context.Context) error {
	if m.checkRun() {
		return nil
	}

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

	go m.metricPrint(ctx)
	go m.waitingGlobalStopAndClose()

	return nil
}

func (m *Manager) checkCountWorkers() error {
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
		case nodes.Action_NOTHING: // nothing
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

	m.logf("stopOneWorker: %s", m.workers[0].ID())

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

	//m.logf("addOneWorker: %s", workerName)

	w, err := NewWorker(workerName, m.name, m.workBench, m.in, m.out, m.errCh, m.handler, m.wgLocal, m.tracer, m.activeWorkers)
	if err != nil {
		return err
	}

	// if it's test session
	w.SetTestMode(m.testObject)

	nextManagerName := faces.UnknownName
	if m.next != nil {
		nextManagerName = m.next.Name()
	}

	w.SetBorderCond(m.typ, m.isLast, nextManagerName)
	m.workers = append(m.workers, w)

	return w.Start(m.ctx)
}

func (m *Manager) logf(format string, a ...interface{}) {
	if m.tracer != nil {
		m.tracer.LazyPrintf(string(m.name)+":: "+format, a...)
	}
}
