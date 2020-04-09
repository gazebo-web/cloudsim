package main

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
)

// RegisterApplications registers every application by calling their Register method.
// The platform is passed to the
func RegisterApplications(p *platform.Platform, apps *map[string]application.IApplication) {
	RegisterApplication(apps, subt.Register(p))
	// p.Applications = application.RegisterApplication(p.Applications, app.Register)
}

// RegisterApplication adds a given application to the map of applications.
func RegisterApplication(applications *map[string]application.IApplication, app application.IApplication) {
	if app == nil || applications == nil {
		panic("Invalid application")
	}
	name := app.Name()
	(*applications)[name] = app
}
