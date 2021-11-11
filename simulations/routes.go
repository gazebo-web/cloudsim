package simulations

import (
	"gitlab.com/ignitionrobotics/web/ign-go"
	"net/http"
	"net/http/pprof"
)

// Routes declares the routes related to simulations. See also IGN's router.go
var Routes = ign.Routes{

	/////////////////
	// Simulations //
	/////////////////

	// Route for all simulations
	ign.Route{
		Name:        "Simulations",
		Description: "Information about all simulations",
		URI:         "/simulations",
		Headers:     ign.AuthHeadersRequired,
		Methods: ign.Methods{
			// swagger:route GET /simulations simulations listSimulations
			//
			// Get list of simulations.
			//
			// Get the list of simulations. Simulations will be returned paginated,
			// with pages of 20 items by default. The user can request a
			// different page with query parameter 'page', and the page size
			// can be defined with query parameter 'per_page'.
			// The route supports the 'order' parameter, with values 'asc' and
			// 'desc' (default: desc).
			//
			// This route also supports getting simulations of a particular
			// organization with the 'owner' parameter
			//
			// This route also supports filtering by privacy with the 'private'
			// parameter. By default the route returns simulations based on the
			// permissions of the JWT user, if present. Otherwise, it returns public
			// simulations
			//
			// This route also supports the 'status' parameter
			// which filters the results based on a status string with one of the
			// following options: ["Pending", "LaunchingNodes", "LaunchingPods",
			// "Running", "ToBeTerminated", "DeletingPods", "DeletingNodes",
			// "TerminatingInstances", "Terminated"]. Prefixing the status string with
			// and exclamation mark, !, will invert the filter logic.
			//
			// This route also supports the 'errorStatus' parameter
			// which filters the results based on an error status string with one of
			// the following options: ["InitializationFailed", "TerminationFailed",
			// "AdminReview"]. Prefixing the error status string with
			// and exclamation mark, !, will invert the filter logic.
			//
			//   Produces:
			//   - application/json
			//
			//   Schemes: https
			//
			//   Responses:
			//     default: Error
			//     200: jsonSims
			ign.Method{
				Type:        "GET",
				Description: "Get all simulations",
				Handlers: ign.FormatHandlers{
					ign.FormatHandler{Extension: ".json", Handler: ign.JSONResult(WithUserOrAnonymous(CloudsimSimulationList))},
					ign.FormatHandler{Handler: ign.JSONResult(WithUserOrAnonymous(CloudsimSimulationList))},
				},
			},
		},
		SecureMethods: ign.SecureMethods{
			// swagger:route POST /simulations simulations createSimulation
			//
			// Launches a new cloudsim simulation
			//
			//   Produces:
			//   - application/json
			//
			//   Schemes: https
			//
			//   Responses:
			//     default: Error
			//     200: jsonSim
			ign.Method{
				Type:        "POST",
				Description: "Starts a simulation, creating all needed resources",
				Handlers: ign.FormatHandlers{
					ign.FormatHandler{Handler: ign.JSONResultNoTx(WithUser(CloudsimSimulationCreate))},
				},
			},
		},
	},

	// Get single simulation route
	ign.Route{
		Name:        "Works with a single simulation based on its groupID",
		Description: "Single simulation based on its groupID",
		URI:         "/simulations/{group}",
		Headers:     ign.AuthHeadersRequired,
		Methods: ign.Methods{
			// swagger:route GET /simulations/{group} simulations getSimulation
			//
			// Get a single simulation based on its groupID
			//
			//   Produces:
			//   - application/json
			//
			//   Schemes: https
			//
			//   Responses:
			//     default: Error
			//     200: jsonSim
			ign.Method{
				Type:        "GET",
				Description: "Get a single simulation based on its groupID",
				Handlers: ign.FormatHandlers{
					ign.FormatHandler{Handler: ign.JSONResultNoTx(WithUserOrAnonymous(GetCloudsimSimulation))},
				},
			},
		},
		SecureMethods: ign.SecureMethods{
			// swagger:route DELETE /simulations/{group} simulations deleteSimulation
			//
			// Deletes a cloudsim simulation
			//
			//   Produces:
			//   - application/json
			//
			//   Schemes: https
			//
			//   Responses:
			//     default: Error
			//     200: jsonSim
			ign.Method{
				Type:        "DELETE",
				Description: "shutdowns a simulation, removing all associated resources",
				Handlers: ign.FormatHandlers{
					ign.FormatHandler{Handler: ign.JSONResultNoTx(WithUser(CloudsimSimulationDelete))},
				},
			},
		},
	},

	// Launch a simulation by the given GroupID that is currently being held by Cloudsim.
	// This route will launch all child simulations for a multisim.
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
						Extension: "",
						Handler:   ign.JSONResult(WithUser(CloudsimSimulationLaunch))},
				},
			},
		},
	},

	// Restart a simulation based on its Group ID
	ign.Route{
		Name:        "Restarts a failed simulation based on its groupID",
		Description: "Restarts a failed simulation based on its groupID",
		URI:         "/simulations/{group}/restart",
		Headers:     ign.AuthHeadersRequired,
		Methods:     ign.Methods{},
		SecureMethods: ign.SecureMethods{
			// swagger:route POST /simulations/{group}/restart simulations restartSimulation
			//
			// Restarts a failed simulation based on its groupID
			//
			//   Produces:
			//   - application/json
			//
			//   Schemes: https
			//
			//   Responses:
			//     default: Error
			//     200: jsonSim
			ign.Method{
				Type:        "POST",
				Description: "Restart a failed simulation based on its groupID",
				Handlers: ign.FormatHandlers{
					ign.FormatHandler{Handler: ign.JSONResult(WithUser(CloudsimSimulationRestart))},
				},
			},
		},
	},

	// Gateway route to get the URL to live logs or downloadable logs of a simulation.
	ign.Route{
		Name:        "Get logs depending on the simulation status",
		Description: "Get file logs or live logs depending on the simulation status",
		URI:         "/simulations/{group}/logs",
		Headers:     ign.AuthHeadersRequired,
		Methods:     ign.Methods{},
		SecureMethods: ign.SecureMethods{
			// swagger:route GET /simulations/{group}/logs/ simulations getLogsGateway
			//
			// Get the current log depending on the simulation status
			//
			// Get the current log depending on the simulation status
			//
			//   Schemes: https
			//
			//   Responses:
			//     default: Error
			//     200: application/json
			ign.Method{
				Type:        "GET",
				Description: "Get logs from a simulation",
				Handlers: ign.FormatHandlers{
					ign.FormatHandler{
						Extension: "",
						Handler:   ign.JSONResult(WithUser(SimulationLogGateway)),
					},
				},
			},
		},
	},

	// Route to get the live logs of a single simulation
	ign.Route{
		Name:        "Live logs of a single simulation",
		Description: "Live logs of a single simulation",
		URI:         "/simulations/{group}/logs/live",
		Headers:     ign.AuthHeadersRequired,
		Methods:     ign.Methods{},
		SecureMethods: ign.SecureMethods{
			// swagger:route GET /simulations/{group}/logs/live/ simulations getLogFile
			//
			// Get live logs from a running simulation
			//
			// Get live logs from a running simulation
			//
			//   Schemes: https
			//
			//   Responses:
			//     default: Error
			//     200: text/plain
			ign.Method{
				Type:        "GET",
				Description: "Get a live log",
				Handlers: ign.FormatHandlers{
					ign.FormatHandler{Handler: ign.JSONResult(WithUser(SimulationLogLive))},
				},
			},
		},
	},

	// Route to get the logs of a single simulation
	ign.Route{
		Name:        "Download the logs of a single simulation",
		Description: "Download the logs of a single simulation",
		URI:         "/simulations/{group}/logs/file",
		Headers:     ign.AuthHeadersRequired,
		Methods:     ign.Methods{},
		SecureMethods: ign.SecureMethods{
			// swagger:route GET /simulations/{group}/logs/file simulations downloadLogFile
			//
			// Download simulation's log files
			//
			// Download simulation's log files
			//
			//   Schemes: https
			//
			//   Responses:
			//     default: Error
			//     200: File
			ign.Method{
				Type:        "GET",
				Description: "Get a log file",
				Handlers: ign.FormatHandlers{
					ign.FormatHandler{Handler: ign.JSONResultNoTx(WithUser(SimulationLogFileDownload))},
				},
			},
		},
	},

	// Route to get the address of a simulation's websocket server and authorization token
	ign.Route{
		Name:        "Gets a simulation's websocket server address",
		Description: "Gets a simulation's websocket server address",
		URI:         "/simulations/{group}/websocket",
		Headers:     ign.AuthHeadersRequired,
		Methods: ign.Methods{
			// swagger:route GET /simulations/{group}/websocket simulations websocket
			//
			// Gets a simulation's websocket server address and authorization token
			//
			// Gets a simulation's websocket server address and authorization token
			//
			//   Schemes: https
			//
			//   Responses:
			//     200: json
			//     503: Error
			ign.Method{
				Type:        "GET",
				Description: "Get a simulation's websocket server address and authorization token",
				Handlers: ign.FormatHandlers{
					ign.FormatHandler{Handler: ign.JSONResultNoTx(WithUserOrAnonymous(SimulationWebsocketAddress))},
				},
			},
		},
		SecureMethods: ign.SecureMethods{},
	},

	ign.Route{
		Name:        "Reconnect websocket",
		Description: "Allow admins to reconnect a specific simulation to its respective websocket server",
		URI:         "/simulations/{group}/reconnect",
		Headers:     ign.AuthHeadersRequired,
		Methods:     ign.Methods{},
		SecureMethods: ign.SecureMethods{
			ign.Method{
				Type:        "POST",
				Description: "Reconnect websocket",
				Handlers: ign.FormatHandlers{
					ign.FormatHandler{
						Extension: "",
						Handler:   ign.JSONResult(WithUser(ReconnectWebsocket)),
					},
				},
			},
		},
	},

	// Route to get machine information
	ign.Route{
		Name:        "Machines",
		Description: "Information about machines",
		URI:         "/machines",
		Headers:     ign.AuthHeadersRequired,
		Methods:     ign.Methods{},
		SecureMethods: ign.SecureMethods{
			// swagger:route GET /machines simulations listMachines
			//
			// Get list of cloud machines.
			//
			// Get the list of machines. Machines will be returned paginated,
			// with pages of 20 items by default. The user can request a
			// different page with query parameter 'page', and the page size
			// can be defined with query parameter 'per_page'.
			// The route supports the 'order' parameter, with values 'asc' and
			// 'desc' (default: desc).
			//
			// This route also supports the 'status' parameter
			// which filters the results based on a status string with one of the
			// following options: ["initializing", "running", "terminating",
			// "terminated", "error"]. Prefixing the status string with
			// and exclamation mark, !, will invert the filter logic.
			//
			//   Produces:
			//   - application/json
			//
			//   Schemes: https
			//
			//   Responses:
			//     default: Error
			//     200: jsonMachines
			ign.Method{
				Type:        "GET",
				Description: "Get all machines",
				Handlers: ign.FormatHandlers{
					ign.FormatHandler{Extension: ".json", Handler: ign.JSONResult(WithUser(CloudMachineList))},
					ign.FormatHandler{Handler: ign.JSONResult(WithUser(CloudMachineList))},
				},
			},
		},
	},

	///////////
	// Rules //
	///////////

	// Route to get the remaining number of submissions for an owner in a circuit
	ign.Route{
		Name:        "RemainingSubmissions",
		Description: "Returns the number of remaining submissions for an owner in a circuit",
		URI:         "/{circuit}/remaining_submissions/{owner}",
		Headers:     ign.AuthHeadersRequired,
		Methods:     ign.Methods{},
		SecureMethods: ign.SecureMethods{
			// swagger:route GET /{circuit}/remaining_submissions/{owner} submissions getSubmissions
			//
			// Returns the number of remaining submissions for an owner in
			// the specified circuit.
			//
			//   Produces:
			//   - application/json
			//
			//   Schemes: https
			//
			//   Responses:
			//     default: Error
			//     200: jsonSim
			ign.Method{
				Type:        "GET",
				Description: "Gets the number of remaining submissions in a circuit for an owner ",
				Handlers: ign.FormatHandlers{
					ign.FormatHandler{Handler: ign.JSONResultNoTx(WithUser(GetRemainingSubmissions))},
				},
			},
		},
	},

	// Route to get all circuit custom rules
	ign.Route{
		Name:        "Rules",
		Description: "Gets the list of all circuit custom rules.",
		URI:         "/rules",
		Headers:     ign.AuthHeadersRequired,
		Methods:     ign.Methods{},
		SecureMethods: ign.SecureMethods{
			// swagger:route GET /rules rules listRules
			//
			// Gets the list of all circuit custom rules. Rules will be returned
			// paginated, with pages of 20 items by default. The user can request a
			// different page with query parameter 'page', and the page size
			// can be defined with query parameter 'per_page'.
			// The route supports the 'order' parameter, with values 'asc' and
			// 'desc' (default: desc).
			//
			// This route also supports the `circuit`, `owner` and `rule_type` parameters,
			// which filter the results.
			//
			//   Produces:
			//   - application/json
			//
			//   Schemes: https
			//
			//   Responses:
			//     default: Error
			//     200: jsonSim
			ign.Method{
				Type:        "GET",
				Description: "Gets the list of all circuit custom rules.",
				Handlers: ign.FormatHandlers{
					ign.FormatHandler{Handler: ign.JSONResult(WithUser(CustomRuleList))},
					ign.FormatHandler{Extension: ".json", Handler: ign.JSONResult(WithUser(CustomRuleList))},
				},
			},
		},
	},

	// Route to create/update a custom rule for an owner in a circuit
	ign.Route{
		Name:        "SetRule",
		Description: "Creates or updates a custom rule for an owner in a circuit.",
		URI:         "/rules/{circuit}/{owner}/{rule}/{value}",
		Headers:     ign.AuthHeadersRequired,
		Methods:     ign.Methods{},
		SecureMethods: ign.SecureMethods{
			// swagger:route PUT /rules/{circuit}/{owner}/{rule}/{value} rules setRules
			//
			// Creates or updates a custom rule for an owner in a circuit.
			//
			//   Produces:
			//   - application/json
			//
			//   Schemes: https
			//
			//   Responses:
			//     default: Error
			//     200: jsonSim
			ign.Method{
				Type:        "PUT",
				Description: "Creates or updates a custom rule for an owner in a circuit.",
				Handlers: ign.FormatHandlers{
					ign.FormatHandler{Handler: ign.JSONResult(WithUser(SetCustomRule))},
				},
			},
		},
	},

	// Route to delete a custom rule for an owner in a circuit
	ign.Route{
		Name:        "DeleteRule",
		Description: "Deletes a custom rule for an owner in a circuit.",
		URI:         "/rules/{circuit}/{owner}/{rule}",
		Headers:     ign.AuthHeadersRequired,
		Methods:     ign.Methods{},
		SecureMethods: ign.SecureMethods{
			// swagger:route DELETE /rules/{circuit}/{owner}/{rule} rules setRules
			//
			// Creates or updates a custom rule for an owner in a circuit.
			//
			//   Schemes: https
			//
			//   Responses:
			//     default: Error
			//     200: jsonSim
			ign.Method{
				Type:        "DELETE",
				Description: "Deletes a custom rule for an owner in a circuit.",
				Handlers: ign.FormatHandlers{
					ign.FormatHandler{Handler: ign.JSONResult(WithUser(DeleteCustomRule))},
				},
			},
		},
	},

	// Route to get robots from competition
	ign.Route{
		Name:        "Competition robots",
		Description: "Gets the list of robots from the competition",
		URI:         "/competition/robots",
		Headers:     ign.AuthHeadersRequired,
		Methods: ign.Methods{
			// swagger:route GET /competition/robots competition robots
			//
			// Gets the list of all competition robots.
			//
			//   Produces:
			//   - application/json
			//
			//   Schemes: https
			//
			//   Responses:
			//     default: Error
			//     200: Robots
			ign.Method{
				Type:        "GET",
				Description: "Gets the list of robots from the competition",
				Handlers: ign.FormatHandlers{
					ign.FormatHandler{Handler: ign.JSONResult(WithUserOrAnonymous(GetCompetitionRobots))},
					ign.FormatHandler{Extension: ".json", Handler: ign.JSONResult(WithUserOrAnonymous(GetCompetitionRobots))},
				},
			},
		},
		SecureMethods: ign.SecureMethods{},
	},

	//////////////
	//	Queue	//
	//////////////

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
						Handler:   ign.JSONResult(WithUser(QueueGet)),
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
						Handler:   ign.JSONResult(WithUser(QueueCount)),
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
						Handler:   ign.JSONResult(WithUser(QueueSwap)),
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
						Handler:   ign.JSONResult(WithUser(QueueMoveToFront)),
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
						Handler:   ign.JSONResult(WithUser(QueueMoveToBack)),
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
						Handler:   ign.JSONResult(WithUser(QueueRemove)),
					},
				},
			},
		},
	},

	//////////////
	// Billing	//
	//////////////

	ign.Route{
		Name:        "Get credits balance",
		Description: "Get credits balance of the current user",
		URI:         "/billing/credits",
		Headers:     ign.AuthHeadersRequired,
		Methods:     ign.Methods{},
		SecureMethods: ign.SecureMethods{
			ign.Method{
				Type:        "GET",
				Description: "Get credits balance",
				Handlers: ign.FormatHandlers{
					ign.FormatHandler{
						Extension: "",
						Handler:   ign.JSONResult(WithUser(GetCreditsBalance)),
					},
				},
			},
		},
	},

	ign.Route{
		Name:        "Create payment session",
		Description: "Start a payment session",
		URI:         "/billing/session",
		Headers:     ign.AuthHeadersRequired,
		Methods:     ign.Methods{},
		SecureMethods: ign.SecureMethods{
			ign.Method{
				Type:        "POST",
				Description: "Create a payment session",
				Handlers: ign.FormatHandlers{
					ign.FormatHandler{
						Extension: "",
						Handler:   ign.JSONResult(WithUser(CreateSession)),
					},
				},
			},
		},
	},

	/////////////////////
	ign.Route{
		Name:        "Debug",
		Description: "Debug multi region support",
		URI:         "/debug/{groupID}",
		Headers:     ign.AuthHeadersRequired,
		Methods:     ign.Methods{},
		SecureMethods: ign.SecureMethods{
			ign.Method{
				Type:        "GET",
				Description: "Debug websocket messages",
				Handlers: ign.FormatHandlers{
					ign.FormatHandler{
						Extension: "",
						Handler:   ign.JSONResult(WithUser(Debug)),
					},
				},
			},
		},
	},
}

// MonitoringRoutes contains the different routes used for service monitoring.
var MonitoringRoutes = ign.Routes{

	///////////////
	//  Healthz  //
	///////////////

	ign.Route{
		Name:        "Cloudsim healthcheck",
		Description: "Get cloudsim status",
		URI:         "/healthz",
		Headers:     nil,
		Methods: ign.Methods{
			ign.Method{
				Type:        "GET",
				Description: "Get cloudsim status",
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

// ProfileRoutes contains the different routes to perform CPU profiling.
var ProfileRoutes = ign.Routes{
	ign.Route{
		Name:        "CPU Profile",
		Description: "Get cloudsim CPU profile data",
		URI:         "/profile",
		SecureMethods: ign.SecureMethods{
			ign.Method{
				Type:        "GET",
				Description: "Get cloudsim CPU profile data",
				Handlers: ign.FormatHandlers{
					ign.FormatHandler{
						Extension: "",
						Handler:   http.HandlerFunc(pprof.Profile),
					},
				},
			},
		},
	},
}
