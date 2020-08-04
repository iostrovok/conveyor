// WorkersCounter rules the number of current worked handlers.
package workerscounter

import (
	"github.com/iostrovok/conveyor/faces"
	"github.com/iostrovok/conveyor/protobuf/go/nodes"
)

/*
	....
*/

type WorkersCounter struct {
}

func NewWorkersCounter() faces.IWorkersCounter {
	return &WorkersCounter{}
}

func findActive(chs []*nodes.ChanData) (*nodes.ChanData, bool) {
	for _, ch := range chs {
		if ch.IsExisted {
			return ch, true
		}
	}
	return nil, false
}

func (m *WorkersCounter) Check(mc *nodes.ManagerData) (*nodes.ManagerAction, error) {

	out := &nodes.ManagerAction{
		Action: nodes.Action_NOTHING,
		Delta:  0,
	}

	if len(mc.ChanBefore) == 0 {
		return out, nil
	}

	if mc.Workers.Number < mc.Workers.Min {
		out.Action = nodes.Action_UP
		out.Delta = 1
		return out, nil
	}

	if activeChanAfter, find := findActive(mc.ChanAfter); find && activeChanAfter.Length > 0 {
		// if next manager can't do his work...
		if float32(activeChanAfter.NumberInCh) > 0.9*float32(activeChanAfter.Length) {
			if mc.Workers.Number > mc.Workers.Min {
				out.Delta = 1
				out.Action = nodes.Action_DOWN
				return out, nil
			}
		}
	}

	activeChanBefore, findBefore := findActive(mc.ChanBefore)
	if !findBefore {
		return out, nil
	}

	if float32(activeChanBefore.NumberInCh) > 0.5*float32(activeChanBefore.Length) {
		if mc.Workers.Number < mc.Workers.Max {
			out.Delta = 1
			out.Action = nodes.Action_UP
			return out, nil
		}
	} else if float32(activeChanBefore.NumberInCh) < 0.1*float32(activeChanBefore.Length) {
		if mc.Workers.Number > mc.Workers.Min {
			out.Delta = 1
			out.Action = nodes.Action_DOWN
			return out, nil
		}
	}

	return out, nil
}
