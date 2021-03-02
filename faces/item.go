package faces

import (
	"context"
)

/*
	....
*/

type IItem interface {
	GetID() int64
	SetID(id int64)

	Get() (data interface{})
	Set(data interface{})
	CheckData()

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

	LogTraceFinishTime(format string, a ...interface{})
	LogTrace(format string, a ...interface{})

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

	// Using for test mode only
	GetTestObject() ITestObject
	SetTestObject(ITestObject)
}
