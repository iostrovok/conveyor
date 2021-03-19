/*
Package tracer supports the simple realization of ITrace.
*/
package tracer

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/iostrovok/conveyor/faces"
)

// Trace is an implementation of faces.ITrace Interface .
type Trace struct {
	sync.RWMutex

	data    []string
	isError bool
}

// New is a constructor.
func New() faces.ITrace {
	return &Trace{
		data: make([]string, 0),
	}
}

// LazyPrintf add data for output.
func (t *Trace) LazyPrintf(format string, a ...interface{}) {
	t.Lock()
	defer t.Unlock()

	t.data = append(t.data, fmt.Sprintf(format, a...))
}

// SetError sets that output has prefix "ERROR: ".
func (t *Trace) SetError() {
	t.Lock()
	defer t.Unlock()

	t.isError = true
}

// Flush calls ForceFlush as goroutine.
func (t *Trace) Flush() {
	go t.ForceFlush()
}

// ForceFlush prints collected data.
func (t *Trace) ForceFlush() {
	out := ""
	t.Lock()
	if len(t.data) > 0 {
		out = strings.Join(t.data, " #-# ")
		if t.isError {
			out = "ERROR: " + out
		}
		t.data = make([]string, 0)
	}
	t.Unlock()

	if out != "" {
		log.Print(out)
	}
}
