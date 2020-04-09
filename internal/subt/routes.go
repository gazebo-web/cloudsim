package subt

import "gitlab.com/ignitionrobotics/web/ign-go"

func (app *SubT) RegisterRoutes() ign.Routes {
	return ign.Routes{
		ign.Route{
			Name:          "Simulations",
			Description:   "Information about all simulations",
			URI:           "/simulations",
			Headers:       ign.AuthHeadersRequired,
			Methods:       ign.Methods{},
			SecureMethods: ign.SecureMethods{
				ign.Method{
					Type:        "GET",
					Description: "Get all simulations",
					Handlers: ign.FormatHandlers{
						ign.FormatHandler{Extension: ".json", Handler: app.Controllers.Simulations.GetAll},
						ign.FormatHandler{Handler: nil},
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
