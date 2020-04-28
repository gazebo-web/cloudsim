package simulations

import (
	"github.com/go-playground/form"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/users"
	"gitlab.com/ignitionrobotics/web/fuelserver/permissions"
	"gopkg.in/go-playground/validator.v9"
)

// IController represents a group of methods to expose in the API Rest.
type IController interface {
	simulations.IController
}

// Controller is an IController implementation.
type Controller struct {
	simulations.IController
	services services
}

type services struct {
	Simulation IService
	User users.IService
}

type NewControllerInput struct {
	Service IService
	Decoder *form.Decoder
	Validator *validator.Validate
	Permissions *permissions.Permissions
	UserService users.IService
}

func NewController(input NewControllerInput) IController {
	var c IController
	userService, err := users.NewService(input.Permissions, "")
	if err != nil {
		panic("Couldn't create new SubT Simulations Controller")
	}
	c = &Controller{
		IController: simulations.NewController(simulations.NewControllerInput{
			SimulationService: input.Service,
			UserService:       userService,
			FormDecoder:       input.Decoder,
			Validator:         input.Validator,
		}),
		services:    services{
			Simulation: input.Service,
			User: input.UserService,
		},
	}
	return c
}