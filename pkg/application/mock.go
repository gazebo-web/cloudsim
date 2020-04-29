package application

import (
	"context"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/tasks"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

type Mock struct {
}

func (m *Mock) Name() string {
	return "application_test"
}

func (m *Mock) Version() string {
	return "1.0"
}

func (m *Mock) RegisterTasks() []tasks.Task {
	panic("implement me")
}

func (m *Mock) RegisterMonitors(ctx context.Context) {
	panic("implement me")
}

func (m *Mock) RebuildState(ctx context.Context) error {
	panic("implement me")
}

func (m *Mock) Shutdown(ctx context.Context) error {
	panic("implement me")
}

func (m *Mock) Launch(payload interface{}) (interface{}, *ign.ErrMsg) {
	panic("implement me")
}

func (m *Mock) ValidateLaunch(ctx context.Context, simulation *simulations.Simulation) error {
	panic("implement me")
}

func NewMock() *Mock {
	app := &Mock{}
	return app
}
