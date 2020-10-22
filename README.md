# Conveyor

The support of the multi-concurrent design patterns.

Common Conveyor schema:

![alt text](https://github.com/iostrovok/conveyor/blob/master/images/conveyor_schema.jpg?raw=true "Common Conveyor Schema")

[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](docs/conveyor/index.html)


## Installation

conveyor requires a Go version with [Modules](https://github.com/golang/go/wiki/Modules) support and
uses import versioning. So please make sure to initialize a Go module before installing conveyor:

```shell
go mod init github.com/my/repo
go get github.com/iostrovok/conveyor
```

Import:

```go
import (
	"github.com/iostrovok/conveyor"
	"github.com/iostrovok/conveyor/faces"
	"github.com/iostrovok/conveyor/input"
)
```

## Quickstart

See example/quickstart

```go

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/iostrovok/conveyor"
	"github.com/iostrovok/conveyor/faces"
	"github.com/iostrovok/conveyor/input"
)

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
        // process our string-item
		item := input.New().Data(fmt.Sprintf("item: %d", i+1))
		myMaster.Run(item)
	}

	// wait while conveyor is working
	myMaster.WaitAndStop()
}

```
