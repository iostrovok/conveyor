package faces

import (
	"context"
	"time"

	"github.com/iostrovok/conveyor/protobuf/go/nodes"
)

/*
	....
*/

type IConveyor interface {
	Start(ctx context.Context) error
	Stop()
	WaitAndStop()

	// simple pushing
	Run(data interface{})
	RunCtx(ctx context.Context, data interface{})
	RunTrace(ctx context.Context, tr ITrace, data interface{})

	RunPriority(data interface{}, priority int)
	RunPriorityCtx(ctx context.Context, data interface{}, priority int)
	RunPriorityTrace(ctx context.Context, tr ITrace, data interface{}, priority int)

	RunRes(data interface{}, priority int) (interface{}, error)
	RunResCtx(ctx context.Context, data interface{}, priority int) (interface{}, error)
	RunResTrace(ctx context.Context, tr ITrace, data interface{}, priority int) (interface{}, error)

	SetDefaultPriority(defaultPriority int)
	GetDefaultPriority() int

	SetName(name string) IConveyor
	GetName() string

	SetWorkersCounter(wc IWorkersCounter) IConveyor
	AddHandler(manageName Name, minCount, maxCount int, handler GiveBirth) error
	AddErrorHandler(manageName Name, minCount, maxCount int, handler GiveBirth) error
	AddFinalHandler(manageName Name, minCount, maxCount int, handler GiveBirth) error
	Statistic() *nodes.SlaveNodeInfoRequest
	SetTracer(tr ITrace) IConveyor

	// The period between metric evaluations.
	// By default 10 second
	MetricPeriod(duration time.Duration) IConveyor

	// Master node is single node for control the conveyor
	SetMasterNode(addr string, masterNodePeriod time.Duration)

	DefaultPriority() int
}
