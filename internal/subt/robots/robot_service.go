package robots

import (
	"gitlab.com/ignitionrobotics/web/ign-go"
	"gopkg.in/go-playground/validator.v9"
)

type IService interface {
	GetAllTypes() (*map[string]Robot, *ign.ErrMsg)
	GetByType(robotType string) (*Robot, error)
	IsValidRobotType(fl validator.FieldLevel) bool
}

type Service struct {
	repository IRepository
}

func NewService() IService {
	var s IService
	s = &Service{}
	return s
}

func (s *Service) GetAllTypes() (*map[string]Robot, *ign.ErrMsg) {
	mappedRobots, err := s.repository.GetAll()
	if err != nil {
		// TODO: Change error type
		return nil, ign.NewErrorMessageWithBase(ign.ErrorUnexpected, err)
	}
	return &mappedRobots, nil
}

func (s *Service) GetByType(robotType string) (*Robot, error) {
	return s.repository.GetByType(robotType)
}

func (s *Service) IsValidRobotType(fl validator.FieldLevel) bool {
	_, err := s.GetByType(fl.Field().String())
	if err != nil {
		return false
	}
	return true
}
