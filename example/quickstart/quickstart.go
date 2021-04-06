package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/iostrovok/conveyor"
	"github.com/iostrovok/conveyor/faces"
	"github.com/iostrovok/conveyor/input"
)

const chanLength = 20

// >>>>>>>>>>>>>>>>>>>> simple worker handler START.

// MySimpleHandler is a handler.
type MySimpleHandler struct {
	faces.EmptyHandler
	name faces.Name
}

func handler(name faces.Name) (faces.IHandler, error) {
	return &MySimpleHandler{name: name}, nil
}

// Run is interface method.
func (m *MySimpleHandler) Run(item faces.IItem) error {
	fmt.Printf("MySimpleHandler %s => %d]: %s\n", m.name, item.GetID(), item.Get().(string))
	// make you job here >>>>>>>
	time.Sleep(5 * time.Millisecond)
	// make you job here <<<<<<<
	return nil
}

// <<<<<<<<<<<<<<<<<<<< simple worked handler. END.

func main() {
	// create new conveyor
	myMaster := conveyor.New(chanLength, faces.ChanStdGo, "my-app")

	// set up first simple handler
	if err := myMaster.AddHandler("handler", 2, 6, handler); err != nil {
		log.Fatalf("%+v", err)
	}

	// start our conveyor
	if err := myMaster.Start(context.Background()); err != nil {
		log.Fatalf("%+v", err)
	}

	for i := 0; i < 100; i++ {
		item := input.New().Data(fmt.Sprintf("item: %d", i+1))
		myMaster.Run(item)
	}

	// wait while conveyor is working
	myMaster.WaitAndStop()
}
