package nps

// This file contains monitoring routes. 

import (
	"gitlab.com/ignitionrobotics/web/ign-go"
)

// GetMonitoringRoutes returns the monitoring routes used by this
// application. 
// See the definition of the Application interface in application.go
func (app *application) GetMonitoringRoutes() ign.Routes {
  return ign.Routes{}
}
