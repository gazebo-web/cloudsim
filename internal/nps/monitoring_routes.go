package nps

// This file contains monitoring routes.

import (
	"gitlab.com/ignitionrobotics/web/ign-go"
)

// GetMonitoringRoutes returns the monitoring routes used by this
// application.
// See the definition of the Application interface in application.go
func (app *application) GetMonitoringRoutes() ign.Routes {
	return ign.Routes{

		// This is the Healthz route. This route is used by Flagger during
		// deployment to promote a Green deployment to Blue. Without this route your
		// application will not deploy properly.
		//
		// Example usage:
		//   curl -k -X GET http://localhost:8000/healthz
		ign.Route{
			Name:        "Cloudsim NPS healthcheck",
			Description: "Get cloudsim nps status",
			URI:         "/healthz",
			Headers:     nil,
			Methods: ign.Methods{
				ign.Method{
					Type:        "GET",
					Description: "Get cloudsim nps status",
					Handlers: ign.FormatHandlers{
						ign.FormatHandler{
							Extension: "",
							Handler:   ign.JSONResult(Healthz),
						},
					},
				},
			},
			SecureMethods: ign.SecureMethods{},
		},
	}
}
