package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/iostrovok/conveyor"
	"github.com/iostrovok/conveyor/faces"
	"github.com/iostrovok/conveyor/input"
	"github.com/iostrovok/conveyor/tracer"
)

var total = 20
var totalOnline = 5

const (
	FirstHandler     faces.Name = "1th-handler"
	SecondHandler    faces.Name = "2th-handler"
	ThirdHandler     faces.Name = "3th-handler"
	FourthHandler    faces.Name = "4th-handler"
	ErrorHandler     faces.Name = "error-user-handler"
	FinalHandlerName faces.Name = "final-user-handler"
)

type MyMessage struct {
	msg string
	id  int
}

// >>>>>>>>>>>>>>>>>>>> simple final handler. START

type FinalHandler struct {
	total int
}

func NewFinalHandler(_ faces.Name) (faces.IHandler, error) {
	return &FinalHandler{}, nil
}

func (m *FinalHandler) Start(_ context.Context) error { return nil }
func (m *FinalHandler) Stop() {
	fmt.Printf("final handler: total proceesed %d from %d + 5\n", m.total, total)
}

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

type ErrHandler struct {
	name faces.Name
}

func NewErrHandler(name faces.Name) (faces.IHandler, error) {
	return &ErrHandler{
		name: name,
	}, nil
}

func (m *ErrHandler) Start(_ context.Context) error { return nil }
func (m *ErrHandler) Stop()                         { /* nothing */ }
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
type MySimpleHandler struct {
	counter     int
	sleepSecond int
	id          string
	name        faces.Name
}

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

func (m *MySimpleHandler) Start(ctx context.Context) error {
	return nil
}

func (m *MySimpleHandler) Stop() { /* nothing */ }

func (m *MySimpleHandler) Run(item faces.IItem) error {

	// increase internal counter
	m.counter++

	// prints incoming message and sleeps few second (imitation of real life)
	msg := item.Get().(*MyMessage)
	fmt.Printf("Global ID: %d. Handler data: %s, %d. incoming message: %s\n", item.GetID(), m.name, m.counter, msg.msg)
	time.Sleep(time.Duration(m.sleepSecond) * time.Second)

	// should it do something special-1?
	if item.GetID()%4 == 0 {
		return fmt.Errorf("%s marked the items with error", m.name)
	}

	// should it do something special-2?
	if item.GetID()%5 == 0 && m.name == FirstHandler {
		fmt.Printf("Global ID: %d. setSkipToName: %s\n", item.GetID(), FirstHandler)
		item.SetSkipToName(FirstHandler)
	}

	return nil
}

// <<<<<<<<<<<<<<<<<<<< simple worked handler. END

func main() {
	// create new conveyor
	myMaster := conveyor.New(20, faces.ChanStack, "my-app")

	// set up default tracer and period for collected metric to with tracer
	myMaster.SetTracer(tracer.NewTrace(), time.Second*5).MetricPeriod(3 * time.Second)

	// do we want to remote control? make it here
	//myMaster.SetMasterNode("127.0.0.1:5101", 2*time.Second)

	// optional method to trace process.
	go func() {
		for {
			fmt.Printf("Statistic...........: %+v\n", myMaster.Statistic())
			time.Sleep(10 * time.Second)
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
		log.Fatal(tracer.NewTrace())
	}

	for i := 0; i < total; i++ {
		tr := tracer.NewTrace()
		tr.LazyPrintf("item N %d", i)

		item := input.New().Trace(tr).
			Data(&MyMessage{msg: fmt.Sprintf("item: %d", i)}).
			Priority(i)

		myMaster.Run(item)
	}

	for i := total; i < totalOnline+total; i++ {

		// process one message "online" with reading result

		item := input.New().
			Data(&MyMessage{msg: fmt.Sprintf("online item: %d", 100)}).
			Priority(100)

		res, err := myMaster.RunRes(item)
		fmt.Printf("\nProccesed online: Result: %+v, Error: %+v\n", res, err)
	}

	// wait while conveyor is working
	myMaster.WaitAndStop()
}
