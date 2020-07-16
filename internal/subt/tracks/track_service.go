package tracks

import (
	"fmt"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"gopkg.in/go-playground/validator.v9"
)

// Service groups a set of methods that have the business logic to perform CRUD operations for a Track.
type Service interface {
	serviceCreate
	serviceRead
	serviceUpdate
	serviceDelete
}

// serviceCreate has the business logic for creating a Track.
type serviceCreate interface {
	Create(input CreateTrackInput) (*Track, error)
}

// serviceRead has the business logic for reading one or multiple Tracks.
type serviceRead interface {
	Get(name string) (*Track, error)
	GetAll() ([]Track, error)
}

// serviceUpdate has the business logic for updating a Track.
type serviceUpdate interface {
	Update(name string, input UpdateTrackInput) (*Track, error)
}

// serviceDelete has the business logic for deleting a Track.
type serviceDelete interface {
	Delete(name string) (*Track, error)
}

type service struct {
	repository Repository
	validator  *validator.Validate
	logger     ign.Logger
}

// Create creates a new Track from the given Input.
func (s service) Create(input CreateTrackInput) (*Track, error) {
	s.logger.Debug(fmt.Sprintf(" [Track.Service] Creating track. Input: %+v", input))
	if err := s.validator.Struct(&input); err != nil {
		s.logger.Debug(fmt.Sprintf(" [Track.Service] Validation failed. Error: %+v", err))
		return nil, err
	}
	track := CreateTrackFromInput(input)
	output, err := s.repository.Create(track)
	if err != nil {
		s.logger.Debug(fmt.Sprintf(" [Track.Service] Creation failed failed. Error: %+v", err))
		return nil, err
	}
	s.logger.Debug(fmt.Sprintf(" [Track.Service] Track created. Output: %+v", output))
	return output, nil
}

// Get gets a Track with the given name.
func (s service) Get(name string) (*Track, error) {
	s.logger.Debug(" [Track.Service] Getting Track with name: ", name)
	track, err := s.repository.Get(name)
	if err != nil {
		s.logger.Debug(fmt.Sprintf(" [Track.Service] Getting track with name %s failed. Error: %+v", name, err))
		return nil, err
	}
	s.logger.Debug(fmt.Sprintf(" [Track.Service] Track found. Output: %+v", track))
	return track, nil
}

// GetAll returns a slice with all the tracks.
func (s service) GetAll() ([]Track, error) {
	s.logger.Debug(" [Track.Service] Getting all tracks")
	tracks, err := s.repository.GetAll()
	if err != nil {
		s.logger.Debug(fmt.Sprintf(" [Track.Service] Getting tracks failed. Error: %+v", err))
		return nil, err
	}
	s.logger.Debug(fmt.Sprintf(" [Track.Service] Getting all tracks succeeded. Tracks: %+v", tracks))
	return tracks, nil
}

// Update updates a track with the given name and populates it with information provided by the given input.
func (s service) Update(name string, input UpdateTrackInput) (*Track, error) {
	s.logger.Debug(fmt.Sprintf(" [Track.Service] Updating track with name: %s. Input: %+v", name, input))
	track, err := s.Get(name)
	if err != nil {
		return nil, err
	}
	updatedTrack := UpdateTrackFromInput(*track, input)
	track, err = s.repository.Update(name, updatedTrack)
	if err != nil {
		s.logger.Debug(fmt.Sprintf(" [Track.Service] Updating track with name: %s failed. Error: %+v", name, err))
		return nil, err
	}
	s.logger.Debug(fmt.Sprintf(" [Track.Service] Updating track with name: %s succeeded. Output: %+v", name, track))
	return track, nil
}

// Delete removes the track with the given name.
func (s service) Delete(name string) (*Track, error) {
	panic("implement me")
}

// NewService initializes a new Service implementation
func NewService(r Repository, v *validator.Validate, logger ign.Logger) Service {
	return &service{
		repository: r,
		validator:  v,
		logger:     logger,
	}
}
