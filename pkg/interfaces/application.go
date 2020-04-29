package interfaces

import (
	"context"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/tasks"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

// IApplication describes a set of methods for an Application.
type IApplication interface {
	Name() string
	Version() string
	RegisterRoutes() ign.Routes
	RegisterTasks() []tasks.Task
	RegisterMonitors(ctx context.Context)
	RebuildState(ctx context.Context) error
	Shutdown(ctx context.Context) error
	Launch(ctx context.Context, simulation *simulations.Simulation) error
	ValidateLaunch(ctx context.Context, simulation *simulations.Simulation) error
}
