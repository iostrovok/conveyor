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

var MakeEmptyHandler GiveBirth = func(name Name) (IHandler, error) {
	return &EmptyHandler{}, nil
}

// Start does nothing
func (m *EmptyHandler) Start(_ context.Context) error {
	return nil
}

// Stop does nothing
func (m *EmptyHandler) Stop() { /* nothing */ }

// Run does nothing
func (m *EmptyHandler) Run(_ IItem) error {
	return nil
}

// does nothing
func (m *EmptyHandler) TickerRun(ctx context.Context) { /* nothing */ }

// return 0
func (m *EmptyHandler) TickerDuration() time.Duration {
	return time.Duration(0)
}
