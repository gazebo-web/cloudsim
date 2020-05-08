package subt

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/handlers"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

// RegisterRoutes returns a slice of routes for the SubT application.
func (app *SubT) RegisterRoutes() ign.Routes {
	return ign.Routes{
		ign.Route{
			Name:        "Simulations",
			Description: "Information about all simulations",
			URI:         "/simulations",
			Headers:     ign.AuthHeadersRequired,
			Methods:     ign.Methods{},
			SecureMethods: ign.SecureMethods{
				ign.Method{
					Type:        "GET",
					Description: "Get all simulations",
					Handlers: ign.FormatHandlers{
						ign.FormatHandler{
							Extension: ".json",
							Handler:   ign.JSONResultNoTx(handlers.WithUser(app.Services.User, app.Controllers.Simulation.GetAll)),
						},
						ign.FormatHandler{
							Handler: ign.JSONResultNoTx(handlers.WithUser(app.Services.User, app.Controllers.Simulation.GetAll)),
						},
					},
				},
				ign.Method{
					Type:        "POST",
					Description: "Starts a simulation, creating all needed resources",
					Handlers: ign.FormatHandlers{
						ign.FormatHandler{
							Handler: ign.JSONResultNoTx(
								handlers.AfterFn(
									handlers.WithUser(app.Services.User, app.Controllers.Simulation.Start),
									app.Launch,
								),
							),
						},
					},
				},
			},
		},
		ign.Route{
			Name:        "Works with a single simulation based on its groupID",
			Description: "Single simulation based on its groupID",
			URI:         "/simulations/{group}",
			Headers:     ign.AuthHeadersRequired,
			Methods:     ign.Methods{},
			SecureMethods: ign.SecureMethods{
				ign.Method{
					Type:        "GET",
					Description: "Get a single simulation based on its groupID",
					Handlers: ign.FormatHandlers{
						ign.FormatHandler{
							Handler: ign.JSONResultNoTx(
								handlers.WithUser(app.Services.User, app.Controllers.Simulation.Get),
							),
						},
					},
				},
				ign.Method{
					Type:        "DELETE",
					Description: "shutdowns a simulation, removing all associated resources",
					Handlers: ign.FormatHandlers{
						ign.FormatHandler{
							Handler: ign.JSONResultNoTx(
								handlers.AfterFn(
									handlers.WithUser(app.Services.User, app.Controllers.Simulation.Shutdown),
									app.Shutdown,
								),
							),
						},
					},
				},
			},
		},
		ign.Route{
			Name:        "Launches a held simulation based on its groupID",
			Description: "Launches a held simulation based on its groupID",
			URI:         "/simulations/{group}/launch",
			Headers:     ign.AuthHeadersRequired,
			Methods:     ign.Methods{},
			SecureMethods: ign.SecureMethods{
				ign.Method{
					Type:        "POST",
					Description: "Launch a simulation that is being held",
					Handlers: ign.FormatHandlers{
						ign.FormatHandler{
							Handler: ign.JSONResultNoTx(
								handlers.AfterFn(
									handlers.WithUser(app.Services.User, app.Controllers.Simulation.LaunchHeld),
									app.LaunchHeld,
								),
							),
						},
					},
				},
			},
		},
		ign.Route{
			Name:          "",
			Description:   "",
			URI:           "",
			Headers:       nil,
			Methods:       nil,
			SecureMethods: nil,
		},
		ign.Route{
			Name:          "",
			Description:   "",
			URI:           "",
			Headers:       nil,
			Methods:       nil,
			SecureMethods: nil,
		},
	}
}
