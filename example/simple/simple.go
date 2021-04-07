package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"

	"github.com/iostrovok/conveyor"
	"github.com/iostrovok/conveyor/faces"
	"github.com/iostrovok/conveyor/input"
	"github.com/iostrovok/conveyor/item"
	"github.com/iostrovok/conveyor/tracer"
)

var (
	total       = 50
	totalOnline = 10
)

const (
	chanLength              = 20
	defaultPriority         = 100
	defaultStatisticTimeOut = 10 * time.Second
	defaultTracerTimeOut    = 5 * time.Second
	defaultMetricTimeOut    = 3 * time.Second

	// FirstHandler is a name for first handler.
	FirstHandler faces.Name = "1th-handler"
	// SecondHandler is a name for second handler.
	SecondHandler faces.Name = "2th-handler"
	// ThirdHandler is a name for third handler.
	ThirdHandler faces.Name = "3th-handler"
	// FourthHandler is a name for fourth handler.
	FourthHandler faces.Name = "4th-handler"
	// ErrorHandler is a name for error handler.
	ErrorHandler faces.Name = "error-user-handler"
	// FinalHandlerName is a name for final handler.
	FinalHandlerName faces.Name = "final-user-handler"
)

// MyMessage implements the one part on conveyor.
type MyMessage struct {
	item.Item
	msg string
	id  int
}

// >>>>>>>>>>>>>>>>>>>> simple final handler. START

// FinalHandler is a simple final handler.
type FinalHandler struct {
	faces.EmptyHandler
	total int
}

// NewFinalHandler is a constructor.
func NewFinalHandler(_ faces.Name) (faces.IHandler, error) {
	return &FinalHandler{}, nil
}

// Start is interface function.
func (m *FinalHandler) Start(_ context.Context) error { return nil }

// Stop is interface function.
func (m *FinalHandler) Stop(_ context.Context) {
	fmt.Printf("final handler: total proceesed %d from %d + 5\n", m.total, total)
}

// Run is interface function.
func (m *FinalHandler) Run(item faces.IItem) error {
	/*
		It gets each item as after success, as with error.
	*/

	m.total++
	msg := item.Get().(*MyMessage)
	fmt.Printf("FinalHandler. Global ID: %d. incoming message: %s\n", item.GetID(), msg.msg)

	// it makes no sense to return error from final handler - it will not processed.
	return nil
}

// <<<<<<<<<<<<<<<<<<<< simple final handler. END

// >>>>>>>>>>>>>>>>>>>> simple error handler. START

// ErrHandler is a simple error handler.
type ErrHandler struct {
	faces.EmptyHandler
	name faces.Name
}

// NewErrHandler is s constructor.
func NewErrHandler(name faces.Name) (faces.IHandler, error) {
	return &ErrHandler{
		name: name,
	}, nil
}

// Run is interface function.
func (m *ErrHandler) Run(item faces.IItem) error {
	msg := item.Get().(*MyMessage)
	fmt.Printf("ErrHandler => %d]: %d ==> %s\n", item.GetID(), msg.id, msg.msg)

	// if error is returned from error handler it replaces the previous error in item.
	return nil
}

// <<<<<<<<<<<<<<<<<<<< simple error handler. END

// >>>>>>>>>>>>>>>>>>>> simple worker handler START
/*
	Example of simple handler. Actions:

	1) It sleeps few second.
	2) It may setup returns error. After that item is passing to error handler.
	3) First handler may send item direct to fourth handler.
*/

// MySimpleHandler implements of simple handler.
type MySimpleHandler struct {
	faces.EmptyHandler

	counter     int
	sleepSecond int
	name        faces.Name
}

// First is s constructor.
func First(name faces.Name) (faces.IHandler, error) {
	/*
		if we want to setup single for all "first" workers connection to database or grpc, we may do it here.
	*/
	return &MySimpleHandler{
		name:        name,
		counter:     0,
		sleepSecond: 1,
	}, nil
}

// Second is s constructor.
func Second(name faces.Name) (faces.IHandler, error) {
	/*
		if we want to setup single for all "second" workers connection to database or grpc, we may do it here.
	*/
	return &MySimpleHandler{
		name:        name,
		counter:     0,
		sleepSecond: 1,
	}, nil
}

// Third is s constructor.
func Third(name faces.Name) (faces.IHandler, error) {
	/*
		if we want to setup single for all "third" workers connection to database or grpc, we may do it here.
	*/
	return &MySimpleHandler{
		name:        name,
		counter:     0,
		sleepSecond: 1,
	}, nil
}

// Fourth is s constructor.
func Fourth(name faces.Name) (faces.IHandler, error) {
	/*
		if we want to setup single for all "fourth" workers connection to database or grpc, we may do it here.
	*/
	return &MySimpleHandler{
		name:        name,
		counter:     0,
		sleepSecond: 1,
	}, nil
}

// TickerRun does nothing. Just for print info about action.
func (m *MySimpleHandler) TickerRun(_ context.Context) {
	//fmt.Printf("MySimpleHandler: TickerRun: %s!\n", m.name)
}

// TickerDuration return 1 second.
func (m *MySimpleHandler) TickerDuration() time.Duration {
	return time.Second
}

// Start is interface function.
func (m *MySimpleHandler) Start(_ context.Context) error {
	return nil
}

// Stop is interface function.
func (m *MySimpleHandler) Stop(_ context.Context) { /* nothing */ }

// Run is interface function.
func (m *MySimpleHandler) Run(item faces.IItem) error {
	// increase internal counter
	m.counter++

	// prints incoming message and sleeps few second (imitation of real life)
	msg := item.Get().(*MyMessage)
	fmt.Printf("Global ID: %d. Handler data: %s, %d. incoming message: %s\n", item.GetID(), m.name, m.counter, msg.msg)
	time.Sleep(time.Duration(m.sleepSecond) * time.Second)

	// should it do something special-1?
	if item.GetID()%4 == 0 {
		return errors.New(string(m.name) + " marked the items with error")
	}

	// should it do something special-2?
	if item.GetID()%5 == 0 && m.name == FirstHandler {
			fmt.Printf("Global ID: %d. setSkipToName: %s\n", item.GetID(), FourthHandler)
		item.SetSkipToName(FourthHandler)
	}

	return nil
}

// <<<<<<<<<<<<<<<<<<<< simple worked handler. END

func main() {
	// create new conveyor
	myMaster := conveyor.New(chanLength, faces.ChanStack, "my-app")

	// set up default tracer and period for collected metric to with tracer
	myMaster.SetTracer(tracer.New(), defaultTracerTimeOut).MetricPeriod(defaultMetricTimeOut)

	// do we want to remote control? make it here
	// myMaster.SetMasterNode("127.0.0.1:5101", 2*time.Second)

	// optional method to trace process.
	go func() {
		for {
			fmt.Printf("Statistic...........: %+v\n", myMaster.Statistic())
			time.Sleep(defaultStatisticTimeOut)
		}
	}()

	// set up our error handler
	if err := myMaster.AddErrorHandler(ErrorHandler, 2, 6, NewErrHandler); err != nil {
		log.Fatal(err)
	}

	// set up our final handler
	if err := myMaster.AddFinalHandler(FinalHandlerName, 2, 6, NewFinalHandler); err != nil {
		log.Fatal(err)
	}

	// set up 4 simple handlers
	if err := myMaster.AddHandler(FirstHandler, 2, 6, First); err != nil {
		log.Fatal(err)
	}

	if err := myMaster.AddHandler(SecondHandler, 2, 6, Second); err != nil {
		log.Fatal(err)
	}

	if err := myMaster.AddHandler(ThirdHandler, 2, 6, Third); err != nil {
		log.Fatal(err)
	}

	if err := myMaster.AddHandler(FourthHandler, 2, 6, Fourth); err != nil {
		log.Fatal(err)
	}

	// start our conveyor
	if err := myMaster.Start(context.Background()); err != nil {
		log.Fatal(tracer.New())
	}

	for i := 0; i < total; i++ {
		tr := tracer.New()
		tr.LazyPrintf("item N %d", i)

		item := input.New().Trace(tr).
			Data(&MyMessage{msg: fmt.Sprintf("item: %d", i)}).
			Priority(i)

		myMaster.Run(item)
	}


	fmt.Printf("\n\n-------------------------------\n\n")
	fmt.Printf("START ONLINE\n")
	fmt.Printf("\n\n-------------------------------\n\n")


	for i := total; i < totalOnline+total; i++ {
		// process one message "online" with reading result

		item := input.New().
			Data(&MyMessage{msg: fmt.Sprintf("online item: %d", defaultPriority)}).
			Priority(defaultPriority)

		res, err := myMaster.RunRes(item)
		fmt.Printf("\nProccesed online: Result: %+v, Error: %+v\n", res, err)
	}

	// wait while conveyor is working
	myMaster.WaitAndStop()
}
