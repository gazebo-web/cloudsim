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

// Controller represents a group of methods to expose in the API Rest.
type Controller interface {
	Start(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg)
	LaunchHeld(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg)
	Restart(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg)
	Shutdown(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg)
	GetAll(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg)
	Get(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg)
	GetDownloadableLogs(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg)
	GetLiveLogs(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg)
}

// controller is an Controller implementation. This implementation serves as an example for the different applications.
type controller struct {
	services    services
	formDecoder *form.Decoder
	validator   *validator.Validate
}

// NewControllerInput is the input needed to create a new Controller implementation.
type NewControllerInput struct {
	Services    services
	FormDecoder *form.Decoder
	Validator   *validator.Validate
}

// NewController receives a NewControllerInput to initialize a new Controller implementation.
func NewController(input NewControllerInput) Controller {
	var c Controller

	if input.Services.Simulation == nil {
		panic("Simulation Service should not be nil")
	}
	if input.Services.User == nil {
		panic("UserService should not be nil")
	}
	if input.FormDecoder == nil {
		panic("formDecoder should not be nil")
	}
	if input.Validator == nil {
		panic("validator should not be nil")
	}

	c = &controller{
		services: services{
			Simulation: input.Services.Simulation,
			User:       input.Services.User,
		},
		formDecoder: input.FormDecoder,
		validator:   input.Validator,
	}
	return c
}

// services represents a set of services used by the controller.
type services struct {
	Simulation Service
	User       users.Service
}

// Start is the handler to create a new simulation.
func (c *controller) Start(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
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

// LaunchHeld is the handler to launch a given held simulation.
func (c *controller) LaunchHeld(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	groupID, ok := mux.Vars(r)["group"]
	if !ok {
		return nil, ign.NewErrorMessage(ign.ErrorIDNotInRequest)
	}
	return c.services.Simulation.Launch(r.Context(), groupID, user)
}

// Restart is the handler to restart a given simulation.
func (c *controller) Restart(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	groupID, ok := mux.Vars(r)["group"]
	if !ok {
		return nil, ign.NewErrorMessage(ign.ErrorIDNotInRequest)
	}
	return c.services.Simulation.Restart(r.Context(), groupID, user)
}

// Shutdown is the handler to stop a given simulation.
func (c *controller) Shutdown(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	groupID, ok := mux.Vars(r)["group"]
	if !ok {
		return nil, ign.NewErrorMessage(ign.ErrorIDNotInRequest)
	}
	return c.services.Simulation.Shutdown(r.Context(), groupID, user)
}

func (c *controller) GetAll(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
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
	// var circuit *string
	// if len(params["circuit"]) > 0 && len(params["circuit"][0]) > 0 {
	// 	circuit = &params["circuit"][0]
	// }

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

func (c *controller) Get(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	groupID, ok := mux.Vars(r)["group"]
	if !ok {
		return nil, ign.NewErrorMessage(ign.ErrorIDNotInRequest)
	}

	return c.services.Simulation.Get(groupID, user)
}

func (c *controller) GetDownloadableLogs(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	panic("Not implemented")

}

func (c *controller) GetLiveLogs(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	panic("Not implemented")
}
