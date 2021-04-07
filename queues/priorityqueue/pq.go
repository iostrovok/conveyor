/*
Package priorityqueue supports the priority queues for using them in conveyor.

Items with the same priority has FIFO order.
Realization has 2-element delay on empty channel because it uses standard GO channel.

*/
package priorityqueue

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

// PQ is main package structure.
type PQ struct {
	sync.RWMutex

	limit int
	last  int
	chIn  faces.MainCh
	chOut faces.MainCh
	body  []int
	cond  chan struct{}
	wb    faces.IWorkBench

	isActive bool
}

// New is a constructor, creates a new stack.
func New(wb faces.IWorkBench, length int) faces.IChan {
	return Init(context.Background(), wb, length)
}

// Init is full constructor for accurate configuration.
func Init(ctx context.Context, wb faces.IWorkBench, limit int) *PQ {
	stack := &PQ{
		chIn:  make(faces.MainCh, 1),
		chOut: make(faces.MainCh, 1),
		limit: limit,
		last:  0,
		body:  make([]int, limit+1, limit+1),
		cond:  make(chan struct{}, limit-1),

		isActive: true,
		wb:       wb,
	}

	for i := range stack.body {
		stack.body[i] = -1
	}

	go stack.runIn(ctx)
	go stack.runOut(ctx)

	return stack
}

// Print is helper for tests, don't use it for production.
func (pq *PQ) Print() string {
	pq.Lock()
	defer pq.Unlock()

	out := make([]string, pq.last)
	for i := 0; i < pq.last; i++ {
		item, err := pq.wb.Get(pq.body[i])
		if err == nil {
			out[i] = strconv.Itoa(i) + "] " + strconv.FormatInt(item.GetID(), 10) + " => " +
				strconv.Itoa(item.GetPriority())
		}
	}

	return strings.Join(out, "\n")
}

// Body is debug function.
func (pq *PQ) Body() []int {
	pq.Lock()
	defer pq.Unlock()

	return pq.body
}

// Push adds item index to queue.
func (pq *PQ) Push(i int) {
	pq.chIn <- i
}

// Close stops the queue.
func (pq *PQ) Close() {
	pq.Lock()
	defer pq.Unlock()
	close(pq.chIn)
}

// Count returns the number of items in the stack channel.
func (pq *PQ) Count() int {
	return pq.last
}

// IsActive is a simple getter.
func (pq *PQ) IsActive() bool {
	return pq.isActive
}

// Len returns the max available number items in the stack channel.
func (pq *PQ) Len() int {
	return pq.limit
}

// ChanIn returns reference to input channel.
func (pq *PQ) ChanIn() faces.MainCh {
	return pq.chIn
}

// ChanOut returns reference to output channel.
func (pq *PQ) ChanOut() faces.MainCh {
	return pq.chOut
}

// runIn gets the item index and save it in queue
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

// runOut pops the top item of the queue and returns it.
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
			pq.body[pq.last] = 1 // just clean
			pq.Unlock()

			if x > -1 {
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
func (pq *PQ) insert(x int) {
	pq.Lock()
	defer func() {
		pq.last++
		pq.Unlock()
	}()

	if pq.last == 0 {
		pq.body[0] = x

		return
	}

	p := pq.wb.GetPriority(x)
	if pq.last == 1 {
		if pq.wb.GetPriority(0) >= p {
			pq.body[0], pq.body[1] = x, pq.body[0]
		} else {
			pq.body[1] = x
		}

		return
	}

	if pq.wb.GetPriority(0) >= p {
		insertToArray(&(pq.body), x, 0)

		return
	}

	if pq.wb.GetPriority(pq.last-1) <= p {
		pq.body[pq.last] = x

		return
	}

	position := findPosition(pq.wb, pq.wb.GetPriority(x), pq.last)
	if position == pq.last {
		pq.body[pq.last] = x
	} else {
		insertToArray(&(pq.body), x, position)
	}
}

// findPosition is a simple binary search.
func findPosition(wb faces.IWorkBench, priority int, lastItemIndex int) int {
	lastItemIndex--

	if priority >= wb.GetPriority(lastItemIndex) {
		return lastItemIndex + 1
	}

	if priority <= wb.GetPriority(0) {
		return 0
	}

	first, mid := 0, 0

forLoop:
	for first <= lastItemIndex {
		mid = (first + lastItemIndex) / half // half == 2

		p := wb.GetPriority(mid)
		switch {
		case p == priority:
			break forLoop
		case p < priority:
			first = mid + 1
		default:
			lastItemIndex = mid - 1
		}
	}

	// adjust for not equally priority
	if wb.GetPriority(mid) < priority {
		return mid + 1
	}

	return mid
}

func insertToArray(array *[]int, itemIndex int, position int) {
	// shift values
	copy((*array)[position+1:], (*array)[position:])

	// insert value
	(*array)[position] = itemIndex
}

// Info returns the information about current stage of queue.
func (pq *PQ) Info() *nodes.ChanData {
	return &nodes.ChanData{
		Type:            nodes.ChanType_CHAN_PRIORITY_QUEUE,
		IsExisted:       true,
		Length:          uint32(pq.Len()),
		NumberOfWorkers: 0,
		NumberInCh:      uint32(pq.Count()),
	}
}
