package platform

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/handlers"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

func (p *platform) RegisterRoutes() ign.Routes {
	return ign.Routes{
		// Launch queue - Get all elements
		ign.Route{
			Name:        "Get all elements from queue",
			Description: "Get all elements from queue. This route should optionally be able to handle pagination parameters.",
			URI:         "/queue",
			Headers:     ign.AuthHeadersRequired,
			Methods:     ign.Methods{},
			SecureMethods: ign.SecureMethods{
				ign.Method{
					Type:        "GET",
					Description: "Get all elements from queue. This route should optionally be able to handle pagination parameters",
					Handlers: ign.FormatHandlers{
						ign.FormatHandler{
							Extension: "",
							Handler:   ign.JSONResultNoTx(handlers.WithUser(p.Services().User(), p.Controllers().Queue().GetAll)),
						},
					},
				},
			},
		},
		// Launch queue - Count elements
		ign.Route{
			Name:        "Count elements in the queue",
			Description: "Get the amount of elements in the queue",
			URI:         "/queue/count",
			Headers:     ign.AuthHeadersRequired,
			Methods:     ign.Methods{},
			SecureMethods: ign.SecureMethods{
				ign.Method{
					Type:        "GET",
					Description: "Get the amount of elements in the queue",
					Handlers: ign.FormatHandlers{
						ign.FormatHandler{
							Extension: "",
							Handler:   ign.JSONResultNoTx(handlers.WithUser(p.Services().User(), p.Controllers().Queue().Count)),
						},
					},
				},
			},
		},
		// Launch queue - Swap elements
		ign.Route{
			Name:        "Swap queue elements moving A to B and vice versa",
			Description: "Swap queue elements moving A to B and vice versa",
			URI:         "/queue/{groupIDA}/swap/{groupIDB}",
			Headers:     ign.AuthHeadersRequired,
			Methods:     ign.Methods{},
			SecureMethods: ign.SecureMethods{
				ign.Method{
					Type:        "PATCH",
					Description: "Swap queue elements moving A to B and vice versa",
					Handlers: ign.FormatHandlers{
						ign.FormatHandler{
							Extension: "",
							Handler:   ign.JSONResultNoTx(handlers.WithUser(p.Services().User(), p.Controllers().Queue().Swap)),
						},
					},
				},
			},
		},
		// Launch queue - Move to front
		ign.Route{
			Name:        "Move an element to the front of the queue",
			Description: "Move an element to the front of the queue",
			URI:         "/queue/{groupID}/move/front",
			Headers:     ign.AuthHeadersRequired,
			Methods:     ign.Methods{},
			SecureMethods: ign.SecureMethods{
				ign.Method{
					Type:        "PATCH",
					Description: "Move an element to the front of the queue",
					Handlers: ign.FormatHandlers{
						ign.FormatHandler{
							Extension: "",
							Handler:   ign.JSONResultNoTx(handlers.WithUser(p.Services().User(), p.Controllers().Queue().MoveToFront)),
						},
					},
				},
			},
		},
		// Launch queue - Move to back
		ign.Route{
			Name:        "Move an element to the back of the queue",
			Description: "Move an element to the back of the queue",
			URI:         "/queue/{groupID}/move/back",
			Headers:     ign.AuthHeadersRequired,
			Methods:     ign.Methods{},
			SecureMethods: ign.SecureMethods{
				ign.Method{
					Type:        "PATCH",
					Description: "Move an element to the back of the queue",
					Handlers: ign.FormatHandlers{
						ign.FormatHandler{
							Extension: "",
							Handler:   ign.JSONResultNoTx(handlers.WithUser(p.Services().User(), p.Controllers().Queue().MoveToBack)),
						},
					},
				},
			},
		},
		// Launch queue - Remove an element
		ign.Route{
			Name:        "Remove an element from the queue",
			Description: "Remove an element from the queue",
			URI:         "/queue/{groupID}",
			Headers:     ign.AuthHeadersRequired,
			Methods:     ign.Methods{},
			SecureMethods: ign.SecureMethods{
				ign.Method{
					Type:        "DELETE",
					Description: "Remove an element from the queue",
					Handlers: ign.FormatHandlers{
						ign.FormatHandler{
							Extension: "",
							Handler:   ign.JSONResultNoTx(handlers.WithUser(p.Services().User(), p.Controllers().Queue().Remove)),
						},
					},
				},
			},
		},
	}
}
