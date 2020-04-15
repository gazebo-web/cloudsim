package subt

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
)

// SubT is an IApplication implementation
type SubT struct {
	*application.Application
}

// New creates a new SubT application.
func New(p *platform.Platform) *SubT {
	app := application.New(p)
	subt := &SubT{
		Application: app,
	}
	repository := simulations.NewRepository(p.Server.Db)
	app.Services.Simulation = simulations.NewService(repository)
	return subt
}

// Name returns the SubT application's name.
func (app *SubT) Name() string {
	return "subt"
}

// Version returns the SubT application's version.
func (app *SubT) Version() string {
	return "2.0"
}

// Register runs a set of instructions to initialize an application for the given platform.
func Register(p *platform.Platform) application.IApplication {
	return New(p)
}
