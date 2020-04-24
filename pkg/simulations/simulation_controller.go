package simulations

import (
	"github.com/go-playground/form"
	"github.com/gorilla/mux"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/users"
	"gitlab.com/ignitionrobotics/web/cloudsim/tools"
	fuel "gitlab.com/ignitionrobotics/web/fuelserver/bundles/users"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strconv"
)

// IController represents a group of methods to expose in the API Rest.
type IController interface {
	Start(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg)
	LunchHeld(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg)
	Restart(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg)
	Shutdown(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg)
	GetAll(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg)
	Get(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg)
	GetDownloadableLogs(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg)
	GetLiveLogs(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg)
}

// Controller is an IController implementation.
type Controller struct {
	services services
	formDecoder *form.Decoder
	validator *validator.Validate
}

// NewControllerInput is the input needed to create a new IController implementation.
type NewControllerInput struct {
	SimulationService IService
	UserService users.IService
	FormDecoder *form.Decoder
	Validator *validator.Validate
}

// NewController receives a NewControllerInput to initialize a new IController implementation.
func NewController(input NewControllerInput) IController {
	var c IController

	if input.SimulationService == nil {
		panic("SimulationService should not be nil")
	}
	if input.UserService == nil {
		panic("UserService should not be nil")
	}
	if input.FormDecoder == nil {
		panic("FormDecoder should not be nil")
	}
	if input.Validator == nil {
		panic("Validator should not be nil")
	}

	c = &Controller{
		services:    services{
			Simulation: input.SimulationService,
			User: input.UserService,
		},
		formDecoder: input.FormDecoder,
		validator:   input.Validator,
	}
	return c
}

// services represents a set of services used by the Controller.
type services struct {
	Simulation IService
	User users.IService
}

// Start is the handler to create a new simulation.
func (c *Controller) Start(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	if err := r.ParseMultipartForm(0); err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorForm, err)
	}
	defer r.MultipartForm.RemoveAll()

	var createSim SimulationCreate
	// TODO: Create function for Parse & Validate.
	if em := tools.ParseFormStruct(&createSim, r, c.formDecoder); em != nil {
		return nil, em
	}

	if em := tools.ValidateStruct(&createSim, c.validator); em != nil {
		return nil, em
	}

	return c.services.Simulation.Create(r.Context(), &createSim, user)
}

// LunchHeld is the handler to launch a given held simulation.
func (c *Controller) LunchHeld(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	groupID, ok := mux.Vars(r)["group"]
	if !ok {
		return nil, ign.NewErrorMessage(ign.ErrorIDNotInRequest)
	}
	return c.services.Simulation.Launch(r.Context(), groupID, user)
}

// Restart is the handler to restart a given simulation.
func (c *Controller) Restart(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg)  {
	groupID, ok := mux.Vars(r)["group"]
	if !ok {
		return nil, ign.NewErrorMessage(ign.ErrorIDNotInRequest)
	}
	return c.services.Simulation.Restart(r.Context(), groupID, user)
}

// Shutdown is the handler to stop a given simulation.
func (c *Controller) Shutdown(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	groupID, ok := mux.Vars(r)["group"]
	if !ok {
		return nil, ign.NewErrorMessage(ign.ErrorIDNotInRequest)
	}
	return c.services.Simulation.Shutdown(r.Context(), groupID, user)
}

func (c *Controller) GetAll(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	// Prepare pagination
	pr, em := ign.NewPaginationRequest(r)
	if em != nil {
		return nil, em
	}

	// Get the parameters
	params := r.URL.Query()
	var status *Status
	invertStatus := false
	invertErrStatus := false
	if len(params["status"]) > 0 && len(params["status"][0]) > 0 {
		statusStr := params["status"][0]
		invertStatus = statusStr[0] == '!'
		sliceIndex := 0
		if invertStatus {
			sliceIndex = 1
		}
		status = NewStatus(statusStr[sliceIndex:])
	}
	var errStatus *ErrorStatus
	if len(params["errorStatus"]) > 0 && len(params["errorStatus"][0]) > 0 {
		statusStr := params["errorStatus"][0]
		invertErrStatus = statusStr[0] == '!'
		sliceIndex := 0
		if invertErrStatus {
			sliceIndex = 1
		}
		err := ErrorStatus(statusStr[sliceIndex:])
		errStatus = &err
	}

	includeChildren := false
	if len(params["children"]) > 0 && len(params["children"][0]) > 0 {
		if flag, err := strconv.ParseBool(params["children"][0]); err == nil {
			includeChildren = flag
		}
	}

	// TODO: This is SubT specific and should be moved
	var circuit *string
	if len(params["circuit"]) > 0 && len(params["circuit"][0]) > 0 {
		circuit = &params["circuit"][0]
	}

	sims, pagination, em := c.services.Simulation.GetAll(r.Context(), GetAllInput{
		p:               pr,
		byStatus:        status,
		invertStatus:    invertStatus,
		byErrStatus:     errStatus,
		invertErrStatus: invertErrStatus,
		user:            user,
		includeChildren: includeChildren,
	})

	if em != nil {
		return nil, em
	}

	ign.WritePaginationHeaders(*pagination, w, r)
	return sims, nil
}

func (c *Controller) Get(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {

}

func (c *Controller) GetDownloadableLogs(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {

}

func (c *Controller) GetLiveLogs(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {

}