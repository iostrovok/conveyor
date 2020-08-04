/*
	The package support the stack queue LIFO (or FILO) for using them in conveyor.
*/
package stack

/*

 */

import (
	"context"
	"sync"

	"github.com/iostrovok/conveyor/faces"
	"github.com/iostrovok/conveyor/protobuf/go/nodes"
)

type IStack interface {
	Len() int32
	ChanIn() chan interface{}
	ChanOut() chan interface{}
}

type Stack struct {
	sync.RWMutex

	current int
	limit   int
	last    int
	chIn    faces.MainCh
	chOut   faces.MainCh
	body    []faces.IItem
	cond    chan struct{}

	isActive bool
}

// Create a new stack
func New(length int) faces.IChan {
	return Init(length, context.Background())
}

func Init(limit int, ctx context.Context) *Stack {

	stack := &Stack{
		chIn:  make(faces.MainCh, 1),
		chOut: make(faces.MainCh, 1),
		limit: limit,
		last:  0,
		body:  make([]faces.IItem, limit+1, limit+1),
		cond:  make(chan struct{}, limit+1),

		isActive: true,
	}

	go stack.runIn(ctx)
	go stack.runOut(ctx)

	return stack
}

func (stack *Stack) Push(item faces.IItem) {
	stack.chIn <- item
}

// Close()
func (stack *Stack) Close() {
	stack.Lock()
	defer stack.Unlock()
	close(stack.chIn)
}

// Returns the number of items in the stack
func (stack *Stack) Count() int {
	return stack.last
}

// IsActive
func (stack *Stack) IsActive() bool {
	return stack.isActive
}

// Returns the max available number items in the stack
func (stack *Stack) Len() int {
	return int(stack.limit)
}

func (stack *Stack) ChanIn() faces.MainCh {
	return stack.chIn
}

func (stack *Stack) ChanOut() faces.MainCh {
	return stack.chOut
}

func (stack *Stack) runIn(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case x, ok := <-stack.chIn:
			if !ok {
				close(stack.cond)
				return
			}

			// put to array
			stack.Lock()
			stack.body[stack.last] = x
			stack.last++
			stack.Unlock()

			select {
			case <-ctx.Done():
				return
			case stack.cond <- struct{}{}:
				// nothing
			}
		}
	}
}

func (stack *Stack) runOut(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case _, ok := <-stack.cond:
			if !ok {
				close(stack.chOut)
				return
			}

			stack.Lock()
			stack.last--
			x := stack.body[stack.last]
			stack.Unlock()

			if x != nil {
				select {
				case <-ctx.Done():
					return
				case stack.chOut <- x:
					// nothing
				}
			}
		}
	}
}

func (stack *Stack) Info() *nodes.ChanData {
	return &nodes.ChanData{
		Type:            nodes.ChanType_CHAN_STACK,
		IsExisted:       true,
		Length:          uint32(stack.Len()),
		NumberOfWorkers: 0,
		NumberInCh:      uint32(stack.Count()),
	}
}
