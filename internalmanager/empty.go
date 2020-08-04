/*
	Internal package. Package supports the empty IChan interface.
*/
package internalmanager

import (
	"context"

	"github.com/iostrovok/conveyor/faces"
)

type EmptyHandler struct{}

var MakeEmptyHandler faces.GiveBirth = func(name faces.Name) (faces.IHandler, error) {
	return &EmptyHandler{}, nil
}

func (m *EmptyHandler) Start(ctx context.Context) error {
	return nil
}
func (m *EmptyHandler) Stop() { /* nothing */ }

func (m *EmptyHandler) Run(item faces.IItem) error {
	return nil
}
