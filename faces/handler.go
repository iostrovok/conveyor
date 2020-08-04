package faces

import (
	"context"
)

/*
	....
*/

// GiveBirth return new handler. Name sis string which was passed with AddHandler(...)
type GiveBirth func(name Name	) (IHandler, error)

type IHandler interface {
	// The Start() function is called one time for each handler right after it's created with GiveBirth().
	Start(ctx context.Context) error

	// The Run() function is called for processing single item.
	Run(item IItem) error

	// The Stop() function is called before destruction of handler.
	Stop()
}
