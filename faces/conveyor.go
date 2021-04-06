/*
Package faces implements the full list of Interfaces.
*/
package faces

import (
	"context"
	"time"

	"github.com/iostrovok/conveyor/protobuf/go/nodes"
)

/*
	....
*/

// IInput is interface for support input data to conveyor.
type IInput interface {
	Context(ctx context.Context) IInput // by default the context.Background()
	Trace(tr ITrace) IInput             // nil is by default
	Data(data interface{}) IInput       // nil is by default
	Priority(priority int) IInput       // by default is IConveyor.DefaultPriority()
	SkipToName(name Name) IInput        // by default is ""

	// return all data above
	Values() (ctx context.Context, tr ITrace, data interface{}, priority *int, name Name)
	Ctx() context.Context
}

// IConveyor is interface for support the conveyor.
type IConveyor interface {
	Start(ctx context.Context) error
	Stop()
	WaitAndStop()

	// simple pushing
	Run(IInput)
	RunRes(IInput) (interface{}, error)

	// simple pushing in test mode
	RunTest(i IInput, object ITestObject)
	RunResTest(i IInput, object ITestObject) (interface{}, error)

	SetDefaultPriority(defaultPriority int)
	GetDefaultPriority() int

	SetName(name string) IConveyor
	GetName() string

	SetWorkersCounter(wc IWorkersCounter) IConveyor
	AddHandler(manageName Name, minCount, maxCount int, handler GiveBirth) error
	AddErrorHandler(manageName Name, minCount, maxCount int, handler GiveBirth) error
	AddFinalHandler(manageName Name, minCount, maxCount int, handler GiveBirth) error
	Statistic() *nodes.SlaveNodeInfoRequest
	SetTracer(tr ITrace, duration time.Duration) IConveyor

	// The period between metric evaluations.
	// By default 10 second
	MetricPeriod(duration time.Duration) IConveyor

	// Master node is single node for control the conveyor
	SetMasterNode(addr string, masterNodePeriod time.Duration)

	DefaultPriority() int

	// Simple getter
	WorkBench() IWorkBench
}
