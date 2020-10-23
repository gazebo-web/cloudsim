package fake

import (
	"github.com/stretchr/testify/mock"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/tracks"
)

// fakeService is a fake tracks.Service implementation.
type fakeService struct {
	*mock.Mock
}

// Create mocks the Create method.
func (s *fakeService) Create(input tracks.CreateTrackInput) (*tracks.Track, error) {
	args := s.Called(input)
	return args.Get(0).(*tracks.Track), args.Error(1)
}

// Get mocks the Get method.
func (s *fakeService) Get(name string) (*tracks.Track, error) {
	args := s.Called(name)
	return args.Get(0).(*tracks.Track), args.Error(1)
}

// GetAll mocks the GetAll method.
func (s *fakeService) GetAll() ([]tracks.Track, error) {
	args := s.Called()
	return args.Get(0).([]tracks.Track), args.Error(1)
}

// Update mocks the Update method.
func (s *fakeService) Update(name string, input tracks.UpdateTrackInput) (*tracks.Track, error) {
	args := s.Called(name, input)
	return args.Get(0).(*tracks.Track), args.Error(1)
}

// Delete mocks the Delete method.
func (s *fakeService) Delete(name string) (*tracks.Track, error) {
	args := s.Called(name)
	return args.Get(0).(*tracks.Track), args.Error(1)
}

// NewService initializes a fake tracks.Service implementation.
func NewService() *fakeService {
	return &fakeService{
		Mock: new(mock.Mock),
	}
}
