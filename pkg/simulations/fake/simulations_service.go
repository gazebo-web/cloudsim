package fake

import (
	"github.com/stretchr/testify/mock"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

type Service struct {
	*mock.Mock
}

func (s *Service) Get(groupID simulations.GroupID) (simulations.Simulation, error) {
	args := s.Called(groupID)
	sim := args.Get(0).(simulations.Simulation)
	return sim, args.Error(1)
}

func NewService() *Service {
	return &Service{
		Mock: new(mock.Mock),
	}
}
