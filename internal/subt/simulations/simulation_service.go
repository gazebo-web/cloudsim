package simulations

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

type IService interface {
	simulations.IService
}

type Service struct {
	*simulations.Service
}