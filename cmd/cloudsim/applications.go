package main

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
)

// Register applications inserts an application to the map of applications in the platform
func RegisterApplications(p *platform.Platform, apps *map[string]application.IApplication) {
	application.RegisterApplication(apps, subt.Register(p))
	// p.Applications = application.RegisterApplication(p.Applications, app.Register)
}