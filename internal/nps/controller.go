package nps

// This file defines the controller, which handles route requests. An
// application creates an instance of a controller by calling `NewController`.

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/form"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	useracc "gitlab.com/ignitionrobotics/web/cloudsim/pkg/users"
	"gitlab.com/ignitionrobotics/web/fuelserver/bundles/users"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"net/http"
	"strings"
	"time"
)

// Controller is an interface designed to handle route requests.
type Controller interface {
	Start(user *users.User, tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg)
	Stop(user *users.User, tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg)
	ListSimulations(user *users.User, tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg)
	GetSimulation(user *users.User, tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg)
	GetUser(user *users.User, tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg)
	AddUser(user *users.User, tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg)
	DeleteUser(user *users.User, tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg)
	ModifyUser(user *users.User, tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg)
	ListUsers(user *users.User, tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg)
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

	// If not a system admin, then check that the user is registered.
	if ok := HTTPHandlerInstance.UserAccessor.IsSystemAdmin(*user.Username); !ok {
		// Check that the user is registered.
		var registeredUser RegisteredUser
		if err := tx.Where("username = ?", *user.Username).Find(&registeredUser).Error; err != nil {
			return nil, ign.NewErrorMessage(ign.ErrorUnauthorized)
		}

		// Check simulation limits.
		if registeredUser.SimulationLimit >= 0 {
			var simulations Simulations
			tx.Where("owner = ? and status != ?", registeredUser.Username, "stopped").Find(&simulations)
			if len(simulations) >= registeredUser.SimulationLimit {
				return nil, ign.NewErrorMessageWithArgs(ign.ErrorUnauthorized, errors.New("Simulation limit reached"), []string{"Simulation limit reached"})
			}
		}
	}

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

	// The name must be lowercase
	req.Name = strings.ToLower(req.Name)
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

	if groupID == "all" {
		var simulations Simulations

		// Get all the simulation if the user is a system admin
		if ok := HTTPHandlerInstance.UserAccessor.IsSystemAdmin(*user.Username); ok {
			tx.Find(&simulations)
		} else {

			// Get all the simulations owned by the user
			if err := tx.Where("owner=?", user.Username).Find(&simulations).Error; err != nil {
				return nil, ign.NewErrorMessageWithBase(ign.ErrorIDNotFound, err)
			}
		}

		for _, sim := range simulations {
			// Construct the stop request to send to the service
			req := StopRequest{
				GroupID: sim.GroupID,
			}

			if !strings.Contains(sim.Status, "stop") && !strings.Contains(sim.Status, "remov") {
				ctrl.service.Stop(user, tx, r.Context(), req)
			}
		}

		return "All instances stopped", nil

	} else {
		// Get the matching simulation
		var simulation Simulation
		if err := tx.Where("group_id=?", groupID).First(&simulation).Error; err != nil {
			return nil, ign.NewErrorMessageWithBase(ign.ErrorIDNotFound, err)
		}

		if !strings.Contains(simulation.Status, "stop") && !strings.Contains(simulation.Status, "remov") {
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
	}

	return nil, nil
}

// ListSimulations handles the `/simulations` route.
//
// Origin: user --> GET /simulations --> controller.ListSimulations()
// Next:
//     * On success --> return ListResponse
//     * On fail --> return error
func (ctrl *controller) ListSimulations(user *users.User, tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {

	var simulations Simulations

	var simQuery Simulation

	// Get and parse the filter parameters
	if query, ok := r.URL.Query()["q"]; ok {
		for _, q := range query {
			parts := strings.Split(q, ":")
			if len(parts) == 2 {
				if strings.ToLower(parts[0]) == "status" {
					simQuery.Status = strings.ToLower(parts[1])
				} else if strings.ToLower(parts[0]) == "name" {
					simQuery.Name = strings.ToLower(parts[1])
				} else if strings.ToLower(parts[0]) == "groupid" {
					simQuery.GroupID = strings.ToLower(parts[1])
				}
			}
		}
	}

	// Return only owner simulations if the user is not a system admin
	if ok := HTTPHandlerInstance.UserAccessor.IsSystemAdmin(*user.Username); !ok {
		simQuery.Owner = *user.Username
	}
	tx.Where(&simQuery).Find(&simulations)

	var response ListResponse
	for _, sim := range simulations {
		response.Simulations = append(response.Simulations, GetSimulationResponse{
			Name:    sim.Name,
			GroupID: sim.GroupID,
			Status:  sim.Status,
			Image:   sim.Image,
			Args:    sim.Args,
			URI:     sim.URI,
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
// ModifyUser handles the PATCH `/user/{USERNAME}` route.
//
// Origin: user --> PATCH /user/{USERNAME} --> controller.ModifyUser()
// Next:
//     * On success --> return User
//     * On fail --> return error
func (ctrl *controller) ModifyUser(user *users.User, tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {

	// Fail if not system admin
	if ok := HTTPHandlerInstance.UserAccessor.IsSystemAdmin(*user.Username); !ok {
		return nil, ign.NewErrorMessage(ign.ErrorUnauthorized)
	}

	// Get the username from the route
	userName, ok := mux.Vars(r)["username"]
	if !ok {
		return nil, ign.NewErrorMessage(ign.ErrorUserNotInRequest)
	}

	// Check that the user already exists
	var existing RegisteredUser
	if err := tx.Where("username = ?", userName).Find(&existing).Error; err != nil {
		return nil, ign.NewErrorMessage(ign.ErrorNonExistentResource)
	}

	// Parse form's values and files.
	if err := r.ParseMultipartForm(0); err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorForm, err)
	}
	defer r.MultipartForm.RemoveAll()

	var req AddModifyUserRequest
	if errs := ctrl.formDecoder.Decode(&req, r.Form); errs != nil {
		return nil, ign.NewErrorMessageWithArgs(ign.ErrorFormInvalidValue, errs,
			getDecodeErrorsExtraInfo(errs))
	}

	// If modifying the username, first check that the new username doesn't exists
	// and is not empty.
	if req.Username != "" && *existing.Username != req.Username {
		var userCheck RegisteredUser
		if err := tx.Where("username = ?", req.Username).Find(&userCheck).Error; err != nil {
			existing.Username = &req.Username
		}
	}

	// Update the simulation limit.
	existing.SimulationLimit = req.SimulationLimit

	if err := tx.Save(&existing).Error; err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorDbSave, err)
	}

	return AddModifyUserResponse{
		Username:        *existing.Username,
		SimulationLimit: existing.SimulationLimit,
	}, nil
}

////////////////////////////////////////////////
// AddUser handles the POST `/users` route.
//
// Origin: user --> POST /users --> controller.AddUser()
// Next:
//     * On success --> return User
//     * On fail --> return error
func (ctrl *controller) AddUser(user *users.User, tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {

	// Fail if not system admin
	if ok := HTTPHandlerInstance.UserAccessor.IsSystemAdmin(*user.Username); !ok {
		return nil, ign.NewErrorMessage(ign.ErrorUnauthorized)
	}

	// Parse form's values and files.
	if err := r.ParseMultipartForm(0); err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorForm, err)
	}
	defer r.MultipartForm.RemoveAll()

	var req AddModifyUserRequest

	if errs := ctrl.formDecoder.Decode(&req, r.Form); errs != nil {
		return nil, ign.NewErrorMessageWithArgs(ign.ErrorFormInvalidValue, errs,
			getDecodeErrorsExtraInfo(errs))
	}

	// Check that the user doesn't already exist
	var existing RegisteredUser
	if err := tx.Where("username = ?", req.Username).Find(&existing).Error; err == nil {
		return nil, ign.NewErrorMessage(ign.ErrorResourceExists)
	}

	// If the simulation_limit form value is not present, then assume unlimited
	// simulations
	if _, ok := r.Form["simulation_limit"]; !ok {
		req.SimulationLimit = -1
	}

	newUser := RegisteredUser{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),

		Username:        &req.Username,
		SimulationLimit: req.SimulationLimit,
	}

	if err := tx.Create(&newUser).Error; err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorDbSave, err)
	}

	return AddModifyUserResponse{
		Username:        req.Username,
		SimulationLimit: req.SimulationLimit,
	}, nil
}

////////////////////////////////////////////////
// GetUser handles the GET `/user/{username}` route.
//
// Origin: user --> GET /user/{username} --> controller.GetUser()
// Next:
//     * On success --> return UserResponse
//     * On fail --> return error
func (ctrl *controller) GetUser(user *users.User, tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {

	// Get the username from the route
	userName, ok := mux.Vars(r)["username"]
	if !ok {
		return nil, ign.NewErrorMessage(ign.ErrorUserNotInRequest)
	}

	// Return if not an admin and not the user.
	if ok := HTTPHandlerInstance.UserAccessor.IsSystemAdmin(*user.Username); !ok {
		if *user.Username != userName {
			return nil, ign.NewErrorMessage(ign.ErrorUnauthorized)
		}
	}

	var userEntry RegisteredUser
	if err := tx.Where("username=?", userName).Find(&userEntry).Error; err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorIDNotFound, err)
	}

	return AddModifyUserResponse{
		Username:        *userEntry.Username,
		SimulationLimit: userEntry.SimulationLimit,
	}, nil
}

////////////////////////////////////////////////
// DeleteUser handles the DELETE `/user/{username}` route.
//
// Origin: user --> DELETE /user/{username} --> controller.DeleteUser()
// Next:
//     * On success --> return UserResponse
//     * On fail --> return error
func (ctrl *controller) DeleteUser(user *users.User, tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {

	// Get the username from the route
	userName, ok := mux.Vars(r)["username"]
	if !ok {
		return nil, ign.NewErrorMessage(ign.ErrorUserNotInRequest)
	}

	// Return if not an admin.
	if ok := HTTPHandlerInstance.UserAccessor.IsSystemAdmin(*user.Username); !ok {
		return nil, ign.NewErrorMessage(ign.ErrorUnauthorized)
	}

	var userEntry RegisteredUser
	if err := tx.Where("username=?", userName).Find(&userEntry).Error; err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorIDNotFound, err)
	}

	result := AddModifyUserResponse{
		Username:        *userEntry.Username,
		SimulationLimit: userEntry.SimulationLimit,
	}

	if err := tx.Delete(&userEntry).Error; err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorDbDelete, err)
	}

	return result, nil
}

////////////////////////////////////////////////
// ListUsers handles the GET `/users` route.
//
// Origin: user --> GET /users --> controller.ListUsers()
// Next:
//     * On success --> return ListResponse
//     * On fail --> return error
func (ctrl *controller) ListUsers(user *users.User, tx *gorm.DB, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {

	// Return only owner simulations if the user is not a system admin
	if ok := HTTPHandlerInstance.UserAccessor.IsSystemAdmin(*user.Username); !ok {
		fmt.Printf("Not admin[%s]\n", *user.Username)
		return nil, ign.NewErrorMessage(ign.ErrorUnauthorized)
	}

	var users RegisteredUsers
	tx.Find(&users)

	return users, nil
}
