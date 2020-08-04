/*
	Internal package. Package provides a handler for supporting online processing.
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

// global vars
var allResults *sync.Map
var mx *sync.RWMutex

type OnwFinalHandler struct{}

func init() {
	mx = new(sync.RWMutex)
	allResults = &sync.Map{}
}

func Init() faces.GiveBirth {
	return func(name faces.Name) (faces.IHandler, error) {
		return &OnwFinalHandler{}, nil
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

func (m *OnwFinalHandler) Start(_ context.Context) error {
	return nil
}

func (m *OnwFinalHandler) Stop() {
	mx.Lock()
	defer mx.Unlock()

	closeFunc := func(key, value interface{}) bool {
		res := value.(*oneResult)
		close(res.ch)
		res.ch = nil
		return true
	}

	allResults.Range(closeFunc)
	allResults = &sync.Map{}
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

func (m *OnwFinalHandler) Run(item faces.IItem) error {

	id := item.GetID()
	if value, loaded := allResults.LoadAndDelete(id); loaded {
		go runOne(value.(*oneResult), item)
	}

	return nil
}
