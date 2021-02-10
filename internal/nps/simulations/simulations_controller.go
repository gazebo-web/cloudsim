package simulations

import (
	"fmt"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/nps/server"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"net/http"
)

type Controller interface {
	server.RoutesGetter
	Start(w http.ResponseWriter, r *http.Request)
	Stop(w http.ResponseWriter, r *http.Request)
}

type controller struct {
	service Service
}

func (c *controller) GetRoutes() ign.Routes {
	return ign.Routes{
		ign.Route{
			Name:        "Start simulation",
			Description: "This is a description for starting a simulation",
			URI:         "/start",
			SecureMethods: []ign.Method{
				{
					Type:        "POST",
					Description: "Start simulations",
					Handlers: ign.FormatHandlers{
						ign.FormatHandler{Handler: http.HandlerFunc(c.Start)},
					},
				},
			},
		},
		ign.Route{
			Name:        "Stop simulation",
			Description: "This is a description for stopping a simulation",
			URI:         "/stop",
			SecureMethods: []ign.Method{
				{
					Type:        "POST",
					Description: "Stop simulations",
					Handlers: ign.FormatHandlers{
						ign.FormatHandler{Handler: http.HandlerFunc(c.Stop)},
					},
				},
			},
		},
	}
}

func (c *controller) Start(w http.ResponseWriter, r *http.Request) {
	// Parse request

	// Get needed data to start simulation from the HTTP request, pass it to the Start Request
	req := StartRequest{}

	res, err := c.service.Start(r.Context(), req)
	if err != nil {
		// Send error message
	}

	// Remove after addressing next comment
	fmt.Println(res)

	// Send response to the user
}

func (c *controller) Stop(w http.ResponseWriter, r *http.Request) {
	// Parse request

	// Get needed data to stop simulation from the HTTP request, pass it to the Stop Request
	req := StopRequest{}

	res, err := c.service.Stop(r.Context(), req)
	if err != nil {
		// Send error message
	}

	// Remove after addressing next comment
	fmt.Println(res)

	// Send response to the user
}

func NewController(service Service) Controller {
	return &controller{}
}
