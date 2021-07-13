// Package item implements the faces.IItem interface.
package item

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/iostrovok/conveyor/faces"
	"github.com/iostrovok/conveyor/testobject"
)

const lastHandlerErrorNote = "it is the last handler: skipped name can't be processed"

// Item implements the faces.IItem interface.
type Item struct {
	sync.RWMutex

	data *Data
}

// Data is internal item's storage.
// It exits as property of Item for OOP/inheritance goals.
type Data struct {
	id     int64
	data   interface{}
	ctx    context.Context
	cancel context.CancelFunc
	tracer faces.ITrace
	err    error

	startTime      time.Time
	localStartTime time.Time

	lastHandler faces.Name
	skipToName  faces.Name
	skipNames   []faces.Name
	stopped     bool

	handlerNameWithError faces.Name
	priority             int

	// need to use in test mode
	testObject faces.ITestObject
}

// New is a constructor.
func New(ctx context.Context, tr faces.ITrace) faces.IItem {
	item := &Item{}
	item.Init(ctx, tr)

	return item
}

// InitEmpty makes empty item.
func (i *Item) InitEmpty() {
	i.Init(context.Background(), nil)
}

// Init is full constructor for accurate configuration.
func (i *Item) Init(ctxIn context.Context, tr faces.ITrace) faces.IItem {
	if i.data != nil {
		return i
	}

	i.Lock()
	defer i.Unlock()

	if i.data != nil {
		return i
	}

	// sometime it happens
	if ctxIn == nil {
		ctxIn = context.Background()
	}

	ctx, cancel := context.WithCancel(ctxIn)

	i.data = &Data{
		data:           make([]interface{}, 0),
		ctx:            ctx,
		cancel:         cancel,
		skipToName:     faces.EmptySkipName,
		tracer:         tr,
		startTime:      time.Now(),
		localStartTime: time.Now(),
		skipNames:      make([]faces.Name, 0),
	}

	return i
}

// SetLock is a interface function.
func (i *Item) SetLock() {
	i.Lock()
}

// SetUnlock is a interface function.
func (i *Item) SetUnlock() {
	i.Unlock()
}

// GetTestObject is a interface function. It's a simple getter.
func (i *Item) GetTestObject() faces.ITestObject {
	i.RLock()
	defer i.RUnlock()

	return i.data.testObject
}

// SetTestObject is a interface function. It's a setter.
// It sets up the Empty test object if it gets nil input parameter.
func (i *Item) SetTestObject(testObject faces.ITestObject) {
	i.RLock()
	defer i.RUnlock()

	if testObject == nil {
		testObject = testobject.Empty()
	}

	i.data.testObject = testObject
}

// GetID is a interface function. It's a simple getter.
func (i *Item) GetID() int64 {
	i.RLock()
	defer i.RUnlock()

	return i.data.id
}

// SetID is a simple setter.
func (i *Item) SetID(id int64) {
	i.Lock()
	defer i.Unlock()

	i.data.id = id
}

// Get is a interface function. It's a simple getter.
func (i *Item) Get() interface{} {
	i.RLock()
	defer i.RUnlock()

	return i.data.data
}

// Set is a simple setter. It sets up the data of item.
func (i *Item) Set(data interface{}) {
	i.Lock()
	defer i.Unlock()

	i.data.data = data
}

// AddError fix error in item and push to logger if it's possible.
func (i *Item) AddError(err error) {
	i.Lock()
	defer i.Unlock()

	i.data.err = err
	if err != nil && i.data.tracer != nil {
		i.data.tracer.SetError()
		i.data.tracer.LazyPrintf("%s", err.Error())
	}
}

// GetError is a interface function. It's a simple getter.
func (i *Item) GetError() error {
	i.RLock()
	defer i.RUnlock()

	return i.data.err
}

// CleanError just removes error from item.
func (i *Item) CleanError() {
	i.Lock()
	defer i.Unlock()

	i.data.err = nil
}

// GetContext is a interface function. It's a simple getter.
func (i *Item) GetContext() context.Context {
	i.RLock()
	defer i.RUnlock()

	return i.data.ctx
}

// LogTracef pushes message to tracer.
func (i *Item) LogTracef(format string, a ...interface{}) {
	i.Lock()
	defer i.Unlock()

	if i.data.tracer != nil {
		i.data.tracer.LazyPrintf(format, a...)
	}
}

// Finish writes tracer for item and flush the tracer.
func (i *Item) Finish() {
	i.Lock()
	defer i.Unlock()

	if i.data.tracer != nil {
		i.data.tracer.LazyPrintf(time.Since(i.data.startTime).String() + " : total")
		i.data.tracer.Flush()
	}
}

// Start restarts the timers for measuring the periods for each stage of process and whole process.
func (i *Item) Start() {
	i.Lock()
	defer i.Unlock()

	if i.data.tracer != nil {
		i.data.startTime = time.Now()
		i.data.localStartTime = time.Now()
	}
}

// Cancel emergency breaks the processing by global context.
func (i *Item) Cancel() {
	i.Lock()
	defer i.Unlock()

	if i.data.cancel != nil {
		i.data.cancel()
	}
}

// LogTraceFinishTimef adds the message and during of period from
// last call of item.LogTraceFinishTimef or item.Start to tracer.
func (i *Item) LogTraceFinishTimef(format string, a ...interface{}) {
	i.Lock()
	defer i.Unlock()

	if i.data.tracer != nil {
		i.data.tracer.LazyPrintf(time.Since(i.data.localStartTime).String()+" : "+format, a...)
		i.data.localStartTime = time.Now()
	}
}

// GetPriority returns priority for item.
func (i *Item) GetPriority() int {
	i.RLock()
	defer i.RUnlock()

	return i.data.priority
}

// SetPriority sets priority for item if priority queue is used.
func (i *Item) SetPriority(priority int) {
	i.Lock()
	defer i.Unlock()

	i.data.priority = priority
}

// GetHandlerError returns the error which got within processing of the item.
func (i *Item) GetHandlerError() faces.Name {
	i.RLock()
	defer i.RUnlock()

	return i.data.handlerNameWithError
}

// SetHandlerError sets the error which got within processing of the item.
func (i *Item) SetHandlerError(handlerNameWithError faces.Name) {
	i.Lock()
	defer i.Unlock()

	i.data.handlerNameWithError = handlerNameWithError
}

// GetLastHandler returns the last handler name which processed the item.
func (i *Item) GetLastHandler() faces.Name {
	i.RLock()
	defer i.RUnlock()

	return i.data.lastHandler
}

// Stopped sets up that item should only be processed by the Final or Error Handlers
func (i *Item) Stopped() {
	i.Lock()
	defer i.Unlock()

	i.data.stopped = true
}

// IsStopped indicates that item should only be processed by the Final or Error Handlers
func (i *Item) IsStopped() bool {
	i.RLock()
	defer i.RUnlock()

	return i.data.stopped
}

// SetLastHandler sets the last handler name which processed the item.
func (i *Item) SetLastHandler(handlerName faces.Name) {
	i.Lock()
	defer i.Unlock()

	i.data.lastHandler = handlerName
}

// PushedToChannel does nothing. It should be redefined.
func (i *Item) PushedToChannel(_ faces.Name) {
}

// ReceivedFromChannel does nothing. It should be redefined.
func (i *Item) ReceivedFromChannel() {
}

// BeforeProcess does nothing. It should be redefined.
func (i *Item) BeforeProcess(_ faces.Name) {
}

// AfterProcess does nothing. It should be redefined.
func (i *Item) AfterProcess(_ faces.Name, _ error) {
}

// SetSkipToName sets the handler name. Conveyor skips all handlers until that.
// When conveyor reaches that name is set up as EmptySkipName.
// If conveyor finishes and name is not found item gets error.
func (i *Item) SetSkipToName(name faces.Name) {
	i.Lock()
	defer i.Unlock()

	i.data.skipToName = name
}

// GetSkipToName returns handler name which was set up with SetSkipToName and is not yet processed.
func (i *Item) GetSkipToName() faces.Name {
	i.RLock()
	defer i.RUnlock()

	return i.data.skipToName
}

// SetSkipNames sets the handler names. Conveyor skips all these handlers.
func (i *Item) SetSkipNames(names ...faces.Name) {
	i.Lock()
	defer i.Unlock()

	i.data.skipNames = append(i.data.skipNames, names...)
}

// GetSkipNames returns handler name which was set up with SetSkipNames.
func (i *Item) GetSkipNames() []faces.Name {
	i.RLock()
	defer i.RUnlock()

	return i.data.skipNames
}

// NeedToSkip checks should be handler skipped or not.
func (i *Item) NeedToSkip(worker faces.IWorker) (bool, error) {
	name, typ, isLast := worker.GetBorderCond()

	i.Lock()
	defer i.Unlock()

	// never skip system handlers
	if typ != faces.WorkerManagerType {
		i.data.skipToName = faces.EmptySkipName

		return false, nil
	}

	// check personal skipping
	for _, skip := range i.data.skipNames {
		if skip == name {
			return true, nil
		}
	}

	if i.data.skipToName == faces.EmptySkipName || i.data.skipToName == faces.SkipAll {
		// processing has been not requested
		return false, nil
	}

	if i.data.skipToName == name {
		// Attention, last station. skipName is found! processing...
		i.data.skipToName = faces.EmptySkipName

		return false, nil
	}

	if isLast && typ == faces.WorkerManagerType {
		// no more handlers after that. Fix error.
		return false, errors.WithMessage(errors.New(lastHandlerErrorNote), "skipped name is "+string(i.data.skipToName))
	}

	return true, nil
}
