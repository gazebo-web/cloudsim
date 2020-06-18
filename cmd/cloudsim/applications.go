package main

import (
	"fmt"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/router"
)

// RegisterApplications registers the given applications by calling their Register method.
// The platform is passed to each application in order to allow the application have a reference to the platform.
func RegisterApplications(p *platform.platform, apps *map[string]application.Application) {
	RegisterApplication(apps, subt.Register(p))
	// RegisterApplication(p.Applications, app.Register(p))
}

// RegisterApplication sets the given application on the map of applications.
func RegisterApplication(applications *map[string]application.Application, app application.Application) {
	if app == nil || applications == nil {
		panic("Invalid application")
	}
	name := app.Name()
	(*applications)[name] = app
}

// RebuildState calls the RebuildState method for all the given applications.
func RebuildState(p *platform.platform, applications map[string]application.Application) {
	for _, app := range applications {
		if err := app.RebuildState(p.Context()); err != nil {
			panic(fmt.Sprintf("Error rebuilding state for application. Name: %s. Version: %s", app.Name(), app.Version()))
		}
	}
}

// RegisterMonitors calls the RegisterMonitors method for all the given applications.
func RegisterMonitors(p *platform.platform, applications map[string]application.Application) {
	for _, app := range applications {
		app.RegisterMonitors(p.Context())
	}
}

// ShutdownApplications calls the Stop method for all given applications.
func ShutdownApplications(p *platform.platform, applications map[string]application.Application) {
	for _, app := range applications {
		if err := app.Stop(p.Context()); err != nil {
			panic(fmt.Sprintf("Error shutting down an application. Name: %s. Version: %s", app.Name(), app.Version()))
		}
	}
}

// RegisterRoutes appends an slice of routes by the given applications to the platform's router.
func RegisterRoutes(p *platform.platform, apps map[string]application.Application) {
	for _, app := range apps {
		router.ConfigureRoutes(p.Server, app.Version(), app.Name(), app.RegisterRoutes())
	}
}

// ScheduleTasks gets all the tasks from each application and add them to the platform's scheduler.
func ScheduleTasks(p *platform.platform, apps map[string]application.Application) {
	for _, app := range apps {
		tasks := app.RegisterTasks()
		for _, task := range tasks {
			p.Scheduler().DoAt(task.Job, task.Date)
		}
	}
}

func RegisterValidators(p *platform.platform, apps map[string]application.Application) {
	for _, app := range apps {
		app.RegisterValidators(p.Context())
	}
}
