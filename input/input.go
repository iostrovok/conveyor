/*
	Package realizes the IItem interface.
*/
package input

import (
	"context"

	"github.com/iostrovok/conveyor/faces"
)

type Input struct {
	data     interface{}
	ctx      context.Context
	tracer   faces.ITrace
	priority *int
	name     faces.Name

	testSuffix string
}

func New() faces.IInput {
	return &Input{}
}

func (i *Input) Context(ctx context.Context) faces.IInput {
	i.ctx = ctx
	return i
}

func (i *Input) Trace(tracer faces.ITrace) faces.IInput {
	i.tracer = tracer
	return i
}

func (i *Input) Data(data interface{}) faces.IInput {
	i.data = data
	return i
}

func (i *Input) Priority(priority int) faces.IInput {
	i.priority = &priority
	return i
}

func (i *Input) SkipToName(name faces.Name) faces.IInput {
	i.name = name
	return i
}

func (i *Input) Values() (context.Context, faces.ITrace, interface{}, *int, faces.Name) {
	return i.ctx, i.tracer, i.data, i.priority, i.name
}

func (i *Input) Ctx() context.Context {
	return i.ctx
}
