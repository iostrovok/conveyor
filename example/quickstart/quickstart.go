package main

import (
	"context"
	"fmt"
	"log"

	"github.com/iostrovok/conveyor"
	"github.com/iostrovok/conveyor/faces"
	"github.com/iostrovok/conveyor/input"
)

// >>>>>>>>>>>>>>>>>>>> simple worker handler START
type MySimpleHandler struct {
	faces.EmptyHandler
	name faces.Name
}

func Handler(name faces.Name) (faces.IHandler, error) {
	return &MySimpleHandler{name: name}, nil
}

func (m *MySimpleHandler) Run(item faces.IItem) error {
	fmt.Printf("MySimpleHandler %s => %d]: %s\n", m.name, item.GetID(), item.Get().(string))
	return nil
}

// <<<<<<<<<<<<<<<<<<<< simple worked handler. END

func main() {
	// create new conveyor
	myMaster := conveyor.New(20, faces.ChanStdGo, "my-app")

	// set up simple handler
	if err := myMaster.AddHandler("handler", 2, 6, Handler); err != nil {
		log.Fatal(err)
	}

	// start our conveyor
	if err := myMaster.Start(context.Background()); err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 100; i++ {
		item := input.New().Data(fmt.Sprintf("item: %d", i+1))
		myMaster.Run(item)
	}

	// wait while conveyor is working
	myMaster.WaitAndStop()
}
