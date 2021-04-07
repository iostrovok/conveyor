/*
Package workbench supports the simple realization of ITrace.
*/
package workbench

import (
	"github.com/pkg/errors"
	"sync"

	"github.com/iostrovok/conveyor/faces"
)

// WorkBench is an implementation of faces.IWorkBench Interface .
type WorkBench struct {
	sync.RWMutex

	chWait chan int

	data         []faces.IItem
	last         int
	activeNumber int
}

// New is a constructor.
func New(len int) faces.IWorkBench {
	out := &WorkBench{
		chWait: make(chan int, len),
		last:   len - 1,
		data:   make([]faces.IItem, len, len),
	}

	// ready for start
	for i := 0; i < len; i++ {
		out.chWait <- i
	}

	return out
}

// Add puts new IItem by number in WorkBench
func (w *WorkBench) Add(item faces.IItem) int {
	i := 0
	ok := false
	select {
	case i, ok = <-w.chWait:
		if !ok {
			return -1
		}

		w.Lock()
		if w.data[i] == nil {
			w.activeNumber++
		}
		w.data[i] = item
		w.Unlock()
	}

	return i
}

// Get returns IItem by number in WorkBench
func (w *WorkBench) Get(i int) (faces.IItem, error) {
	if w.last < i || i < 0 {
		return nil, errors.New("")
	}

	w.RLock()
	defer w.RUnlock()

	return w.data[i], nil
}

// Len returns the total length of WorkBench
func (w *WorkBench) Len() int {
	return w.last + 1
}

// Number returns the number of active IItem in WorkBench
func (w *WorkBench) Count() int {
	return w.activeNumber
}

// Clean removes IItem from WorkBench (makes no-active)
func (w *WorkBench) Clean(i int) {
	if w.last < i || i < 0 {
		return
	}

	w.Lock()
	if w.data[i] != nil {
		w.data[i] = nil
		w.activeNumber--
	}
	w.Unlock()

	// return index to reusing
	w.chWait <- i
}

// GetPriority returns the priority for item by number. If item is not fund, return 0.
func (w *WorkBench) GetPriority(i int) int {
	if w.last < i || i < 0 {
		return 0
	}

	w.RLock()
	defer w.RUnlock()
	if w.data[i] == nil {
		return 0
	}
	return w.data[i].GetPriority()
}
