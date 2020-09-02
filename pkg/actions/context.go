package actions

import (
	ctx "context"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

// Context is used to send context data to jobs. Context should be used to provide access to platforms, services,
// loggers, etc.
// Applications are free to add values to a Context through a context.Context object, but should be careful not to
// replace values set by other actors. It is recommended to create a set of constants to use as keys, where each key is
// composed of an application prefix followed by an identifier.
// Note that any information that must be passed between jobs should not be passed using a context, and should be
// returned by the job instead. This is because return values are automatically persisted, which allows an action to
// recover after an unexpected stop (e.g. server restart), while context is lost. The context should only be used to
// pass application-specific values used by jobs, and can be used to support a simple dependency injection scheme.
type Context interface {
	ctx.Context
	Logger() ign.Logger
}

type context struct {
	ctx.Context
	logger ign.Logger
}

// TODO: Open branch to add logger
func (c *context) Logger() ign.Logger {
	return c.logger
}

// NewContext returns a new Context to pass context information to action jobs.
func NewContext(ctx ctx.Context, logger ign.Logger) Context {
	return &context{
		Context: ctx,
		logger:  logger,
	}
}
