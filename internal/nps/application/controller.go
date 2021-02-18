package nps

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"net/http"
)

// Controller is an interface designed to handle route requests.
type Controller interface {
	Start(w http.ResponseWriter, r *http.Request)
	Stop(w http.ResponseWriter, r *http.Request)
}

type controller struct {
	// service is this controller's implementation of the
	// cloudsim/pkg/simulations service. See the simulations_service.go file.
	service Service
}

// NewController creates a new controller
func NewController(db *gorm.DB, logger ign.Logger) Controller {
	return &controller{
		// Create a simulation service to manage simulation instances
		service: NewService(db, logger),
	}
}

// Start handles the `/start` route.
func (ctrl *controller) Start(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\n\nHERE\n\n")
	// Parse request

	// Get needed data to start simulation from the HTTP request, pass it to the Start Request
	req := StartRequest{}

	res, err := ctrl.service.Start(r.Context(), req)
	if err != nil {
		// Send error message
	}

	// Remove after addressing next comment
	fmt.Println(res)

	// Send response to the user
}

// Stop handles the `/stop` route.
func (ctrl *controller) Stop(w http.ResponseWriter, r *http.Request) {
	// Parse request

	// Get needed data to stop simulation from the HTTP request, pass it to the Stop Request
	req := StopRequest{}

	res, err := ctrl.service.Stop(r.Context(), req)
	if err != nil {
		// Send error message
	}

	// Remove after addressing next comment
	fmt.Println(res)

	// Send response to the user
}
