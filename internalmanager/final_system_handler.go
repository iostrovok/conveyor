/*
Package internalmanager is an internal package. Package provides a handler for support online processing.
*/
package internalmanager

import (
	"context"
	"sync"
	"time"

	"github.com/iostrovok/conveyor/faces"
)

const defaultResultChLen = 2

type oneResult struct {
	ch  chan faces.IItem
	ctx context.Context
}

type myMap struct {
	sync.RWMutex
	data map[int64]*oneResult
}

func (m *myMap) LoadAndDelete(id int64) (*oneResult, bool) {
	m.Lock()
	defer m.Unlock()

	out, find := m.data[id]
	if find {
		delete(m.data, id)
	}

	return out, find
}

func (m *myMap) Store(id int64, res *oneResult) {
	m.Lock()
	defer m.Unlock()

	m.data[id] = res
}

func (m *myMap) Range(f func(key int64, res *oneResult) bool) {
	m.Lock()
	defer m.Unlock()

	for k, v := range m.data {
		if !f(k, v) {
			return
		}
	}
}

// global vars.
var allResults *myMap

// SystemFinalHandler implements the final manager with support the online processing.
type SystemFinalHandler struct {
	faces.EmptyHandler // defines unused methods
}

func init() {
	allResults = &myMap{
		data: map[int64]*oneResult{},
	}
}

// Init returns the SystemFinalHandler init method.
func Init() faces.GiveBirth {
	return func(name faces.Name) (faces.IHandler, error) {
		return &SystemFinalHandler{}, nil
	}
}

// AddID adds new item to waiting of result.
func AddID(ctx context.Context, id int64) chan faces.IItem {
	if ctx == nil {
		ctx = context.Background()
	}

	ch := make(chan faces.IItem, defaultResultChLen)
	res := &oneResult{
		ch:  ch,
		ctx: ctx,
	}

	allResults.Store(id, res)

	return ch
}

// Start is an interface method.
func (m *SystemFinalHandler) Start(_ context.Context) error {
	return nil
}

// Stop is an interface method.
func (m *SystemFinalHandler) Stop(_ context.Context) {
	closeFunc := func(key int64, res *oneResult) bool {
		close(res.ch)
		res.ch = nil
		return true
	}

	allResults.Range(closeFunc)
	allResults = &myMap{
		data: map[int64]*oneResult{},
	}
}

func runOne(res *oneResult, item faces.IItem) {
	select {
	case res.ch <- item: /* nothing */
	case <-res.ctx.Done(): /* nothing */
	case <-time.After(time.Minute): /* nothing */
	}

	// Out of sight, out of mind
	close(res.ch)
	res.ch = nil
}

// Run is an interface method.
// Check the uniq id of item and returns the result if it's necessary.
func (m *SystemFinalHandler) Run(item faces.IItem) error {

	id := item.GetID()
	if value, loaded := allResults.LoadAndDelete(id); loaded {
		go runOne(value, item)
	}

	return nil
}
