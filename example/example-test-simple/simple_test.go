package main_test

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	. "github.com/iostrovok/check"
	"github.com/pkg/errors"

	"github.com/iostrovok/conveyor"
	"github.com/iostrovok/conveyor/faces"
	"github.com/iostrovok/conveyor/input"
	"github.com/iostrovok/conveyor/item"
	"github.com/iostrovok/conveyor/testobject"
)

// common section

var (
	CountMySimpleHandlerStop                 = new(int32)
	CountMySimpleHandlerStopStopRockAndRoll  = new(int32)
	CountMySimpleHandlerStart                = new(int32)
	CountMySimpleHandlerStartStopRockAndRoll = new(int32)
	countRunJazz                             = new(int32)
	countRun                                 = new(int32)
)

const (
	total                       = 20
	FirstHandler     faces.Name = "1th-handler"
	SecondHandler    faces.Name = "2th-handler"
	ThirdHandler     faces.Name = "3th-handler"
	FourthHandler    faces.Name = "4th-handler"
	ErrorHandler     faces.Name = "error-user-handler"
	FinalHandlerName faces.Name = "final-user-handler"
)

type MyMessage struct {
	item.Item

	msg string
	id  int
}

func (m *MyMessage) PushedToChannel(label faces.Name) {
	// fmt.Printf("\n\nPushedToChannel %s!!!!!!!!\n\n", string(label)+"channel")
}

func (m *MyMessage) ReceivedFromChannel() {
	// fmt.Printf("\n\nReceivedFromChannel !!!!!!!!\n\n")
}

// >>>>>>>>>>>>>>>>>>>> simple final handler. START

type FinalHandler struct {
	faces.EmptyHandler

	total int
}

func NewFinalHandler(_ faces.Name) (faces.IHandler, error) {
	return &FinalHandler{}, nil
}

func (m *FinalHandler) Start(_ context.Context) error {
	return nil
}

func (m *FinalHandler) Stop(_ context.Context) {
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
	faces.EmptyHandler

	name faces.Name
}

func NewErrHandler(name faces.Name) (faces.IHandler, error) {
	return &ErrHandler{
		name: name,
	}, nil
}

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
	3) Jazz handler may send item direct to fourth handler.
*/
type MySimpleHandler struct {
	faces.EmptyHandler

	counter     int
	sleepSecond int
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

//// does nothing
//func (m *MySimpleHandler) TickerRun(_ context.Context) {
//	fmt.Printf("MySimpleHandler: TickerRun: %s!\n", m.name)
//}
//
//// return 1 second
//func (m *MySimpleHandler) TickerDuration() time.Duration {
//	return time.Second * 1
//}

func (m *MySimpleHandler) Start(_ context.Context) error {
	/* does nothing */

	// just debug message
	fmt.Printf("Start!\n")

	// just debug counter
	atomic.AddInt32(CountMySimpleHandlerStart, 1)

	return nil
}

func (m *MySimpleHandler) StartTestRockAndRoll(_ context.Context, testObject *C) error {
	/* does nothing */

	// just debug message
	fmt.Printf("StartTestRockAndRoll!\n")

	// just debug counter
	atomic.AddInt32(CountMySimpleHandlerStartStopRockAndRoll, 1)

	// Check RockAndRoll case with testObject here
	testObject.Assert(1, Equals, 1)

	return nil
}

func (m *MySimpleHandler) Stop(_ context.Context) {
	/* does nothing */

	// just debug message
	fmt.Printf("Stop!\n")

	// just debug counter
	atomic.AddInt32(CountMySimpleHandlerStop, 1)
}

func (m *MySimpleHandler) StopTestRockAndRoll(_ context.Context, testObject *C) error {
	// just debug message
	fmt.Printf("StopTestRockAndRoll!\n")
	// just debug counter
	atomic.AddInt32(CountMySimpleHandlerStopStopRockAndRoll, 1)

	// Check RockAndRoll case with testObject here
	testObject.Assert(1, Equals, 1)

	return nil
}

func (m *MySimpleHandler) _run(item faces.IItem) error {
	// increase internal counter
	m.counter++

	// prints incoming message and sleeps few second (imitation of real life)
	msg := item.Get().(*MyMessage)
	fmt.Printf("Global ID: %d. Handler data: %s, %d. incoming message: %s\n", item.GetID(), m.name, m.counter, msg.msg)
	time.Sleep(time.Duration(m.sleepSecond) * time.Second)

	// should it do something special-1?
	if msg.id%4 == 0 {
		return errors.New(string(m.name) + " marked the items with error")
	}

	// should it do something special-2?
	if msg.id%5 == 0 && m.name == FirstHandler {
		fmt.Printf("Global ID: %d. setSkipToName: %s\n", item.GetID(), FirstHandler)
		item.SetSkipToName(FirstHandler)
	}

	return nil
}

func (m *MySimpleHandler) Run(item faces.IItem) error {
	atomic.AddInt32(countRun, 1)

	return m._run(item)
}

func (m *MySimpleHandler) RunTestJazz(item faces.IItem, testObject *C) error {
	// just debug message
	fmt.Printf("RunTestJazz!\n")
	// just debug counter
	atomic.AddInt32(countRunJazz, 1)

	// Check RockAndRoll case with testObject here
	testObject.Assert(1, Equals, 1)

	return m._run(item)
}

// <<<<<<<<<<<<<<<<<<<< simple worked handler. END

func buildConveyor(myMaster faces.IConveyor) error {
	// set up our error handler
	if err := myMaster.AddErrorHandler(ErrorHandler, 2, 6, NewErrHandler); err != nil {
		return err
	}

	// set up our final handler
	if err := myMaster.AddFinalHandler(FinalHandlerName, 2, 6, NewFinalHandler); err != nil {
		return err
	}

	// set up 4 simple handlers
	if err := myMaster.AddHandler(FirstHandler, 2, 6, First); err != nil {
		return err
	}

	if err := myMaster.AddHandler(SecondHandler, 2, 6, Second); err != nil {
		return err
	}

	if err := myMaster.AddHandler(ThirdHandler, 2, 6, Third); err != nil {
		return err
	}

	if err := myMaster.AddHandler(FourthHandler, 2, 6, Fourth); err != nil {
		return err
	}

	// start our conveyor
	if err := myMaster.Start(context.Background()); err != nil {
		return err
	}

	return nil
}

// test section.
type testSuite struct{}

var _ = Suite(&testSuite{})

func TestSuite(t *testing.T) { TestingT(t) }

func (s *testSuite) TestSyntax(c *C) {
	/*
		Test object "RockAndRoll"
		Used for Start and Stop methods only
	*/
	to := testobject.New(true, c, "RockAndRoll")

	// create and build new conveyor
	myMaster := conveyor.NewTest(20, faces.ChanStack, "my-app", to)
	c.Assert(buildConveyor(myMaster), IsNil)
	err := myMaster.Start(context.Background())
	c.Assert(err, IsNil)

	for i := 0; i < total; i++ {
		item := input.New().
			Data(&MyMessage{msg: fmt.Sprintf("online item: %d", i), id: i}).
			Priority(100)

		/*
			 object "RockAndRoll"
			Used for Start and Stop methods only
		*/

		res, err := myMaster.RunResTest(item, testobject.New(true, c, "Jazz"))
		fmt.Printf("\nProccesed online: Result: %+v, Error: %+v\n", res, err)
		if res.(*MyMessage).id%4 == 0 || res.(*MyMessage).id%5 == 0 {
			c.Assert(err, NotNil)
		} else {
			c.Assert(err, IsNil)
		}
	}

	myMaster.WaitAndStop()

	c.Assert(atomic.LoadInt32(countRun), Equals, int32(0))
	c.Assert(atomic.LoadInt32(countRunJazz) > 0, Equals, true)

	c.Assert(atomic.LoadInt32(CountMySimpleHandlerStop), Equals, int32(0))
	c.Assert(atomic.LoadInt32(CountMySimpleHandlerStopStopRockAndRoll) > 0, Equals, true)

	c.Assert(atomic.LoadInt32(CountMySimpleHandlerStart), Equals, int32(0))
	c.Assert(atomic.LoadInt32(CountMySimpleHandlerStartStopRockAndRoll) > 0, Equals, true)
}
