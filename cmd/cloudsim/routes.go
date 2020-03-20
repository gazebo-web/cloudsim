package main

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/router"
)

func RegisterRoutes(p platform.Platform) {
	router.ConfigureRoutes(p.Server, p.Server.Router, "1.0", "subt", subt.RegisterRoutes())
	// router.ConfigureRoutes(p.Server, p.Server.Router, "2.0", "app", app.RegisterRoutes())
}