package faces

import (
	"context"
)

// IItem is interface for support the single part on conveyor.
type IItem interface {
	GetID() int64
	SetID(id int64)

	Get() (data interface{})
	Set(data interface{})
	InitEmpty()

	// processing functions
	GetContext() context.Context
	SetLock()
	SetUnlock()

	AddError(err error)
	GetError() error
	CleanError()

	SetSkipNames(label ...Name)
	SetSkipToName(label Name)
	GetSkipToName() Name
	GetSkipNames() []Name
	NeedToSkip(worker IWorker) (bool, error)

	LogTraceFinishTimef(format string, a ...interface{})
	LogTracef(format string, a ...interface{})

	Start()
	Cancel()
	Finish()

	PushedToChannel(label Name)
	ReceivedFromChannel()
	BeforeProcess(label Name)
	AfterProcess(label Name, err error)

	// >>>>>>> Priority Queue Supports
	GetPriority() int
	SetPriority(priority int)
	// <<<<<<< Priority Queue Support

	SetHandlerError(handlerNameWithError Name)
	GetHandlerError() Name

	SetLastHandler(handlerName Name)
	GetLastHandler() Name

	// Stopped sets up that item should only be processed by the Final or Error Handlers
	Stopped()
	// IsStopped indicates that item should only be processed by the Final or Error Handlers
	IsStopped() bool

	// Using for test mode only
	GetTestObject() ITestObject
	SetTestObject(ITestObject)
}
