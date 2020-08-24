package robots

import (
	"gitlab.com/ignitionrobotics/web/ign-go"
	"gopkg.in/go-playground/validator.v9"
)

type IService interface {
	GetAllConfigs() (*[]RobotConfig, *ign.ErrMsg)
	GetConfigByType(robotType string) (*RobotConfig, error)
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

func (s *Service) GetAllConfigs() (*[]RobotConfig, *ign.ErrMsg) {
	robotCfgs, err := s.repository.GetAllConfigs()
	if err != nil {
		// TODO: Change error type
		return nil, ign.NewErrorMessageWithBase(ign.ErrorUnexpected, err)
	}
	return &robotCfgs, nil
}

func (s *Service) GetConfigByType(robotType string) (*RobotConfig, error) {
	return s.repository.GetConfigByType(robotType)
}

func (s *Service) IsValidRobotType(fl validator.FieldLevel) bool {
	_, err := s.GetConfigByType(fl.Field().String())
	if err != nil {
		return false
	}
	return true
}
