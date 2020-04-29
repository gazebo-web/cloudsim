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
func RegisterApplications(p *platform.Platform, apps *map[string]interfaces.IApplication) {
	RegisterApplication(apps, subt.Register(p))
	// RegisterApplication(p.Applications, app.Register(p))
}

// RegisterApplication sets the given application on the map of applications.
func RegisterApplication(applications *map[string]interfaces.IApplication, app interfaces.IApplication) {
	if app == nil || applications == nil {
		panic("Invalid application")
	}
	name := app.Name()
	(*applications)[name] = app
}

// RebuildState calls the RebuildState method for all the given applications.
func RebuildState(p *platform.Platform, applications map[string]interfaces.IApplication) {
	for _, app := range applications {
		if err := app.RebuildState(p.Context); err != nil {
			panic(fmt.Sprintf("Error rebuilding state for application. Name: %s. Version: %s", app.Name(), app.Version()))
		}
	}
}

// RegisterMonitors calls the RegisterMonitors method for all the given applications.
func RegisterMonitors(p *platform.Platform, applications map[string]interfaces.IApplication) {
	for _, app := range applications {
		app.RegisterMonitors(p.Context)
	}
}

// ShutdownApplications calls the Shutdown method for all given applications.
func ShutdownApplications(p *platform.Platform, applications map[string]interfaces.IApplication) {
	for _, app := range applications {
		if err := app.Shutdown(p.Context); err != nil {
			panic(fmt.Sprintf("Error shutting down an application. Name: %s. Version: %s", app.Name(), app.Version()))
		}
	}
}

// RegisterRoutes appends an slice of routes by the given applications to the Platform's router.
func RegisterRoutes(p *platform.Platform, apps map[string]interfaces.IApplication) {
	for _, app := range apps {
		p.Server.Router = router.ConfigureRoutes(p.Server, p.Server.Router, app.Version(), app.Name(), app.RegisterRoutes())
	}
}

// ScheduleTasks gets all the tasks from each application and add them to the platform's scheduler.
func ScheduleTasks(cloudsim *platform.Platform, apps map[string]interfaces.IApplication) {
	for _, app := range apps {
		tasks := app.RegisterTasks()
		for _, task := range tasks {
			cloudsim.Scheduler.DoAt(task.Job, task.Date)
		}
	}
}
