package nps

// This file defines the API routes.
// Make sure your ~/.kube/config is set appropriately.
//     * For testing use:
//         aws eks update-kubeconfig --name web-cloudsim-testing --kubeconfig=$HOME/.kube/config
//
//
// Debugging commands
// 1. kubectl --kubeconfig=/home/nkoenig/.kube/config -n web-cloudsim-integration get no
// 2. kubectl --kubeconfig=/home/nkoenig/.kube/config -n web-cloudsim-integration get po
//
// Debug deployments
//
// 1. kubectl -n web-cloudsim-nps-staging describe canary
// 2. kubectl -n flagger logs deployment/flagger
// 3. kubectl -n flagger logs deployment/flagger-loadtester
// 4. kubectl -n web-cloudsim-nps-staging edit deploy web-cloudsim-nps-primary
// 5. kubectl -n web-cloudsim-nps-staging edit cm cloudsim-config-nps-primary

import (
	"gitlab.com/ignitionrobotics/web/ign-go"
)

// GetAPIRoutes returns the routes used by this application.
// See the definition of the Application interface in application.go
func (app *application) GetAPIRoutes() ign.Routes {
	ctrl := app.controller

	// Return the routes for this application. See also IGN's router.go
	return ign.Routes{
		// Example usage:
		//     curl -X POST -H "Private-Token: TOKEN" http://localhost:8000/1.0/start -F "image=osrf/ros:melodic-desktop-full" -F "args=gazebo" -F "name=my_test_name"
    ign.Route{
      Name:        "Start simulation",
      Description: "This is a route for starting a simulation",
      URI:         "/start",
      Methods: ign.Methods{},
      SecureMethods: ign.SecureMethods{
        ign.Method{
            Type:        "POST",
            Description: "Start simulations",
            Handlers: ign.FormatHandlers{
              ign.FormatHandler{Handler: ign.JSONResult(WithUser(ctrl.Start))},
            },
        },
      },
    },
		// Example usage:
		//     curl -X POST -H "Private-Token: TOKEN" http://localhost:8000/1.0/stop/{groupid}
		ign.Route{
			Name:        "Stop simulation",
			Description: "This is a route for stopping a simulation",
			URI:         "/stop/{groupid}",
      Methods: ign.Methods{},
			SecureMethods: ign.SecureMethods{
				ign.Method{
					Type:        "POST",
					Description: "Stop simulations",
					Handlers: ign.FormatHandlers{
						ign.FormatHandler{Handler: ign.JSONResult(WithUser(ctrl.Stop))},
					},
				},
			},
		},
		// Example usage:
		//     curl -X GET -H "Private-Token: TOKEN" http://localhost:8000/1.0/simulations
		ign.Route{
			Name:        "List simulations",
			Description: "This is a route for listing simulations",
			URI:         "/simulations",
      Methods: ign.Methods{},
			SecureMethods: ign.SecureMethods{
				ign.Method{
					Type:        "GET",
					Description: "List simulations",
					Handlers: ign.FormatHandlers{
						ign.FormatHandler{Handler: ign.JSONResult(WithUser(ctrl.ListSimulations))},
					},
				},
			},
		},
		// Example usage:
    //     curl -X POST -H "Private-Token: TOKEN" http://localhost:8000/1.0/simulations/{groupid}
		ign.Route{
			Name:        "Get simulation",
			Description: "This is a route for acquiring information about a simulation",
			URI:         "/simulations/{groupid}",
      Methods: ign.Methods{},
			SecureMethods: ign.SecureMethods{
				ign.Method{
					Type:        "GET",
					Description: "Get information about a simulation",
					Handlers: ign.FormatHandlers{
						ign.FormatHandler{Handler: ign.JSONResult(WithUser(ctrl.GetSimulation))},
					},
				},
			},
		},
	}
}
