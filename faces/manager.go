package faces

import (
	"context"
	"sync"
	"time"

	"github.com/iostrovok/conveyor/protobuf/go/nodes"
)

/*
	....
*/

// ManagerType is a global type for define the manager position inside of conveyor.
type ManagerType string

const (
	// WorkerManagerType is a simple handler.
	WorkerManagerType ManagerType = "worker"
	// ErrorManagerType is a manager which gets all defective parts.
	ErrorManagerType ManagerType = "error"
	// FinalManagerType is a manager in th end of conveyor. "Test of the quality".
	FinalManagerType ManagerType = "final"
)

// Name is a global type for define the name of manager.
type Name string

const (
	// UnknownName is a simple default name.
	UnknownName Name = "unknown"

	// ErrorName is a simple default name for error handler.
	ErrorName Name = "error"
)

/*
IWorkersCounter realizes the interface to rule the numbers of workers for each manage.
It gets the current information about workers and in/out channels for a manager
and returns action for increase or decrease numbers of workers.

See simple realization in workerscounter directory.
*/
type IWorkersCounter interface {
	Check(mc *nodes.ManagerData) (*nodes.ManagerAction, error)
}

/*
IManager is an interface to rule the workers for single handler.
*/
type IManager interface {
	SetHandler(handler GiveBirth) IManager

	Start(ctx context.Context) error
	Stop()

	SetWorkersCounter(wc IWorkersCounter) IManager
	SetChanIn(in IChan) IManager
	SetChanOut(out IChan) IManager
	SetChanErr(errCh IChan) IManager

	GetNextManager() IManager
	SetNextManager(next IManager) IManager

	GetPrevManager() IManager
	SetPrevManager(previous IManager) IManager

	SetIsLast(isLast bool) IManager
	IsLast() bool

	SetWaitGroup(wg *sync.WaitGroup) IManager

	MetricPeriod(duration time.Duration) IManager

	Statistic() *nodes.ManagerData

	Name() Name

	// test mode
	// testObject - object for checking tests
	// startTestSuffix - suffix for start and stop workers methods
	SetTestMode(testObject ITestObject) IManager
}

/*
IWorker is an interface to support using one handler from IManager.
*/
type IWorker interface {
	Start(ctx context.Context) error

	Stop()

	SetTestMode(testObject ITestObject)

	SetBorderCond(typ ManagerType, isLast bool, nextManagerName Name)
	GetBorderCond() (Name, ManagerType, bool)

	Name() Name
	ID() string
}
