package application

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/tasks"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

// IApplication describes a set of methods for an Application.
type IApplication interface {
	Name() string
	Version() string
	RegisterRoutes() ign.Routes
	RegisterTasks() []tasks.Task
}

type Application struct {
	Platform *platform.Platform
}

func New(p *platform.Platform) *Application {
	app := &Application{
		Platform: p,
	}
	return app
}

func (app *Application) Name() string {
	panic("Name should be implemented by the application")
}

func (app *Application) Version() string {
	return "1.0"
}

func (app *Application) RegisterRoutes() ign.Routes {
	panic("RegisterRoutes should be implemented by the application")
}

func (app *Application) RegisterTasks() []tasks.Task {
	return []tasks.Task{}
}