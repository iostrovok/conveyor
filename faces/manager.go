package faces

import (
	"context"
	"sync"
	"time"

	"github.com/iostrovok/check"

	"github.com/iostrovok/conveyor/protobuf/go/nodes"
)

/*
	....
*/

type ManagerType string

const (
	WorkerManagerType ManagerType = "worker"
	ErrorManagerType  ManagerType = "error"
	FinalManagerType  ManagerType = "final"
)

type Name string

const (
	UnknownName Name = "unknown"
	ErrorName   Name = "error"
)

type IWorkersCounter interface {
	Check(mc *nodes.ManagerData) (*nodes.ManagerAction, error)
}

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
	SetTestMode(mode bool, testObject *check.C) IManager
}

type IWorker interface {
	Start(ctx context.Context) error
	Stop()

	SetTestMode(mode bool, testObject *check.C)

	SetBorderCond(typ ManagerType, isLast bool, nextManagerName Name)
	GetBorderCond() (Name, ManagerType, bool)

	Name() Name
	ID() string
}
