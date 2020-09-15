package faces

import (
	"context"
)

/*
	....
*/

type IItem interface {
	GetID() int64
	SetID(id int64) IItem

	Get() (data interface{})
	Set(data interface{}) IItem

	GetContext() context.Context

	AddError(err error)
	GetError() error
	CleanError()

	SetSkipToName(label Name)
	GetSkipToName() Name
	NeedToSkip(worker IWorker) (bool, error)

	LogTraceFinishTime(format string, a ...interface{})
	LogTrace(format string, a ...interface{})

	Start() IItem
	Cancel()
	Finish()

	// >>>>>>> Priority Queue Supports
	GetPriority() int
	SetPriority(priority int) IItem
	// <<<<<<< Priority Queue Support

	SetHandlerError(handlerNameWithError Name) IItem
	GetHandlerError() Name

	SetLastHandler(handlerName Name) IItem
	GetLastHandler() Name
}
