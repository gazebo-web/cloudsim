package subt

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
)

// SubT is an IApplication implementation
type SubT struct {}

// New creates a new SubT application.
func New() SubT {
	return SubT{}
}

// Name returns the SubT application name.
func (s SubT) Name() string {
	return "subt"
}

// Register creates a New application to be registered in the platform.
func Register() *application.IApplication {
	var subt application.IApplication
	subt = New()
	return &subt
}
