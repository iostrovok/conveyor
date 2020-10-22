/*
	Internal package. Package provides a handler for support online processing.
*/
package internalmanager

import (
	"context"
	"sync"

	"github.com/iostrovok/conveyor/faces"
)

type oneResult struct {
	ch  chan faces.IItem
	ctx context.Context
}

type Map struct {
	sync.RWMutex
	data map[int64]*oneResult
}

func (m *Map) LoadAndDelete(id int64) (*oneResult, bool) {
	m.Lock()
	defer m.Unlock()
	out, find := m.data[id]
	if find {
		delete(m.data, id)
	}
	return out, find
}

func (m *Map) Store(id int64, res *oneResult) {
	m.Lock()
	defer m.Unlock()
	m.data[id] = res
}

func (m *Map) Range(f func(key int64, res *oneResult) bool) {
	m.Lock()
	defer m.Unlock()
	for k, v := range m.data {
		if !f(k, v) {
			return
		}
	}
}

// global vars
var allResults *Map
var mx *sync.RWMutex

type SystemFinalHandler struct{}

func init() {
	mx = new(sync.RWMutex)
	allResults = &Map{
		data: map[int64]*oneResult{},
	}
}

func Init() faces.GiveBirth {
	return func(name faces.Name) (faces.IHandler, error) {
		return &SystemFinalHandler{}, nil
	}
}

func AddId(id int64, ctx context.Context) chan faces.IItem {
	mx.Lock()
	defer mx.Unlock()

	res := &oneResult{
		ch:  make(chan faces.IItem, 1),
		ctx: ctx,
	}

	allResults.Store(id, res)
	return res.ch
}

func (m *SystemFinalHandler) Start(_ context.Context) error {
	return nil
}

func (m *SystemFinalHandler) Stop() {
	mx.Lock()
	defer mx.Unlock()

	closeFunc := func(key int64, res *oneResult) bool {
		close(res.ch)
		res.ch = nil
		return true
	}

	allResults.Range(closeFunc)
	allResults = &Map{
		data: map[int64]*oneResult{},
	}
}

func runOne(res *oneResult, item faces.IItem) {
	select {
	case res.ch <- item: /* nothing */
	case <-res.ctx.Done(): /* nothing */
	default: /* nothing */
	}

	// Out of sight, out of mind
	close(res.ch)
	res.ch = nil
}

func (m *SystemFinalHandler) Run(item faces.IItem) error {

	id := item.GetID()
	if value, loaded := allResults.LoadAndDelete(id); loaded {
		go runOne(value, item)
	}

	return nil
}
