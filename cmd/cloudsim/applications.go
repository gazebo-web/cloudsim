package main

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
)

// Register applications inserts an application to the map of applications in the platform
func RegisterApplications(p *platform.Platform, apps *map[string]application.IApplication) {
	RegisterApplication(apps, subt.Register(p))
	// p.Applications = application.RegisterApplication(p.Applications, app.Register)
}

// RegisterApplications adds a given application to the platform.
// Returns the list of applications
func RegisterApplication(applications *map[string]application.IApplication, app application.IApplication) {
	if app == nil || applications == nil {
		panic("Invalid application")
	}
	name := app.Name()
	(*applications)[name] = app
}
