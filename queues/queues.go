/*
Package queues implements the IChan interface and provides 3 simples realization of them.

- standard GO channel FIFO, the fastest realization
- stack, LIFO
- priority queues
*/
package queues

import (
	"github.com/iostrovok/conveyor/faces"
	"github.com/iostrovok/conveyor/queues/priorityqueue"
	"github.com/iostrovok/conveyor/queues/stack"
	"github.com/iostrovok/conveyor/queues/std"
)

// New is a wrapper over selected type of queues.
func New(wb faces.IWorkBench, chanType faces.ChanType) faces.IChan {

	lengthChannel := wb.Len() + 1

	switch chanType {
	case faces.ChanStdGo:
		return std.New(lengthChannel)
	case faces.ChanStack:
		return stack.New(lengthChannel)
	case faces.ChaPriorityQueue:
		return priorityqueue.New(wb, lengthChannel)
	}

	return nil
}
