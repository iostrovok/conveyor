/*
Package conveyor implements a conveyor for process items step by step.
*/
package conveyor

import (
	"context"
	"log"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/iostrovok/conveyor/faces"
	"github.com/iostrovok/conveyor/internalmanager"
	"github.com/iostrovok/conveyor/item"
	"github.com/iostrovok/conveyor/protobuf/go/nodes"
	"github.com/iostrovok/conveyor/queues"
	"github.com/iostrovok/conveyor/slavenode"
	"github.com/iostrovok/conveyor/testobject"
	"github.com/iostrovok/conveyor/workbench"
	"github.com/iostrovok/conveyor/workers"
	"github.com/iostrovok/conveyor/workerscounter"
)

const (
	defaultFinalName faces.Name = "final-system-handler"
	defaultErrorName faces.Name = "error-empty-handler"
	defaultPriority             = 0
)

// DefaultMetricPeriodDuration set up by default default for send/print metric.
const DefaultMetricPeriodDuration = 10 * time.Second

// Conveyor is main top structure supported the faces.IConveyor interface.
// It is exported to possibility of debug in outside applications.
type Conveyor struct {
	data *data
}

type data struct {
	sync.RWMutex
	workerGroup *sync.WaitGroup
	errorGroup  *sync.WaitGroup
	finalGroup  *sync.WaitGroup

	clusterID string
	name      string
	itemID    *int64
	isRun     bool

	// storage for items
	workBench faces.IWorkBench

	workBranchLength int
	chanType         faces.ChanType
	errCh            faces.IChan
	inCh             faces.IChan
	outCh            faces.IChan

	managerCounter     int
	userFinalManager   faces.IManager // outside final manager
	systemFinalManager faces.IManager
	firstWorkerManager faces.IManager
	lastWorkerManager  faces.IManager
	firstErrorManager  faces.IManager
	lastErrorManager   faces.IManager

	metricPeriodDuration time.Duration
	workersCounter       faces.IWorkersCounter

	tracer               faces.ITrace
	tracerPeriodDuration time.Duration

	stopContext   context.Context
	cancelContext context.CancelFunc

	masterNodeAddress string
	masterNodePeriod  time.Duration
	slaveNode         *slavenode.SlaveNode

	defaultPriority int
	uniqNames       []faces.Name

	// need to use in test mode
	testObject faces.ITestObject
}

// New is a simple constructor.
func New(workBranchLength int, chanType faces.ChanType, name string) faces.IConveyor {
	c := &Conveyor{}

	return c.Init(workBranchLength, chanType, name, testobject.Empty())
}

// NewTest is full constructor for accurate configuration in test mode.
func NewTest(workBranchLength int, chanType faces.ChanType, name string, testObject faces.ITestObject) faces.IConveyor {
	c := &Conveyor{}

	return c.Init(workBranchLength, chanType, name, testObject)
}

// Init is full constructor for accurate configuration.
func (c *Conveyor) Init(workBranchLength int, chType faces.ChanType, name string, testObject faces.ITestObject) faces.IConveyor {
	if name == "" {
		name = uuid.New().String()
	}

	if workBranchLength < 1 {
		workBranchLength = 1
	}

	c.data = &data{
		workBench:            workbench.New(workBranchLength),
		clusterID:            name + "-" + strconv.FormatInt(time.Now().Unix(), 10),
		name:                 name,
		workBranchLength:     workBranchLength,
		chanType:             chType,
		workerGroup:          &sync.WaitGroup{},
		errorGroup:           &sync.WaitGroup{},
		finalGroup:           &sync.WaitGroup{},
		metricPeriodDuration: DefaultMetricPeriodDuration,

		itemID: new(int64),

		// default IWorkersCounter
		workersCounter: workerscounter.New(),

		uniqNames:       []faces.Name{},
		defaultPriority: defaultPriority,
		testObject:      testObject,
	}

	c.data.outCh = queues.New(c.data.workBench, chType)
	c.data.inCh = queues.New(c.data.workBench, chType)
	c.data.errCh = queues.New(c.data.workBench, chType)

	c.addSystemFinalHandler()

	return c
}

// getItemFrommInput excavator IItem or create new.
func (c *Conveyor) getItemFrommInput(i faces.IInput) faces.IItem {
	ctx, tr, val, priorityRef, skipToName := i.Values()

	if ctx == nil {
		ctx = context.Background()
	}

	priority := c.data.defaultPriority
	if priorityRef != nil {
		priority = *priorityRef
	}

	var it faces.IItem
	if v, ok := val.(faces.IItem); ok {
		it = v
		it.InitEmpty()
	} else {
		it = item.New(ctx, tr)
	}

	it.SetPriority(priority)
	it.Set(val)
	it.SetID(atomic.AddInt64(c.data.itemID, 1))
	if skipToName != "" {
		it.SetSkipToName(skipToName)
	}

	return it
}

// RunTest creates the new item over interface and sends to conveyor.
// If priority queue is used the default priority will be set up.
func (c *Conveyor) RunTest(i faces.IInput, testObject faces.ITestObject) {
	it := c.getItemFrommInput(i)

	// set test suffix
	it.SetTestObject(testObject)

	// marker before pushing to first channel
	it.PushedToChannel(c.data.firstWorkerManager.Name())
	it.Start()
	c.data.inCh.ChanIn() <- c.WorkBench().Add(it)
}

// Run creates the new item over interface and sends to conveyor.
// If priority queue is used the default priority will be set up.
func (c *Conveyor) Run(i faces.IInput) {
	it := c.getItemFrommInput(i)
	// marker before pushing to first channel
	it.PushedToChannel(c.data.firstWorkerManager.Name())
	it.Start()
	c.data.inCh.ChanIn() <- c.WorkBench().Add(it)
}

// RunResTest creates the new item over interface, sends to conveyor and returns result.
func (c *Conveyor) RunResTest(i faces.IInput, testObject faces.ITestObject) (interface{}, error) {
	it := c.getItemFrommInput(i)

	// set test suffix
	it.SetTestObject(testObject)

	return c._runRes(it)
}

// RunRes creates the new item over interface, sends to conveyor and returns result.
func (c *Conveyor) RunRes(i faces.IInput) (interface{}, error) {
	it := c.getItemFrommInput(i)

	return c._runRes(it)
}

func (c *Conveyor) _runRes(it faces.IItem) (interface{}, error) {
	ctx := it.GetContext()

	// it adds id to the latest system handler which will wait for result, get it and return to channel.
	ch := internalmanager.AddID(ctx, it.GetID())

	// marker before pushing to first channel
	it.PushedToChannel(c.data.firstWorkerManager.Name())
	it.Start()
	c.data.inCh.ChanIn() <- c.WorkBench().Add(it)

	select {
	case <-ctx.Done():
		return nil, errors.New("context is canceled in RunRes")
	case it, ok := <-ch:
		if !ok || it == nil {
			return nil, errors.New("unexpected error: result channel is closed")
		}

		return it.Get(), it.GetError()
	}
}

// SetDefaultPriority sets the priority of items.
// It makes sense if priority queue is used.
// If priority is not set up it equals the 0 (defaultPriority constant).
func (c *Conveyor) SetDefaultPriority(defaultPriority int) {
	c.data.defaultPriority = defaultPriority
}

// DefaultPriority returns the current default priority.
// The current default priority is used if a priority is not set up for item.
func (c *Conveyor) DefaultPriority() int {
	return c.data.defaultPriority
}

// GetDefaultPriority returns the default priority of items.
// It makes sense if priority queue is used.
func (c *Conveyor) GetDefaultPriority() int {
	return c.data.defaultPriority
}

// SetMasterNode sets the internet address  master node.
// Master node allow to get information online about current conveyor.
// see more information github.com/iostrovok/conveyormaster.
func (c *Conveyor) SetMasterNode(addr string, masterNodePeriod time.Duration) {
	c.data.masterNodeAddress = addr
	c.data.masterNodePeriod = masterNodePeriod
	if c.data.masterNodePeriod == 0 {
		c.data.masterNodePeriod = time.Minute
	}
}

// StartConnectMasterNode tries to connect to mater node with timeout.
func (c *Conveyor) sendToMasterNode() {
	if c.data.slaveNode == nil {
		return
	}

	// firstWorkerManager
	if _, err := c.data.slaveNode.Send(context.Background(), c.Statistic()); err != nil {
		log.Printf("slaveNode.Send.err: %s\n", err.Error())
	}

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(c.data.masterNodePeriod):
				// TODO: result will use to update configuration online
				_, err := c.data.slaveNode.Send(ctx, c.Statistic())
				if err != nil {
					log.Printf("slaveNode.Send.err: %s\n", err.Error())
				}
			}
		}
	}(c.data.stopContext)
}

func (c *Conveyor) initMasterNode() error {
	sn, err := slavenode.New(c.data.masterNodeAddress)
	if err == nil {
		c.logTracef("success connection to master node")
		c.data.Lock()
		c.data.slaveNode = sn
		c.data.Unlock()
		c.sendToMasterNode()
	} else {
		c.logTracef("error connection to master node: %s", err.Error())
	}

	return err
}

func (c *Conveyor) runMasterNode() {
	if c.data.masterNodeAddress == "" {
		return
	}

	go func(ctx context.Context) {
		if err := c.initMasterNode(); err == nil {
			return
		}

		for {
			select {
			case <-time.After(c.data.masterNodePeriod):
				if err := c.initMasterNode(); err == nil {
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}(c.data.stopContext)
}

// SetWorkersCounter sets up the tracer with IWorkersCounter interface.
// WorkersCounter rules the number of current worked handlers.
func (c *Conveyor) SetWorkersCounter(wc faces.IWorkersCounter) faces.IConveyor {
	c.data.workersCounter = wc

	return c
}

// SetName is a simple setter for name property.
func (c *Conveyor) SetName(name string) faces.IConveyor {
	c.data.name = name

	return c
}

// GetName is a simple getter for name property.
func (c *Conveyor) GetName() string {
	return c.data.name
}

func (c *Conveyor) flushTrace() {
	if c.data.tracer != nil {
		c.data.tracer.ForceFlush()
	}
}

func (c *Conveyor) logTracef(format string, a ...interface{}) {
	if c.data.tracer != nil {
		c.data.tracer.LazyPrintf(format, a...)
	}
}

// SetTracer sets up the tracer with ITrace interface.
func (c *Conveyor) SetTracer(tr faces.ITrace, duration time.Duration) faces.IConveyor {
	c.data.Lock()
	defer c.data.Unlock()

	if tr != nil {
		c.data.tracer = tr
		c.data.tracerPeriodDuration = duration
	}

	return c
}

// MetricPeriod sets up the period between metric evaluations.
// By default 10 second.
func (c *Conveyor) MetricPeriod(duration time.Duration) faces.IConveyor {
	c.data.metricPeriodDuration = duration

	c.data.Lock()
	defer c.data.Unlock()

	mg := c.data.firstWorkerManager
	for {
		if mg == nil {
			break
		}
		mg = mg.MetricPeriod(duration).GetNextManager()
	}

	return c
}

func (c *Conveyor) startGroup(manager faces.IManager) error {
	for {
		if manager == nil {
			break
		}
		if err := manager.Start(c.data.stopContext); err != nil {
			return err
		}
		manager = manager.GetNextManager()
	}

	return nil
}

// Start starts the conveyor.
func (c *Conveyor) Start(ctx context.Context) error {
	if c.data.isRun {
		return nil
	}

	// adds default error manager if it's necessary
	if c.data.firstErrorManager == nil {
		if err := c.AddErrorHandler(defaultErrorName, 1, 2, faces.MakeEmptyHandler); err != nil {
			return err
		}
	}

	c.data.Lock()
	defer c.data.Unlock()

	// fix error handler channels
	c.data.firstErrorManager.SetChanIn(c.data.errCh)
	c.data.lastErrorManager.SetIsLast(true).SetChanErr(c.data.outCh).SetChanOut(c.data.outCh)

	// marks the lastWorkerManager handler. It's lastWorkerManager handler manager, not final.
	c.data.firstWorkerManager.SetChanIn(c.data.inCh)
	c.data.lastWorkerManager.SetIsLast(true).SetChanOut(c.data.outCh)

	// adds default final manager
	c.data.isRun = true
	c.data.stopContext, c.data.cancelContext = context.WithCancel(ctx)

	// start all groups
	for _, first := range []faces.IManager{
		c.data.systemFinalManager, c.data.firstErrorManager, c.data.firstWorkerManager,
	} {
		if err := c.startGroup(first); err != nil {
			return err
		}
	}

	// sending information about cluster to
	c.runMasterNode()

	if c.data.tracer != nil {
		go c.tracerFlush(c.data.stopContext)
	}

	return nil
}

func (c *Conveyor) tracerFlush(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(c.data.tracerPeriodDuration):
			if c.data.tracer != nil {
				c.data.tracer.Flush()
			}
		}
	}
}

// Stop stops the conveyor.
// Processing of items will be  interrupted.
func (c *Conveyor) Stop() {
	c.data.Lock()
	defer c.data.Unlock()

	if !c.data.isRun {
		return
	}

	mg := c.data.firstWorkerManager
	for {
		if mg == nil {
			return
		}
		mg.Stop()
		mg = mg.GetNextManager()
	}
}

// WaitAndStop waits while all handler are finished and exits.
// Processing of items will not be interrupted.
func (c *Conveyor) WaitAndStop() {
	if !c.data.isRun {
		return
	}

	// close income channel and wait for all managers are stopped.
	c.data.inCh.Close()
	c.data.workerGroup.Wait()

	// waiting error handlers
	c.data.errorGroup.Wait()

	// waiting error handlers
	c.data.finalGroup.Wait()

	// lastWorkerManager actions
	if c.data.slaveNode != nil {
		// ignore all errors
		_, _ = c.data.slaveNode.Send(context.Background(), c.Statistic())
	}

	c.data.cancelContext()
	c.flushTrace()
}

func (c *Conveyor) checkUniqName(manageName faces.Name) error {
	for _, n := range c.data.uniqNames {
		if n == manageName {
			return errors.WithStack(errors.New("not uniq handler name '" + string(manageName) + "'"))
		}
	}

	c.data.uniqNames = append(c.data.uniqNames, manageName)

	return nil
}

// AddFinalHandler adds customer final handler as the latest handler from all.
// Only single customer final handler is allow.
// If custom final handler returned error the conveyor doesn't process and log it.
func (c *Conveyor) AddFinalHandler(name faces.Name, minCount, maxCount int, handler faces.GiveBirth) error {
	c.data.Lock()
	defer c.data.Unlock()

	if err := c.checkUniqName(name); err != nil {
		return err
	}

	c.data.managerCounter++
	if name == "" {
		name = faces.Name(c.data.name + "-" + strconv.Itoa(c.data.managerCounter))
	}

	c.data.userFinalManager = workers.NewManager(
		name,
		faces.FinalManagerType,
		c.data.workBench,
		c.data.workBranchLength,
		minCount,
		maxCount,
		c.data.tracer,
	).
		SetHandler(handler).
		SetIsLast(true).
		SetWaitGroup(c.data.finalGroup).
		MetricPeriod(c.data.metricPeriodDuration).
		SetWorkersCounter(c.data.workersCounter).
		SetTestMode(c.data.testObject)

	in := queues.New(c.data.workBench, c.data.chanType)
	c.data.systemFinalManager.SetNextManager(c.data.userFinalManager).SetChanOut(in).SetChanErr(in)
	c.data.userFinalManager.SetPrevManager(c.data.systemFinalManager).SetChanIn(in)

	return nil
}

// addSystemFinalHandler adds default final manager for processing with getting result.
func (c *Conveyor) addSystemFinalHandler() {
	c.data.Lock()
	defer c.data.Unlock()

	c.logTracef("AddOwnFinalHandler")

	c.data.managerCounter++
	c.data.uniqNames = append(c.data.uniqNames, defaultFinalName)

	handler := internalmanager.Init()

	c.data.systemFinalManager = workers.NewManager(
		defaultFinalName,
		faces.FinalManagerType,
		c.data.workBench,
		c.data.workBranchLength,
		1, 2, c.data.tracer,
	).
		SetHandler(handler).
		SetIsLast(false).
		SetWaitGroup(c.data.finalGroup).
		MetricPeriod(c.data.metricPeriodDuration).
		SetWorkersCounter(c.data.workersCounter).
		SetChanIn(c.data.outCh).
		SetTestMode(c.data.testObject)
}

// AddHandler adds customer handler.
// Parameter name should be unique.
// minCount should be less or equal the maxCount and great than zero.
// The order of adding handlers are important. The handlers are called in the same order as they were added.
func (c *Conveyor) AddHandler(name faces.Name, minCount, maxCount int, handler faces.GiveBirth) error {
	c.data.Lock()
	defer c.data.Unlock()

	//c.logTracef("AddHandler %s", name)
	if name == "" {
		return errors.New("handler name can not be empty")
	}

	if err := c.checkUniqName(name); err != nil {
		return err
	}

	c.data.managerCounter++

	next := workers.NewManager(
		name, faces.WorkerManagerType, c.data.workBench,
		c.data.workBranchLength, minCount, maxCount, c.data.tracer,
	).
		SetHandler(handler).
		SetChanErr(c.data.errCh).
		SetWaitGroup(c.data.workerGroup).
		MetricPeriod(c.data.metricPeriodDuration).
		SetWorkersCounter(c.data.workersCounter).
		SetTestMode(c.data.testObject)

	if c.data.lastWorkerManager != nil {
		in := queues.New(c.data.workBench, c.data.chanType)
		c.data.lastWorkerManager.SetNextManager(next).SetChanOut(in)
		next.SetPrevManager(c.data.lastWorkerManager).SetChanIn(in)
	}

	c.data.lastWorkerManager = next
	if c.data.firstWorkerManager == nil {
		c.data.firstWorkerManager = next
	}

	return nil
}

// AddErrorHandler adds custom error handler for processing the errors which were returned with work handler.
// Multiple custom error handlers are allowed.
// If custom error handler returned error the conveyor logs the error but doesn't process.
func (c *Conveyor) AddErrorHandler(manageName faces.Name, minCount, maxCount int, handler faces.GiveBirth) error {
	c.data.Lock()
	defer c.data.Unlock()

	//c.logTracef("AddErrorHandler %s", manageName)
	if manageName == "" {
		return errors.New("error handler name can not be empty")
	}

	if err := c.checkUniqName(manageName); err != nil {
		return err
	}

	c.data.managerCounter++

	next := workers.NewManager(
		manageName,
		faces.ErrorManagerType,
		c.data.workBench,
		c.data.workBranchLength,
		minCount, maxCount,
		c.data.tracer,
	).
		SetHandler(handler).
		SetWaitGroup(c.data.errorGroup).
		MetricPeriod(c.data.metricPeriodDuration).
		SetWorkersCounter(c.data.workersCounter).
		SetTestMode(c.data.testObject)

	if c.data.lastErrorManager != nil {
		in := queues.New(c.data.workBench, c.data.chanType)
		c.data.lastErrorManager.SetNextManager(next).SetChanOut(in).SetChanErr(in)
		next.SetPrevManager(c.data.lastErrorManager).SetChanIn(in)
	}

	c.data.lastErrorManager = next
	if c.data.firstErrorManager == nil {
		c.data.firstErrorManager = next
	}

	return nil
}

func managerStatistic(managers ...faces.IManager) []*nodes.ManagerData {
	out := make([]*nodes.ManagerData, 0)
	for _, manager := range managers {
		if manager != nil {
			out = append(out, manager.Statistic())
		}
	}

	return out
}

// Statistic returns the information about current stage of conveyor.
func (c *Conveyor) Statistic() *nodes.SlaveNodeInfoRequest {
	c.data.RLock()
	defer c.data.RUnlock()

	out := &nodes.SlaveNodeInfoRequest{
		ClusterID:        c.data.clusterID,
		NodeID:           c.data.name,
		ErrorManagerData: []*nodes.ManagerData{},
		FinalManagerData: managerStatistic(
			c.data.systemFinalManager,
			c.data.userFinalManager,
		),
		ManagerData: []*nodes.ManagerData{},
	}

	if !c.data.isRun {
		return out
	}

	mg := c.data.firstWorkerManager
	for {
		if mg == nil {
			break
		}
		out.ManagerData = append(out.ManagerData, mg.Statistic())
		mg = mg.GetNextManager()
	}

	mge := c.data.firstErrorManager
	for {
		if mge == nil {
			break
		}
		out.ErrorManagerData = append(out.ErrorManagerData, mge.Statistic())
		mge = mge.GetNextManager()
	}

	return out
}

// WorkBench is a simple getter
func (c *Conveyor) WorkBench() faces.IWorkBench {
	return c.data.workBench
}
