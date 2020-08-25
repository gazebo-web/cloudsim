package fake

import (
	"github.com/stretchr/testify/mock"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// Service is a fake simulations.Service implementation.
type Service struct {
	*mock.Mock
}

// Reject is a mock for the Reject method.
func (s *Service) Reject(groupID simulations.GroupID) (simulations.Simulation, error) {
	args := s.Called(groupID)
	sim := args.Get(0).(simulations.Simulation)
	return sim, args.Error(1)
}

// GetParent is a mock for the GetParent method.
func (s *Service) GetParent(groupID simulations.GroupID) (simulations.Simulation, error) {
	args := s.Called(groupID)
	sim := args.Get(0).(simulations.Simulation)
	return sim, args.Error(1)
}

// UpdateStatus is a mock for the UpdateStatus method.
func (s *Service) UpdateStatus(groupID simulations.GroupID, status simulations.Status) error {
	args := s.Called(groupID)
	return args.Error(0)
}

// Update is a mock for the Update method.
func (s *Service) Update(groupID simulations.GroupID, simulation simulations.Simulation) error {
	args := s.Called(groupID)
	return args.Error(0)
}

// Get is a mock for the Get method.
func (s *Service) Get(groupID simulations.GroupID) (simulations.Simulation, error) {
	args := s.Called(groupID)
	sim := args.Get(0).(simulations.Simulation)
	return sim, args.Error(1)
}

// NewService initializes a new fake service implementation.
func NewService() *Service {
	return &Service{
		Mock: new(mock.Mock),
	}
}
