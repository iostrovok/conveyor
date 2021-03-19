package faces

import (
	"context"
	"time"
)

// GiveBirth return new handler. Type Name is string which was passed with AddHandler(...).
type GiveBirth func(name Name) (IHandler, error)

// IHandler is interface for support the single handler.
type IHandler interface {
	// Start() function is called one time for each handler right after it's created with GiveBirth().
	Start(ctx context.Context) error

	// Run() function is called for processing single item.
	Run(item IItem) error

	// TickerRun() function is called for processing by timer (ticker).
	TickerRun(ctx context.Context) // error is not processing
	TickerDuration() time.Duration

	// Stop() function is called before destruction of handler.
	Stop(ctx context.Context)
}
