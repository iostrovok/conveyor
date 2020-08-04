package faces

import (
	"github.com/iostrovok/conveyor/protobuf/go/nodes"
)

/*
	File describes the chan interface.

*/

type MainCh chan IItem
type ChanType nodes.ChanType

const (
	ChanStdGo        = ChanType(nodes.ChanType_CHAN_STD_GO)
	ChanStack        = ChanType(nodes.ChanType_CHAN_STACK)
	ChaPriorityQueue = ChanType(nodes.ChanType_CHAN_PRIORITY_QUEUE)
)

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
