package nps

// This file defines the controller, which handles route requests. An
// application creates an instance of a controller by calling `NewController`.

import (
	"context"
	"errors"
	"fmt"
	"gitlab.com/ignitionrobotics/web/fuelserver/bundles/users"
	useracc "gitlab.com/ignitionrobotics/web/cloudsim/pkg/users"
	"github.com/go-playground/form"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"net/http"
)

// Controller is an interface designed to handle route requests.
type Controller interface {
	Start(user *users.User, tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg)
	Stop(user *users.User, tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg)
	ListSimulations(user *users.User, tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg)
	GetSimulation(user *users.User, tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg)
}

type controller struct {
	// service is this controller's implementation of the
	// cloudsim/pkg/simulations service. See the simulations_service.go file.
	service     Service
	formDecoder *form.Decoder
}

// NewController creates a new controller
func NewController(db *gorm.DB, logger ign.Logger) Controller {
	return &controller{
		// Create a simulation service to manage simulation instances
		service:     NewService(db, logger),
		formDecoder: form.NewDecoder(),
	}
}

// Builds the ErrMsg extra info from the given DecodeErrors
// \todo: Move this to a common place so that we don't have to copy it
// every time we create an application.
func getDecodeErrorsExtraInfo(err error) []string {
	errs := err.(form.DecodeErrors)
	extra := make([]string, 0, len(errs))
	for field, er := range errs {
		extra = append(extra, fmt.Sprintf("Field: %s. %v", field, er.Error()))
	}
	return extra
}

// Start handles the `/start` route.
//
// Origin: user --> POST /start --> controller.Start()
// Next:
//     * On success --> service.Start
//     * On fail --> return error
func (ctrl *controller) Start(user *users.User, tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {

	// Parse form's values and files.
	if err := r.ParseMultipartForm(0); err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorForm, err)
	}
	defer r.MultipartForm.RemoveAll()

	// Get needed data to start simulation from the HTTP request, pass it to the
	// Start Request
	var req StartRequest

	if errs := ctrl.formDecoder.Decode(&req, r.Form); errs != nil {
		return nil, ign.NewErrorMessageWithArgs(ign.ErrorFormInvalidValue, errs,
			getDecodeErrorsExtraInfo(errs))
	}

	// A name form field is required. This is the name of the pod.
	if req.Name == "" {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorMissingField,
			errors.New("Missing 'name' form field"))
	}

	// An image form field is required. This is the docker image to run.
	if req.Image == "" {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorMissingField,
			errors.New("Missing 'image' form field"))
	}

	// Make sure the some arguments are set. The arguments are passed to the
	// docker image.
	if req.Args == "" {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorMissingField,
			errors.New("Missing 'args' form field"))
	}

	// Hand off the start request data to the service.
	res, err := ctrl.service.Start(user, tx, r.Context(), req)
	if err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorForm, err)
	}

	// Send response to the user
	return res, nil
}

// Stop handles the `/stop` route.
//
// Origin: user --> POST /start --> controller.Start()
// Next:
//     * On success --> service.Start
//     * On fail --> return error
func (ctrl *controller) Stop(user *users.User, tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
  // Get the groupid from the route
	groupID, ok := mux.Vars(r)["groupid"]
	if !ok {
		return nil, ign.NewErrorMessage(ign.ErrorIDNotInRequest)
	}

  // Get the matching simulation
	var simulation Simulation
	if err := tx.Where("group_id=?", groupID).First(&simulation).Error; err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorIDNotFound, err)
	}

	// Construct the stop request to send to the service
	req := StopRequest{
    GroupID: simulation.GroupID,
  }

	res, err := ctrl.service.Stop(user, tx, r.Context(), req)
	if err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorForm, err)
	}

	// Send response to the user
	return res, nil
}

// ListSimulations handles the `/simulations` route.
//
// Origin: user --> GET /simulations --> controller.ListSimulations()
// Next:
//     * On success --> return ListResponse
//     * On fail --> return error
func (ctrl *controller) ListSimulations(user *users.User, tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {

	var simulations Simulations
	tx.Where("owner = ?", user.Username).Find(&simulations)

	var response ListResponse
	for _, sim := range simulations {
		response.Simulations = append(response.Simulations, GetSimulationResponse{
			Name:    sim.Name,
			GroupID: sim.GroupID,
			Status:  sim.Status,
			Image:   sim.Image,
			Args:    sim.Args,
			URI:     sim.URI,
			IP:      sim.IP,
		})
	}

	// Send the group id to the queue
	return &response, nil

	// Send response to the user
	// return res, nil
}

// GetSimulation handles the `/simulation/{id}` route.
//
// Origin: user --> GET /simulation/{id} --> controller.GetSimulation()
// Next:
//     * On success --> service.GetSimulation
//     * On fail --> return error
func (ctrl *controller) GetSimulation(user *users.User, tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	groupID, ok := mux.Vars(r)["groupid"]
	if !ok {
		return nil, ign.NewErrorMessage(ign.ErrorIDNotInRequest)
	}

	var simulation Simulation
	if err := tx.Where("group_id=?", groupID).First(&simulation).Error; err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorIDNotFound, err)
	}

  if simulation.Owner != *user.Username {
    return nil, ign.NewErrorMessage(ign.ErrorUnauthorized)
  }

	// Send response to the user
	return GetSimulationResponse{
		Name:    simulation.Name,
		GroupID: simulation.GroupID,
		Status:  simulation.Status,
		Image:   simulation.Image,
		Args:    simulation.Args,
		URI:     simulation.URI,
		IP:      simulation.IP,
	}, nil
}

// Healthz returns a string to confirm that cloudsim is running.
func Healthz(tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	return "Cloudsim is up", nil
}

////////////////////////////////////////////////
// All of the following is to handle user access in routes
// HTTPHandler is used to invoke inner logic based on incoming Http requests.
type HTTPHandler struct {
	UserAccessor useracc.Service
}

// HTTPHandlerInstance is the default HTTPHandler instance. It is used by routes.go.
var HTTPHandlerInstance *HTTPHandler

// NewHTTPHandler creates a new HTTPHandler.
func NewHTTPHandler(ctx context.Context, ua useracc.Service) (*HTTPHandler, error) {
	return &HTTPHandler{
		UserAccessor: ua,
	}, nil
}

type handlerWithUser func(user *users.User, tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg)

// WithUser is a middleware that checks for a valid user from the JWT and passes
// the user to the handlerWithUser.
func WithUser(handler handlerWithUser) ign.HandlerWithResult {
	return func(tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
		// Get JWT user. Fail if invalid or missing
		user, ok, em := HTTPHandlerInstance.UserAccessor.UserFromJWT(r)
		if !ok {
			return nil, em
		}
		return handler(user, tx, w, r)
	}
}
////////////////////////////////////////////////
