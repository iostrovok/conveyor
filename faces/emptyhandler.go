package faces

import (
	"context"
	"time"
)

/*
EmptyHandler realizes the empty handler which do nothing.
It is useful for replace empty Start and Stop methods by inherited objects.

Example:

	type MySimpleHandler struct {
	   EmptyHandler
	}

	func Handler(_ faces.Name) (faces.IHandler, error) {
		return &MySimpleHandler{}, nil
	}

	func (m *MySimpleHandler) Run(item faces.IItem) error {s
		// do something here
		return nil
	}
*/
type EmptyHandler struct{}

// MakeEmptyHandler is a constructor for EmptyHandler.
func MakeEmptyHandler(_ Name) (IHandler, error) {
	return &EmptyHandler{}, nil
}

// Start does nothing.
func (m *EmptyHandler) Start(_ context.Context) error {
	return nil
}

// Stop does nothing.
func (m *EmptyHandler) Stop(_ context.Context) { /* nothing */ }

// Run does nothing.
func (m *EmptyHandler) Run(_ IItem) error {
	return nil
}

// It does nothing by default.
func (m *EmptyHandler) TickerRun(ctx context.Context) { /* nothing */ }

// returns 0 by default.
func (m *EmptyHandler) TickerDuration() time.Duration {
	return time.Duration(0)
}
