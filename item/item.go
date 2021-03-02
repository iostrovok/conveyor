/*
	Package realizes the IItem interface.
*/
package item

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/iostrovok/conveyor/faces"
)

type Item struct {
	sync.RWMutex

	data *Data
}

type Data struct {
	id     int64
	data   interface{}
	ctx    context.Context
	cancel context.CancelFunc
	tracer faces.ITrace
	err    error
	isDone bool

	startTime      time.Time
	localStartTime time.Time

	lastHandler faces.Name
	skipToName  faces.Name
	skipNames   []faces.Name

	handlerNameWithError faces.Name
	priority             int

	// need to use in test mode
	testObject faces.ITestObject
}

func New(ctx context.Context, tr faces.ITrace) faces.IItem {

	item := &Item{}
	item.Init(ctx, tr)
	return item
}

func (i *Item) SetLock() {
	i.Lock()
}

func (i *Item) SetUnlock() {
	i.Unlock()
}

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

func (i *Item) GetTestObject() faces.ITestObject {
	i.RLock()
	defer i.RUnlock()

	return i.data.testObject
}

func (i *Item) SetTestObject(testObject faces.ITestObject) {
	i.RLock()
	defer i.RUnlock()

	i.data.testObject = testObject
}

func (i *Item) CheckData() {
	i.Init(nil, nil)
}

func (i *Item) GetID() int64 {
	i.RLock()
	defer i.RUnlock()

	return i.data.id
}

func (i *Item) SetID(id int64) {
	i.Lock()
	defer i.Unlock()

	i.data.id = id
}

func (i *Item) Get() interface{} {
	i.RLock()
	defer i.RUnlock()

	return i.data.data
}

func (i *Item) Set(data interface{}) {
	i.Lock()
	defer i.Unlock()

	i.data.data = data
}

func (i *Item) AddError(err error) {
	i.Lock()
	defer i.Unlock()

	i.data.err = err
	if err != nil && i.data.tracer != nil {
		i.data.tracer.SetError()
		i.data.tracer.LazyPrintf("%s", err.Error())
	}
}

func (i *Item) GetError() error {
	i.RLock()
	defer i.RUnlock()

	return i.data.err
}

func (i *Item) CleanError() {
	i.Lock()
	defer i.Unlock()

	i.data.err = nil
}

func (i *Item) GetContext() context.Context {
	i.RLock()
	defer i.RUnlock()

	return i.data.ctx
}

// LogTrace pushes message to tracer
func (i *Item) LogTrace(format string, a ...interface{}) {
	i.Lock()
	defer i.Unlock()

	if i.data.tracer != nil {
		i.data.tracer.LazyPrintf(format, a...)
	}
}

// Finish writes tracer for item and flush the tracer
func (i *Item) Finish() {
	i.Lock()
	defer i.Unlock()

	if i.data.tracer != nil {
		i.data.tracer.LazyPrintf(time.Now().Sub(i.data.startTime).String() + " : total")
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

// Cancel emergency breaks the processing by global context
func (i *Item) Cancel() {
	i.Lock()
	defer i.Unlock()

	if i.data.cancel != nil {
		i.data.cancel()
	}
}

// LogTraceFinishTime adds the message and during of period from last call of item.LogTraceFinishTime or item.Start to tracer
func (i *Item) LogTraceFinishTime(format string, a ...interface{}) {
	i.Lock()
	defer i.Unlock()

	if i.data.tracer != nil {
		i.data.tracer.LazyPrintf(time.Now().Sub(i.data.localStartTime).String()+" : "+format, a...)
		i.data.localStartTime = time.Now()
	}
}

// GetPriority returns priority for item
func (i *Item) GetPriority() int {
	i.RLock()
	defer i.RUnlock()

	return i.data.priority
}

// SetPriority sets priority for item if priority queue is used
func (i *Item) SetPriority(priority int) {
	i.Lock()
	defer i.Unlock()

	i.data.priority = priority
}

// GetHandlerError returns the error which got within processing of the item
func (i *Item) GetHandlerError() faces.Name {
	i.RLock()
	defer i.RUnlock()

	return i.data.handlerNameWithError
}

// SetHandlerError sets the error which got within processing of the item
func (i *Item) SetHandlerError(handlerNameWithError faces.Name) {
	i.Lock()
	defer i.Unlock()

	i.data.handlerNameWithError = handlerNameWithError
}

// GetLastHandler returns the last handler name which processed the item
func (i *Item) GetLastHandler() faces.Name {
	i.RLock()
	defer i.RUnlock()

	return i.data.lastHandler
}

// SetHandlerError set the last handler name which processed the item
func (i *Item) SetLastHandler(handlerName faces.Name) {
	i.Lock()
	defer i.Unlock()

	i.data.lastHandler = handlerName
}

// PushedToChannel it should be redefined
func (i *Item) PushedToChannel(_ faces.Name) {
}

// ReceivedFromChannel it should be redefined
func (i *Item) ReceivedFromChannel() {
}

// BeforeProcess it should be redefined
func (i *Item) BeforeProcess(_ faces.Name) {
}

// AfterProcess it should be redefined
func (i *Item) AfterProcess(_ faces.Name, _ error) {
}

// SetSkipToName sets the handler name. Conveyor skips all handlers until that.
// When conveyor reaches that name is set up as EmptySkipName
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

// NeedToSkip checks should be handler skipped or not
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
		return false, errors.New("it is the last handler: skipped name [" +
			string(i.data.skipToName) +
			"] can't be processed")
	}

	return true, nil
}
