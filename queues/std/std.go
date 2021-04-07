// Package std supports queue with standard GO-channels for using them in conveyor.
package std

import (
	"github.com/iostrovok/conveyor/faces"
	"github.com/iostrovok/conveyor/protobuf/go/nodes"
)

// Chan is main package object.
type Chan struct {
	ch       faces.MainCh
	isActive bool
}

// New is a constructor.
func New(length int) faces.IChan {
	return &Chan{
		isActive: true,
		ch:       make(faces.MainCh, length),
	}
}

// Push adds item index to queue.
func (c *Chan) Push(i int) {
	c.ch <- i
}

// IsActive is a simple getter.
func (c *Chan) IsActive() bool {
	return c.isActive
}

// Len returns the max available number items in the stack channel.
func (c *Chan) Len() int {
	return cap(c.ch)
}

// ChanIn returns reference to input channel.
func (c *Chan) ChanIn() faces.MainCh {
	return c.ch
}

// ChanOut returns reference to output channel.
func (c *Chan) ChanOut() faces.MainCh {
	return c.ch
}

// Count returns the number of items in the stack channel.
func (c *Chan) Count() int {
	return len(c.ch)
}

// Close stops the queue.
func (c *Chan) Close() {
	close(c.ch)
	c.isActive = false
}

// Info returns the information about current stage of queue.
func (c *Chan) Info() *nodes.ChanData {
	return &nodes.ChanData{
		Type:            nodes.ChanType_CHAN_STD_GO,
		IsExisted:       true,
		Length:          uint32(c.Len()),
		NumberOfWorkers: 0,
		NumberInCh:      uint32(c.Count()),
	}
}
