package main

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
)

// Register applications inserts an application to the map of applications in the platform
func RegisterApplications(p *platform.Platform) {
	application.RegisterApplications(p.Applications, subt.Register())
	// p.Applications = application.RegisterApplications(p.Applications, app.Register)
}