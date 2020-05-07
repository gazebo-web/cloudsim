package circuits

import "gopkg.in/go-playground/validator.v9"

type IService interface {
	GetByName(name string) (*Circuit, error)
	IsValidCircuit(fl validator.FieldLevel) bool
}

type Service struct {
	repository IRepository
}

func NewService(repository IRepository) IService {
	var s IService
	s = &Service{
		repository: repository,
	}
	return s
}

func (s *Service) GetByName(name string) (*Circuit, error) {
	return s.repository.GetByName(name)
}

func (s *Service) IsValidCircuit(fl validator.FieldLevel) bool {
	circuit, err := s.GetByName(fl.Field().String())
	if err != nil {
		return false
	}
	return circuit.Enabled
}
