package subt

import (
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
	return subt
}

// Name returns the SubT application name.
func (app *SubT) Name() string {
	return "subt"
}

func (app *SubT) Version() string {
	return "2.0"
}

// Register creates a New application to be registered in the platform.
func Register(p *platform.Platform) application.IApplication {
	return New(p)
}
