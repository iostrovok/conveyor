/*
Package input implements the faces.IInput interface.
*/
package input

import (
	"context"

	"github.com/iostrovok/conveyor/faces"
)

// Input is an implementation of faces.IInput Interface .
type Input struct {
	data     interface{}
	ctx      context.Context
	tracer   faces.ITrace
	priority *int
	name     faces.Name
}

// New is a constructor.
func New() faces.IInput {
	return &Input{}
}

// Context is a simple setter.
func (i *Input) Context(ctx context.Context) faces.IInput {
	i.ctx = ctx

	return i
}

// Trace is a simple setter.
func (i *Input) Trace(tracer faces.ITrace) faces.IInput {
	i.tracer = tracer

	return i
}

// Data is a simple setter.
func (i *Input) Data(data interface{}) faces.IInput {
	i.data = data

	return i
}

// Priority is a simple setter.
func (i *Input) Priority(priority int) faces.IInput {
	i.priority = &priority

	return i
}

// SkipToName is a simple setter.
func (i *Input) SkipToName(name faces.Name) faces.IInput {
	i.name = name

	return i
}

// Values is a simple getter.
func (i *Input) Values() (context.Context, faces.ITrace, interface{}, *int, faces.Name) {
	return i.ctx, i.tracer, i.data, i.priority, i.name
}

// Ctx is a simple getter.
func (i *Input) Ctx() context.Context {
	return i.ctx
}
