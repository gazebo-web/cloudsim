package tracks

import (
	"errors"
	"fmt"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/pkg/domain"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/pkg/repositories"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"gopkg.in/go-playground/validator.v9"
)

// Service groups a set of methods that have the business logic to perform CRUD operations for a Track.
type Service interface {
	Create(input CreateTrackInput) (*Track, error)
	Get(name string) (*Track, error)
	GetAll(page, pageSize *int) ([]Track, error)
	Update(name string, input UpdateTrackInput) (interface{}, error)
	Delete(name string) (*Track, error)
}

type Config struct {

}

type service struct {
	repository repositories.Repository
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
	output, err := s.repository.Create([]domain.Entity{track})
	if err != nil {
		s.logger.Debug(fmt.Sprintf(" [Track.Service] Creation failed failed. Error: %+v", err))
		return nil, err
	}
	s.logger.Debug(fmt.Sprintf(" [Track.Service] Track created. Output: %+v", output))
	return track, nil
}

// Get gets a Track with the given name.
func (s service) Get(name string) (*Track, error) {
	s.logger.Debug(" [Track.Service] Getting Track with name: ", name)
	var track Track
	err := s.repository.FindOne(&track, repositories.NewGormFilter("name = ?", name))
	if err != nil {
		s.logger.Debug(fmt.Sprintf(" [Track.Service] Getting track with name %s failed. Error: %+v", name, err))
		return nil, err
	}
	s.logger.Debug(fmt.Sprintf(" [Track.Service] Track found. Output: %+v", track))
	return &track, nil
}

// GetAll returns a slice with 10 tracks from the first track in the database.
// If `page` and `pageSize` are not nil, it will return `pageSize` tracks starting from the `page` track.
func (s service) GetAll(page, pageSize *int) ([]Track, error) {
	s.logger.Debug(" [Track.Service] Getting all tracks")
	var tracks []Track

	count, err := s.repository.Count()

	// TODO: Move limit constant to another place.
	limit := 10
	if pageSize == nil {
		pageSize = &limit
	}

	if page != nil && count < (*page+1) * *pageSize {
		// TODO: Extract error into variable.
		return nil, errors.New("invalid page")
	}

	err = s.repository.Find(&tracks, page, pageSize)
	if err != nil {
		s.logger.Debug(fmt.Sprintf(" [Track.Service] Getting tracks failed. Error: %+v", err))
		return nil, err
	}
	s.logger.Debug(fmt.Sprintf(" [Track.Service] Getting all tracks succeeded. Tracks: %+v", tracks))
	return tracks, nil
}

// Update updates a track with the given name and populates it with information provided by the given input.
func (s service) Update(name string, input UpdateTrackInput) (interface{}, error) {
	s.logger.Debug(fmt.Sprintf(" [Track.Service] Updating track with name: %s. Input: %+v", name, input))
	err := s.validator.Struct(&input)
	if err != nil {
		s.logger.Debug(fmt.Sprintf(" [Track.Service] Validating input failed. Error: %+v", err))
		return nil, err
	}
	updatedTrack, err := input.ToMap()
	if err != nil {
		return nil, err
	}
	err = s.repository.Update(updatedTrack, repositories.NewGormFilter("name = ?", name))
	if err != nil {
		s.logger.Debug(fmt.Sprintf(" [Track.Service] Updating track with name: %s failed. Error: %+v", name, err))
		return nil, err
	}
	s.logger.Debug(fmt.Sprintf(" [Track.Service] Updating track with name: %s succeeded. Output: %+v", name, updatedTrack))
	return input, nil
}

// Delete removes the track with the given name.
func (s service) Delete(name string) (*Track, error) {
	s.logger.Debug(fmt.Sprintf(" [Track.Service] Removing track with name: %s", name))
	entity, err := s.Get(name)
	if err != nil {
		return nil, err
	}
	err = s.repository.Delete(repositories.NewGormFilter("name = ?", name))
	if err != nil {
		s.logger.Debug(fmt.Sprintf(" [Track.Service] Deleting the track with name: %s failed. Error: %+v", name, err))
		return nil, err
	}
	s.logger.Debug(fmt.Sprintf(" [Track.Service] Track deleted with name: %s. Track: %+v", name, *entity))
	return entity, nil
}

// NewService initializes a new Service implementation
func NewService(r repositories.Repository, v *validator.Validate, logger ign.Logger) Service {
	return &service{
		repository: r,
		validator:  v,
		logger:     logger,
	}
}
