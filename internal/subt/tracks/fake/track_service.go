package fake

import (
	"github.com/stretchr/testify/mock"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/tracks"
)

// Service is a fake tracks.Service implementation.
type Service struct {
	*mock.Mock
}

// Create mocks the Create method.
func (s *Service) Create(input tracks.CreateTrackInput) (*tracks.Track, error) {
	args := s.Called(input)
	return args.Get(0).(*tracks.Track), args.Error(1)
}

// Get mocks the Get method.
func (s *Service) Get(name string, id int) (*tracks.Track, error) {
	args := s.Called(name, id)
	return args.Get(0).(*tracks.Track), args.Error(1)
}

// GetAll mocks the GetAll method.
func (s *Service) GetAll() ([]tracks.Track, error) {
	args := s.Called()
	return args.Get(0).([]tracks.Track), args.Error(1)
}

// Update mocks the Update method.
func (s *Service) Update(name string, input tracks.UpdateTrackInput) (*tracks.Track, error) {
	args := s.Called(name, input)
	return args.Get(0).(*tracks.Track), args.Error(1)
}

// Delete mocks the Delete method.
func (s *Service) Delete(name string) (*tracks.Track, error) {
	args := s.Called(name)
	return args.Get(0).(*tracks.Track), args.Error(1)
}

// NewService initializes a fake tracks.Service implementation.
func NewService() *Service {
	return &Service{
		Mock: new(mock.Mock),
	}
}
