package nps

// This file defines the API routes.

import (
	"gitlab.com/ignitionrobotics/web/ign-go"
	"net/http"
)

// GetAPIRoutes returns the routes used by this application.
// See the definition of the Application interface in application.go
func (app *application) GetAPIRoutes() ign.Routes {
	ctrl := app.controller

	// Return the routes for this application. See also IGN's router.go
	return ign.Routes{
    // Example usage:
    //     curl -X POST http://localhost:8000/1.0/start -F "image=DOCKER_IMAGE"
		ign.Route{
			Name:        "Start simulation",
			Description: "This is a description for starting a simulation",
			URI:         "/start",
			Methods: []ign.Method{
				{
					Type:        "POST",
					Description: "Start simulations",
					Handlers: ign.FormatHandlers{
						ign.FormatHandler{Handler: http.HandlerFunc(ctrl.Start)},
					},
				},
			},
		},
    // Example usage:
    //     curl -X POST http://localhost:8000/1.0/stop
		ign.Route{
			Name:        "Stop simulation",
			Description: "This is a description for stopping a simulation",
			URI:         "/stop",
			Methods: []ign.Method{
				{
					Type:        "POST",
					Description: "Stop simulations",
					Handlers: ign.FormatHandlers{
						ign.FormatHandler{Handler: http.HandlerFunc(ctrl.Stop)},
					},
				},
			},
		},
	}
}
