package subt

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
)

// SubT is an IApplication implementation
type SubT struct {
	application.Application
}

// New creates a new SubT application.
func New(p *platform.Platform) application.IApplication {
	var subt application.IApplication
	app := application.New(p)

	subt = &SubT{
		Application: app.Instance(),
	}
	return subt
}

// Name returns the SubT application name.
func (s SubT) Name() string {
	return "subt"
}

// Register creates a New application to be registered in the platform.
func Register(p *platform.Platform) application.IApplication {
	return New(p)
}
