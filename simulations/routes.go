package simulations

import (
	"bitbucket.org/ignitionrobotics/ign-go"
)

// Routes declares the routes related to simulations. See also IGN's router.go
var Routes = ign.Routes{

	/////////////////
	// Simulations //
	/////////////////

	// Route for all simulations
	ign.Route{
		"Simulations",
		"Information about all simulations",
		"/simulations",
		ign.AuthHeadersRequired,
		ign.Methods{},
		ign.SecureMethods{
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
				"GET",
				"Get all simulations",
				ign.FormatHandlers{
					ign.FormatHandler{".json", ign.JSONResult(WithUser(CloudsimSimulationList))},
					ign.FormatHandler{"", ign.JSONResult(WithUser(CloudsimSimulationList))},
				},
			},
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
				"POST",
				"Starts a simulation, creating all needed resources",
				ign.FormatHandlers{
					ign.FormatHandler{"", ign.JSONResultNoTx(WithUser(CloudsimSimulationCreate))},
				},
			},
		},
	},

	// Get single simulation route
	ign.Route{
		"Works with a single simulation based on its groupId",
		"Single simulation based on its groupId",
		"/simulations/{group}",
		ign.AuthHeadersRequired,
		ign.Methods{},
		ign.SecureMethods{
			// swagger:route GET /simulations/{group} simulations getSimulation
			//
			// Get a single simulation based on its groupId
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
				"GET",
				"Get a single simulation based on its groupId",
				ign.FormatHandlers{
					ign.FormatHandler{"", ign.JSONResultNoTx(WithUser(GetCloudsimSimulation))},
				},
			},
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
				"DELETE",
				"shutdowns a simulation, removing all associated resources",
				ign.FormatHandlers{
					ign.FormatHandler{"", ign.JSONResultNoTx(WithUser(CloudsimSimulationDelete))},
				},
			},
		},
	},
	// Launch a simulation by the given GroupID that is currently being held by Cloudsim.
	// This route will launch all child simulations for a multisim.
	ign.Route{
		Name:        "Launches a held simulation based on its groupId",
		Description: "Launches a held simulation based on its groupId",
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
		"Restarts a failed simulation based on its groupId",
		"Restarts a failed simulation based on its groupId",
		"/simulations/{group}/restart",
		ign.AuthHeadersRequired,
		ign.Methods{},
		ign.SecureMethods{
			// swagger:route POST /simulations/{group}/restart simulations restartSimulation
			//
			// Restarts a failed simulation based on its groupId
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
				"POST",
				"Restart a failed simulation based on its groupId",
				ign.FormatHandlers{
					ign.FormatHandler{"", ign.JSONResult(WithUser(CloudsimSimulationRestart))},
				},
			},
		},
	},

	// Gateway route to get the URL to live logs or downloadable logs of a simulation.
	ign.Route{
		"Get logs depending on the simulation status",
		"Get file logs or live logs depending on the simulation status",
		"/simulations/{group}/logs",
		ign.AuthHeadersRequired,
		ign.Methods{},
		ign.SecureMethods{
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
				"GET",
				"Get logs from a simulation",
				ign.FormatHandlers{
					ign.FormatHandler{"", ign.JSONResult(WithUser(SimulationLogGateway))},
				},
			},
		},
	},

	// Route to get the live logs of a single simulation
	ign.Route{
		"Live logs of a single simulation",
		"Live logs of a single simulation",
		"/simulations/{group}/logs/live",
		ign.AuthHeadersRequired,
		ign.Methods{},
		ign.SecureMethods{
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
				"GET",
				"Get a live log",
				ign.FormatHandlers{
					ign.FormatHandler{"", ign.JSONResult(WithUser(SimulationLogLive))},
				},
			},
		},
	},

	// Route to get the logs of a single simulation
	ign.Route{
		"Download the logs of a single simulation",
		"Download the logs of a single simulation",
		"/simulations/{group}/logs/file",
		ign.AuthHeadersRequired,
		ign.Methods{},
		ign.SecureMethods{
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
				"GET",
				"Get a log file",
				ign.FormatHandlers{
					ign.FormatHandler{"", ign.JSONResultNoTx(WithUser(SimulationLogFileDownload))},
				},
			},
		},
	},

	// Route to get machine information
	ign.Route{
		"Machines",
		"Information about machines",
		"/machines",
		ign.AuthHeadersRequired,
		ign.Methods{},
		ign.SecureMethods{
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
				"GET",
				"Get all machines",
				ign.FormatHandlers{
					ign.FormatHandler{".json", ign.JSONResult(WithUser(CloudMachineList))},
					ign.FormatHandler{"", ign.JSONResult(WithUser(CloudMachineList))},
				},
			},
		},
	},

	///////////
	// Rules //
	///////////

	// Route to get the remaining number of submissions for an owner in a circuit
	ign.Route{
		"RemainingSubmissions",
		"Returns the number of remaining submissions for an owner in a circuit",
		"/{circuit}/remaining_submissions/{owner}",
		ign.AuthHeadersRequired,
		ign.Methods{},
		ign.SecureMethods{
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
				"GET",
				"Gets the number of remaining submissions in a circuit for an owner ",
				ign.FormatHandlers{
					ign.FormatHandler{"", ign.JSONResultNoTx(WithUser(GetRemainingSubmissions))},
				},
			},
		},
	},

	// Route to get all circuit custom rules
	ign.Route{
		"Rules",
		"Gets the list of all circuit custom rules.",
		"/rules",
		ign.AuthHeadersRequired,
		ign.Methods{},
		ign.SecureMethods{
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
				"GET",
				"Gets the list of all circuit custom rules.",
				ign.FormatHandlers{
					ign.FormatHandler{"", ign.JSONResult(WithUser(CustomRuleList))},
					ign.FormatHandler{".json", ign.JSONResult(WithUser(CustomRuleList))},
				},
			},
		},
	},

	// Route to create/update a custom rule for an owner in a circuit
	ign.Route{
		"SetRule",
		"Creates or updates a custom rule for an owner in a circuit.",
		"/rules/{circuit}/{owner}/{rule}/{value}",
		ign.AuthHeadersRequired,
		ign.Methods{},
		ign.SecureMethods{
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
				"PUT",
				"Creates or updates a custom rule for an owner in a circuit.",
				ign.FormatHandlers{
					ign.FormatHandler{"", ign.JSONResult(WithUser(SetCustomRule))},
				},
			},
		},
	},

	// Route to delete a custom rule for an owner in a circuit
	ign.Route{
		"DeleteRule",
		"Deletes a custom rule for an owner in a circuit.",
		"/rules/{circuit}/{owner}/{rule}",
		ign.AuthHeadersRequired,
		ign.Methods{},
		ign.SecureMethods{
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
				"DELETE",
				"Deletes a custom rule for an owner in a circuit.",
				ign.FormatHandlers{
					ign.FormatHandler{"", ign.JSONResult(WithUser(DeleteCustomRule))},
				},
			},
		},
	},

	// Route to get robots from competition
	ign.Route{
		"Competition robots",
		"Gets the list of robots from the competition",
		"/competition/robots",
		ign.AuthHeadersRequired,
		ign.Methods{},
		ign.SecureMethods{
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
				"GET",
				"Gets the list of robots from the competition",
				ign.FormatHandlers{
					ign.FormatHandler{"", ign.JSONResult(WithUser(GetCompetitionRobots))},
					ign.FormatHandler{".json", ign.JSONResult(WithUser(GetCompetitionRobots))},
				},
			},
		},
	},

	// Extra route to access nodes and their associated hosts
	ign.Route{
		"Remove Nodes and Hosts associated to a group",
		"Deletes cluster nodes and terminates the instances associated to a given Cloudsim GroupId",
		"/k8/nodes",
		ign.AuthHeadersRequired,
		ign.Methods{},
		ign.SecureMethods{
			// swagger:route DELETE /k8/nodes k8 deleteNodesAndHosts
			//
			// Deletes a set of nodes and hosts (instances) associated to a GroupId
			//
			//   Produces:
			//   - application/json
			//
			//   Schemes: https
			//
			//   Responses:
			//     default: Error
			//     200: string
			ign.Method{
				"DELETE",
				"terminates nodes and instances associated to a given groupid",
				ign.FormatHandlers{
					ign.FormatHandler{"", ign.JSONResultNoTx(WithUser(DeleteNodesAndHosts))},
				},
			},
		},
	},

	// Extra route to count the pods in a k8 cluster
	ign.Route{
		"Checks the access to k8 cluster and return the count of running pods",
		"Checks the access to k8 cluster and return the count of running pods",
		"/k8/countpods",
		ign.AuthHeadersRequired,
		ign.Methods{},
		ign.SecureMethods{
			// swagger:route GET /k8/countpods k8 countPods
			//
			// Checks the access to k8 cluster and return the count of running pods.
			// It is used mainly as a test route to ensure the server has correct configuration
			// to access the k8 cluster.
			//
			//   Produces:
			//   - application/json
			//
			//   Schemes: https
			//
			//   Responses:
			//     default: Error
			//     200: string
			ign.Method{
				"GET",
				"Checks the access to k8 cluster and return the count of running pods",
				ign.FormatHandlers{
					ign.FormatHandler{"", ign.JSONResult(WithUser(CountPods))},
				},
			},
		},
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
		URI:         "/queue/{groupIdA}/swap/{groupIdB}",
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
		URI:         "/queue/{groupId}/move/front",
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
		URI:         "/queue/{groupId}/move/back",
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
		URI:         "/queue/{groupId}",
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
}
