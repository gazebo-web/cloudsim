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
	Update(track UpdateTrackInput) (*Track, error)
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

func (s service) Create(input CreateTrackInput) (*Track, error) {
	s.logger.Debug(fmt.Sprintf(" [Track.Service] Creating track. Input: %+v", input))
	if err := s.validator.Struct(&input); err != nil {
		s.logger.Debug(fmt.Sprintf(" [Track.Service] Validation failed. Error: %+v", err))
		return nil, err
	}
	track := CreateTrackFromInput(input)
	output, err := s.repository.Create(track)
	if err != nil {
		s.logger.Debug(fmt.Sprintf(" [Track.Service] Repository.Create() failed. Error: %+v", err))
		return nil, err
	}
	s.logger.Debug(fmt.Sprintf(" [Track.Service] Track created. Output: %+v", output))
	return output, nil
}

func (s service) Get(name string) (*Track, error) {
	panic("implement me")
}

func (s service) GetAll() ([]Track, error) {
	panic("implement me")
}

func (s service) Update(track UpdateTrackInput) (*Track, error) {
	panic("implement me")
}

func (s service) Delete(name string) (*Track, error) {
	panic("implement me")
}

func NewService(r Repository, v *validator.Validate, logger ign.Logger) Service {
	return &service{
		repository: r,
		validator:  v,
		logger:     logger,
	}
}
