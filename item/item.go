/*
	Package realizes the IItem interface.
*/
package item

import (
	"context"
	"time"

	"github.com/iostrovok/conveyor/faces"
)

type Item struct {
	data *Data
}

type Data struct {
	id     int64
	data   interface{}
	ctx    context.Context
	trace  faces.ITrace
	err    error
	isDone bool

	startTime      time.Time
	localStartTime time.Time

	lastHandler faces.Name
	skipToName  faces.Name

	handlerNameWithError faces.Name
	priority             int
}

func NewItem(ctx context.Context, tr faces.ITrace) faces.IItem {
	item := &Item{}
	item.Init(ctx, tr)
	return item
}

func (i *Item) Init(ctx context.Context, tr faces.ITrace) faces.IItem {
	i.data = &Data{
		data:           make([]interface{}, 0),
		ctx:            ctx,
		skipToName:     faces.EmptySkipName,
		trace:          tr,
		startTime:      time.Now(),
		localStartTime: time.Now(),
	}
	return i
}

func (i *Item) GetID() int64 {
	return i.data.id
}

func (i *Item) SetID(id int64) faces.IItem {
	i.data.id = id
	return i
}

func (i *Item) Get() interface{} {
	return i.data.data
}

func (i *Item) Set(data interface{}) faces.IItem {
	i.data.data = data
	return i
}

func (i *Item) AddError(err error) {
	i.data.err = err
	if err != nil && i.data.trace != nil {
		i.data.trace.SetError()
		i.data.trace.LazyPrintf("%s", err.Error())
	}
}

func (i *Item) GetError() error {
	return i.data.err
}

func (i *Item) CleanError() {
	i.data.err = nil
}

func (i *Item) GetContext() context.Context {
	return i.data.ctx
}

func (i *Item) SetSkipToName(name faces.Name) {
	i.data.skipToName = name
}

func (i *Item) GetSkipToName() faces.Name {
	return i.data.skipToName
}

func (i *Item) CleanSkipToName() {
	i.data.skipToName = faces.EmptySkipName
}

func (i *Item) LogTrace(format string, a ...interface{}) {
	if i.data.trace != nil {
		i.data.trace.LazyPrintf(format, a...)
	}
}

func (i *Item) Finish() {
	i.CleanSkipToName()
	if i.data.trace != nil {
		i.data.trace.LazyPrintf(time.Now().Sub(i.data.startTime).String() + " : total")
		i.data.trace.Flush()
	}
}

func (i *Item) Start() faces.IItem {
	if i.data.trace != nil {
		i.data.startTime = time.Now()
		i.data.localStartTime = time.Now()
	}
	return i
}

func (i *Item) LogTraceFinishTime(format string, a ...interface{}) {
	if i.data.trace != nil {
		i.data.trace.LazyPrintf(time.Now().Sub(i.data.localStartTime).String()+" : "+format, a...)
		i.data.localStartTime = time.Now()
	}
}

func (i *Item) GetPriority() int {
	return i.data.priority
}

func (i *Item) SetPriority(priority int) faces.IItem {
	i.data.priority = priority
	return i
}

func (i *Item) GetHandlerError() faces.Name {
	return i.data.handlerNameWithError
}

func (i *Item) SetHandlerError(handlerNameWithError faces.Name) faces.IItem {
	i.data.handlerNameWithError = handlerNameWithError
	return i
}

func (i *Item) GetLastHandler() faces.Name {
	return i.data.lastHandler
}

func (i *Item) SetLastHandler(handlerName faces.Name) faces.IItem {
	i.data.lastHandler = handlerName
	return i
}
