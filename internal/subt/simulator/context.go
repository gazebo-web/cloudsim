package simulator

import (
	"errors"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
)

var (
	// contextPlatform is the key used to get the platform from the context.
	contextPlatform = "cloudsim_platform"

	// contextServices is the key used to get the services from the context.
	contextServices = "subt_services"

	// ErrCtxInvalidPlatform is returned as a panic when casting a platform retrieved from context fails.
	ErrCtxInvalidPlatform = errors.New("invalid platform from context")

	// ErrCtxInvalidAppServices is returned as a panic when casting a group of services retrieved from context fails.
	ErrCtxInvalidAppServices = errors.New("invalid application services from context")
)

// Context is an action's context wrapper that exposes the methods needed to be accessed during jobs.
type Context interface {
	actions.Context
	// Platform returns the platform from context.
	Platform() platform.Platform

	// Services returns the services from context.
	Services() application.Services
}

// simulatorContext is a Context implementation.
type simulatorContext struct {
	actions.Context
}

// Platform gets the platform from context and returns it.
// It panics if the casting fails.
func (ctx *simulatorContext) Platform() platform.Platform {
	value := ctx.Value(contextPlatform)
	output, ok := value.(platform.Platform)
	if !ok {
		panic(ErrCtxInvalidPlatform)
	}
	return output
}

// Services gets the services from context and returns it.
// It panics if the casting fails.
func (ctx *simulatorContext) Services() application.Services {
	value := ctx.Value(contextServices)
	output, ok := value.(application.Services)
	if !ok {
		panic(ErrCtxInvalidAppServices)
	}
	return output
}

// NewContext initializes a new Context implementation from the base action's context.
func NewContext(ctx actions.Context) Context {
	return &simulatorContext{
		Context: ctx,
	}
}
