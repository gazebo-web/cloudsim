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

// Application is a generic implementation of an application to be extended by a specific application.
type Application struct {
	Platform *platform.Platform
}

// New creates a new application for the given platform.
func New(p *platform.Platform) *Application {
	app := &Application{
		Platform: p,
	}
	return app
}

// Name returns the application's name.
// Needs to be implemented by the specific application.
func (app *Application) Name() string {
	panic("Name should be implemented by the application")
}

// Version returns the application's version.
// If the specific application doesn't implement this method, it will return 1.0.
func (app *Application) Version() string {
	return "1.0"
}

// RegisterRoutes returns the slice of the application's routes.
// Needs to be implemented by the specific application.
func (app *Application) RegisterRoutes() ign.Routes {
	panic("RegisterRoutes should be implemented by the application")
}

// RegisterTasks returns an array of the tasks that need to be executed by the scheduler.
// If the specific application doesn't implement this method, it will return an empty slice.
func (app *Application) RegisterTasks() []tasks.Task {
	return []tasks.Task{}
}