package priorityqueue

/*
	The package support the priority queue for using them in conveyor.

	Items with the same priority has FIFO order.
	Realisation has 2-element delay on empty channel because it uses standard GO channel.
*/

import (
	"context"
	"strconv"
	"strings"
	"sync"

	"github.com/iostrovok/conveyor/faces"
	"github.com/iostrovok/conveyor/protobuf/go/nodes"
)

const (
	half = 2
)

type IPQ interface {
	Len() int32
	ChanIn() chan interface{}
	ChanOut() chan interface{}
}

type PQ struct {
	sync.RWMutex

	limit int
	last  int
	chIn  faces.MainCh
	chOut faces.MainCh
	body  []faces.IItem
	cond  chan struct{}

	isActive bool
}

// Create a new stack.
func New(length int) faces.IChan {
	return Init(context.Background(), length)
}

func Init(ctx context.Context, limit int) *PQ {
	stack := &PQ{
		chIn:  make(faces.MainCh, 1),
		chOut: make(faces.MainCh, 1),
		limit: limit,
		last:  0,
		body:  make([]faces.IItem, limit+1, limit+1),
		cond:  make(chan struct{}, limit-1),

		isActive: true,
	}

	go stack.runIn(ctx)
	go stack.runOut(ctx)

	return stack
}

// Print is helper for tests, don't use it for production.
func (pq *PQ) Print() string {
	pq.Lock()
	defer pq.Unlock()

	return PrintBody(pq.body, pq.last)
}

// PrintBody is helper for tests, don't use it for production.
func PrintBody(body []faces.IItem, last int) string {
	out := make([]string, last)
	for i := 0; i < last; i++ {
		item := body[i]
		out[i] = strconv.Itoa(i) + "] " + strconv.FormatInt(item.GetID(), 10) + " => " +
			strconv.Itoa(item.GetPriority())
	}

	return strings.Join(out, "\n")
}

// Body().
func (pq *PQ) Body() []faces.IItem {
	pq.Lock()
	defer pq.Unlock()

	return pq.body
}

func (pq *PQ) Push(item faces.IItem) {
	pq.chIn <- item
}

// Close().
func (pq *PQ) Close() {
	pq.Lock()
	defer pq.Unlock()
	close(pq.chIn)
}

// Returns the number of items in the stack.
func (pq *PQ) Count() int {
	return pq.last
}

// IsActive.
func (pq *PQ) IsActive() bool {
	return pq.isActive
}

// Returns the max available number items in the stack.
func (pq *PQ) Len() int {
	return pq.limit
}

func (pq *PQ) ChanIn() faces.MainCh {
	return pq.chIn
}

func (pq *PQ) ChanOut() faces.MainCh {
	return pq.chOut
}

// Pop the top item of the stack and return it.
func (pq *PQ) runIn(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			close(pq.cond)

			return
		case x, ok := <-pq.chIn:
			if !ok {
				close(pq.cond)

				return
			}

			// put to array
			pq.insert(x)

			select {
			case <-ctx.Done():
				return
			case pq.cond <- struct{}{}: // nothing
			}
		}
	}
}

// Pop the top item of the stack and return it.
func (pq *PQ) runOut(ctx context.Context) {
	sent := 0
	for {
		select {
		case <-ctx.Done():
			return
		case _, ok := <-pq.cond:
			if !ok {
				close(pq.chOut)

				return
			}

			pq.Lock()
			pq.last--
			x := pq.body[pq.last]
			pq.Unlock()

			if x != nil {
				sent++
				select {
				case <-ctx.Done():
					return
				case pq.chOut <- x: // nothing
				}
			}
		}
	}
}

// Pop the top item of the stack and return it.
func (pq *PQ) insert(x faces.IItem) {
	pq.Lock()
	defer func() {
		pq.last++
		pq.Unlock()
	}()

	if pq.last == 0 {
		pq.body[0] = x

		return
	}

	p := x.GetPriority()
	if pq.last == 1 {
		if pq.body[0].GetPriority() >= p {
			pq.body[0], pq.body[1] = x, pq.body[0]
		} else {
			pq.body[1] = x
		}

		return
	}

	if pq.body[0].GetPriority() >= p {
		insertToArray(&(pq.body), x, 0)

		return
	}

	if pq.body[pq.last-1].GetPriority() <= p {
		pq.body[pq.last] = x

		return
	}

	position := findPosition(pq.body, x.GetPriority(), pq.last)
	if position == pq.last {
		pq.body[pq.last] = x
	} else {
		insertToArray(&(pq.body), x, position)
	}
}

// findPosition is a simple binary search.
func findPosition(array []faces.IItem, priority int, last int) int {
	last--

	if priority >= array[last].GetPriority() {
		return last + 1
	}

	if priority <= array[0].GetPriority() {
		return 0
	}

	first, mid := 0, 0

forLoop:
	for first <= last {
		mid = (first + last) / half // half == 2

		switch {
		case array[mid].GetPriority() == priority:
			break forLoop
		case array[mid].GetPriority() < priority:
			first = mid + 1
		default:
			last = mid - 1
		}
	}

	// adjust for not equally priority
	if array[mid].GetPriority() < priority {
		return mid + 1
	}

	return mid
}

func insertToArray(array *[]faces.IItem, item faces.IItem, position int) {
	// shift values
	copy((*array)[position+1:], (*array)[position:])

	// insert value
	(*array)[position] = item
}

func (pq *PQ) Info() *nodes.ChanData {
	return &nodes.ChanData{
		Type:            nodes.ChanType_CHAN_PRIORITY_QUEUE,
		IsExisted:       true,
		Length:          uint32(pq.Len()),
		NumberOfWorkers: 0,
		NumberInCh:      uint32(pq.Count()),
	}
}
