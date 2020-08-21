package actions

import (
	"context"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
)

// Context is used to pass a certain amount of elements needed when running jobs.
type Context interface {
	// Context returns the base context.
	Context() context.Context
	// Platform returns a certain platform to launch simulations.
	Platform() platform.Platform
	// Services returns the services needed to launch simulations.
	Services() application.Services
}

type actionCtx struct {
	ctx      context.Context
	platform platform.Platform
	services application.Services
}

// Context returns the base context.
func (c *actionCtx) Context() context.Context {
	return c.ctx
}

// Platform returns a certain platform to launch simulations.
func (c *actionCtx) Platform() platform.Platform {
	return c.platform
}

// Services returns the services needed to launch simulations.
func (c *actionCtx) Services() application.Services {
	return c.services
}

// NewContext initializes a new action's context from the given base context, platform and services.
func NewContext(ctx context.Context, platform platform.Platform, services application.Services) Context {
	return &actionCtx{
		ctx:      ctx,
		platform: platform,
		services: services,
	}
}
