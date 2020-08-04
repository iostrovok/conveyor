/*
	The package support queue with standard GO-channels for using them in conveyor.
*/
package std

import (
	"github.com/iostrovok/conveyor/faces"
	"github.com/iostrovok/conveyor/protobuf/go/nodes"
)

type Chan struct {
	ch       faces.MainCh
	isActive bool
}

func New(length int) faces.IChan {
	return &Chan{
		isActive: true,
		ch:       make(faces.MainCh, length),
	}
}

func (c *Chan) Push(item faces.IItem) {
	c.ch <- item
}

func (c *Chan) IsActive() bool {
	return c.isActive
}

func (c *Chan) Len() int {
	return cap(c.ch)
}

func (c *Chan) ChanIn() faces.MainCh {
	return c.ch
}

func (c *Chan) ChanOut() faces.MainCh {
	return c.ch
}

func (c *Chan) Count() int {
	return len(c.ch)
}

func (c *Chan) Close() {
	close(c.ch)
	c.isActive = false
}

func (c *Chan) Info() *nodes.ChanData {
	return &nodes.ChanData{
		Type:            nodes.ChanType_CHAN_STD_GO,
		IsExisted:       true,
		Length:          uint32(c.Len()),
		NumberOfWorkers: 0,
		NumberInCh:      uint32(c.Count()),
	}
}
