package main

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/router"
)

// RegisterRoutes appends an slice of routes by the given applications to the Platform's router.
func RegisterRoutes(p *platform.Platform, apps map[string]application.IApplication) {
	for _, app := range apps {
		p.Server.Router = router.ConfigureRoutes(p.Server, p.Server.Router, app.Version(), app.Name(), app.RegisterRoutes())
	}
}