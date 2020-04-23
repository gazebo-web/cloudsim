package simulations

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

type IService interface {
	simulations.IService
	CountByOwnerAndCircuit(owner, circuit string) (*int, error)
}

type Service struct {
	*simulations.Service
}

func NewService(repository IRepository) IService {
	var s IService
	parent := simulations.NewService(repository)
	service := parent.(*simulations.Service)
	s = &Service{
		Service: service,
	}
	return s
}

func (s *Service) CountByOwnerAndCircuit(owner, circuit string) (*int, error) {
	panic("Not implemented")
}