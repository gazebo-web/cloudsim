package simulations

import (
	"github.com/go-playground/form"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/users"
	"gitlab.com/ignitionrobotics/web/cloudsim/tools"
	fuel "gitlab.com/ignitionrobotics/web/fuelserver/bundles/users"
	"gitlab.com/ignitionrobotics/web/fuelserver/permissions"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
)

// IController represents a group of methods to expose in the API Rest.
type IController interface {
	simulations.IController
}

// Controller is an IController implementation.
type Controller struct {
	services    services
	formDecoder *form.Decoder
	validator   *validator.Validate
}

type services struct {
	Simulation IService
	User       users.IService
}

type NewControllerInput struct {
	Service     IService
	Decoder     *form.Decoder
	Validator   *validator.Validate
	Permissions *permissions.Permissions
	UserService users.IService
}

func NewController(input NewControllerInput) IController {
	var c IController
	c = &Controller{
		formDecoder: input.Decoder,
		validator:   input.Validator,
		services: services{
			Simulation: input.Service,
			User:       input.UserService,
		},
	}
	return c
}

func (c *Controller) Start(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
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

func (c *Controller) LaunchHeld(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	panic("implement me")
}

func (c *Controller) Restart(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	panic("implement me")
}

func (c *Controller) Shutdown(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	panic("implement me")
}

func (c *Controller) GetAll(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	panic("implement me")
}

func (c *Controller) Get(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	panic("implement me")
}

func (c *Controller) GetDownloadableLogs(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	panic("implement me")
}

func (c *Controller) GetLiveLogs(user *fuel.User, w http.ResponseWriter, r *http.Request) (interface{}, *ign.ErrMsg) {
	panic("implement me")
}
