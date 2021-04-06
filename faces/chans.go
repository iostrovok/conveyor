package faces

// File describes the chan interface.

import (
	"github.com/iostrovok/conveyor/protobuf/go/nodes"
)

type (
	// MainCh is a global type.
	MainCh chan int

	// ChanType is a global type.
	ChanType nodes.ChanType
)

const (
	// ChanStdGo is wrapper for nodes.ChanType_CHAN_STD_GO.
	ChanStdGo = ChanType(nodes.ChanType_CHAN_STD_GO)

	// ChanStack is wrapper for nodes.ChanType_CHAN_STACK.
	ChanStack = ChanType(nodes.ChanType_CHAN_STACK)

	// ChaPriorityQueue is wrapper for nodes.ChanType_CHAN_PRIORITY_QUEUE.
	ChaPriorityQueue = ChanType(nodes.ChanType_CHAN_PRIORITY_QUEUE)
)

// IChan is interface for support queue oin conveyor.
type IChan interface {

	// ChanIn returns reference to input channel.
	ChanIn() MainCh
	// ChanOut returns reference to output channel.
	ChanOut() MainCh

	// Push adds item index to queue.
	Push(int)
	Close()
	IsActive() bool
	// Count returns the number of items in the stack channel.
	Count() int
	// Len returns the max available number items in the stack channel.
	Len() int

	Info() *nodes.ChanData
}
