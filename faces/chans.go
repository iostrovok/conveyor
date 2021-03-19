package faces

// File describes the chan interface.

import (
	"github.com/iostrovok/conveyor/protobuf/go/nodes"
)

type (
	// MainCh is a global type.
	MainCh chan IItem

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
	// Return
	ChanIn() MainCh
	ChanOut() MainCh
	Push(IItem)

	Close()
	IsActive() bool
	Count() int
	Len() int
	Info() *nodes.ChanData
}
