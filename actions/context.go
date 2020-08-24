package actions

import (
	ctx "context"
)

// Context is used to send context data to jobs. Context should be used to provide access to platforms, services,
// loggers, etc.
// Applications are free to add values to a Context through a context.Context object, but should be careful not to
// replace values set by other actors. It is recommended to create a set of constants to use as keys, where each key is
// composed of an application prefix followed by an identifier.
type Context interface {
	Ctx() ctx.Context
}

type context struct {
	ctx ctx.Context
}

func (c *context) Ctx() ctx.Context {
	return c.ctx
}

// NewContext returns a new Context to pass context information to action jobs.
func NewContext(ctx ctx.Context) Context {
	return &context{
		ctx: ctx,
	}
}
