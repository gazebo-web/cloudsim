package simulations

import (
	"github.com/go-playground/form"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/users"
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
	formDecoder *form.Decoder
	validator *validator.Validate
}

type services struct {
	Simulation IService
	User users.IService
}