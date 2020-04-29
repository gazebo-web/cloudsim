package application

import (
	"context"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/tasks"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

type Mock struct {
	RegisterRoutesMock func() ign.Routes
	RegisterTasksMock func() []tasks.Task
	RegisterMonitorsMock func(ctx context.Context)
	RebuildStateMock func(ctx context.Context) error
	ShutdownMock func(ctx context.Context) error
	LaunchMock func(ctx context.Context, simulation *simulations.Simulation) error
	ValidateLaunchMock func(ctx context.Context, simulation *simulations.Simulation) error
}

func (m *Mock) Name() string {
	return "application_test"
}

func (m *Mock) Version() string {
	return "1.0"
}

func (m *Mock) RegisterRoutes() ign.Routes {
	return m.RegisterRoutesMock()
}

func (m *Mock) RegisterTasks() []tasks.Task {
	return m.RegisterTasksMock()
}

func (m *Mock) RegisterMonitors(ctx context.Context) {
	m.RegisterMonitorsMock(ctx)
}

func (m *Mock) RebuildState(ctx context.Context) error {
	return m.RebuildStateMock(ctx)
}

func (m *Mock) Shutdown(ctx context.Context) error {
	return m.ShutdownMock(ctx)
}

func (m *Mock) Launch(ctx context.Context, simulation *simulations.Simulation) error {
	return m.LaunchMock(ctx, simulation)
}

func (m *Mock) ValidateLaunch(ctx context.Context, simulation *simulations.Simulation) error {
	return m.ValidateLaunchMock(ctx, simulation)
}

func NewMock() *Mock {
	app := &Mock{}
	return app
}