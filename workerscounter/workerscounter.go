// Package workerscounter rules the number of current worked handlers.
package workerscounter

import (
	"github.com/iostrovok/conveyor/faces"
	"github.com/iostrovok/conveyor/protobuf/go/nodes"
)

// WorkersCounter is main structure.
type WorkersCounter struct{}

// New is a constructor.
func New() faces.IWorkersCounter {
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

func makeManagerAction(mc *nodes.ManagerData, action nodes.Action) *nodes.ManagerAction {
	switch action {
	case nodes.Action_NOTHING:
		return &nodes.ManagerAction{Action: nodes.Action_NOTHING, Delta: 0}
	case nodes.Action_UP:
		if mc.Workers.Number < mc.Workers.Max {
			return &nodes.ManagerAction{Action: nodes.Action_UP, Delta: 1}
		}
	case nodes.Action_DOWN:
		if mc.Workers.Number > mc.Workers.Min {
			return &nodes.ManagerAction{Action: nodes.Action_DOWN, Delta: 1}
		}
	}

	return &nodes.ManagerAction{Action: nodes.Action_NOTHING, Delta: 0}
}

// Check checks loading of manager and returns the recommendation to action.
func (m *WorkersCounter) Check(mc *nodes.ManagerData) (*nodes.ManagerAction, error) {
	// the first "input" handler
	if len(mc.ChanBefore) == 0 {
		return makeManagerAction(mc, nodes.Action_NOTHING), nil
	}

	// simple support of the bottom threshold
	if mc.Workers.Number < mc.Workers.Min {
		return makeManagerAction(mc, nodes.Action_UP), nil
	}

	// simple support of the upper threshold
	if mc.Workers.Number > mc.Workers.Min {
		return makeManagerAction(mc, nodes.Action_DOWN), nil
	}

	activeChanBefore, findBefore := findActive(mc.ChanBefore)
	if !findBefore {
		return makeManagerAction(mc, nodes.Action_NOTHING), nil
	}

	if float32(activeChanBefore.NumberInCh) > 0.5*float32(activeChanBefore.Length) {
		return makeManagerAction(mc, nodes.Action_UP), nil
	} else if float32(activeChanBefore.NumberInCh) < 0.1*float32(activeChanBefore.Length) {
		return makeManagerAction(mc, nodes.Action_DOWN), nil
	}

	return makeManagerAction(mc, nodes.Action_NOTHING), nil
}
