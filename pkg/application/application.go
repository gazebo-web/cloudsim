package application

import "gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"

// IApplication describes a set of methods for an Application.
type IApplication interface {
	Name() string
	Platform() platform.IPlatform
	Instance() Application
}

type Application struct {
	parent platform.IPlatform
}

func New(p platform.IPlatform) IApplication {
	var app IApplication
	app = &Application{
		parent: p,
	}
	return app
}

func (app Application) Name() string {
	panic("Name should be implemented by the application")
}

func (app Application) Platform() platform.IPlatform {
	return app.parent
}

func (app Application) Instance() Application {
	return app
}