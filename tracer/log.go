/*
	Package supports the simple realization of ITrace.

	Thanks the "golang.org/x/net/trace" for idea.
*/
package tracer

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/iostrovok/conveyor/faces"
)

type Trace struct {
	sync.RWMutex

	data    []string
	isError bool
}

func NewTrace() faces.ITrace {
	return &Trace{
		data: make([]string, 0),
	}
}

func (t *Trace) LazyPrintf(format string, a ...interface{}) {
	t.Lock()
	t.Unlock()
	t.data = append(t.data, fmt.Sprintf(format, a...))
}

func (t *Trace) SetError() {
	t.isError = true
}

func (t *Trace) Flush() {
	go t.ForceFlush()
}

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
