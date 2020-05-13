package simulations

import (
	"github.com/go-playground/form"
	"github.com/gorilla/mux"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/users"
	"gitlab.com/ignitionrobotics/web/cloudsim/tools"
	fuel "gitlab.com/ignitionrobotics/web/fuelserver/bundles/users"
	"gitlab.com/ignitionrobotics/web/fuelserver/permissions"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
)

// controller represents a group of methods to expose in the API Rest.
type Controller interface {
	simulations.Controller
}

// controller is an controller implementation.
type controller struct {
	services    services
	formDecoder *form.Decoder
	validator   *validator.Validate
}

type services struct {
	Simulation IService
	User       users.Service
}

type NewControllerInput struct {
	Service     IService
	Decoder     *form.Decoder
	Validator   *validator.Validate
	Permissions *permissions.Permissions
	UserService users.Service
}

func NewController(input NewControllerInput) Controller {
	var c Controller
	c = &controller{
		formDecoder: input.Decoder,
		validator:   input.Validator,
		services: services{
			Simulation: input.Service,
			User:       input.UserService,
		},
	}
	return c
}

func (c *controller) Start(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	if err := r.ParseMultipartForm(0); err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorForm, err)
	}
	defer r.MultipartForm.RemoveAll()

	var createSim SimulationCreate
	// TODO: Create function for Parse & Validate.
	if em := tools.ParseFormStruct(&createSim.SimulationCreate, r, c.formDecoder); em != nil {
		return nil, em
	}

	if em := tools.ValidateStruct(&createSim, c.validator); em != nil {
		return nil, em
	}

	return c.services.Simulation.Create(r.Context(), &createSim, user)
}

func (c *controller) LaunchHeld(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	panic("implement me")
}

func (c *controller) Restart(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	panic("implement me")
}

func (c *controller) Shutdown(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	panic("implement me")
}

func (c *controller) GetAll(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	panic("implement me")
}

func (c *controller) Get(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	groupID, ok := mux.Vars(r)["group"]
	if !ok {
		return nil, ign.NewErrorMessage(ign.ErrorIDNotInRequest)
	}

	return c.services.Simulation.Get(groupID, user)
}

func (c *controller) GetDownloadableLogs(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	panic("implement me")
}

func (c *controller) GetLiveLogs(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	panic("implement me")
}
