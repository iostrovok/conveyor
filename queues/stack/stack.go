// Package stack supports the stack queues LIFO (or FILO) for using them in conveyor.
package stack

import (
	"context"
	"sync"

	"github.com/iostrovok/conveyor/faces"
	"github.com/iostrovok/conveyor/protobuf/go/nodes"
)

// Stack is main package object.
type Stack struct {
	sync.RWMutex

	limit int
	last  int
	chIn  faces.MainCh
	chOut faces.MainCh
	body  []int
	cond  chan struct{}

	isActive bool
}

// New is a constructor, creates a new stack.
func New(length int) faces.IChan {
	return Init(context.Background(), length)
}

// Init is full constructor for accurate configuration.
func Init(ctx context.Context, limit int) *Stack {
	stack := &Stack{
		chIn:  make(faces.MainCh, 1),
		chOut: make(faces.MainCh, 1),
		limit: limit,
		last:  0,
		body:  make([]int, limit+1, limit+1),
		cond:  make(chan struct{}, limit+1),

		isActive: true,
	}

	for i := range stack.body {
		stack.body[i] = -1
	}

	go stack.runIn(ctx)
	go stack.runOut(ctx)

	return stack
}

// Push adds item index to queue.
func (stack *Stack) Push(i int) {
	stack.chIn <- i
}

// Close stops the queue.
func (stack *Stack) Close() {
	stack.Lock()
	defer stack.Unlock()
	close(stack.chIn)
}

// Count returns the number of items in the stack channel.
func (stack *Stack) Count() int {
	return stack.last
}

// IsActive is a simple getter.
func (stack *Stack) IsActive() bool {
	return stack.isActive
}

// Len returns the max available number items in the stack channel.
func (stack *Stack) Len() int {
	return stack.limit
}

// ChanIn returns reference to input channel.
func (stack *Stack) ChanIn() faces.MainCh {
	return stack.chIn
}

// ChanOut returns reference to output channel.
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
			case stack.cond <- struct{}{}: // nothing
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
			stack.body[stack.last] = -1
			stack.Unlock()

			if x > -1 {
				select {
				case <-ctx.Done():
					return
				case stack.chOut <- x: // nothing
				}
			}
		}
	}
}

// Info returns the information about current stage of queue.
func (stack *Stack) Info() *nodes.ChanData {
	stack.RLock()
	defer stack.RUnlock()

	return &nodes.ChanData{
		Type:            nodes.ChanType_CHAN_STACK,
		IsExisted:       true,
		Length:          uint32(stack.Len()),
		NumberOfWorkers: 0,
		NumberInCh:      uint32(stack.Count()),
	}
}
