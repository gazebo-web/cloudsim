package circuits

import "gopkg.in/go-playground/validator.v9"

type Service interface {
	GetPending() ([]Circuit, error)
	GetByName(name string) (*Circuit, error)
	IsValidCircuit(fl validator.FieldLevel) bool
}

type service struct {
	repository Repository
}

func (s *service) GetPending() ([]Circuit, error) {
	return s.repository.GetPending()
}

func NewService(repository Repository) Service {
	var s Service
	s = &service{
		repository: repository,
	}
	return s
}

func (s *service) GetByName(name string) (*Circuit, error) {
	return s.repository.GetByName(name)
}

func (s *service) IsValidCircuit(fl validator.FieldLevel) bool {
	circuit, err := s.GetByName(fl.Field().String())
	if err != nil {
		return false
	}
	return circuit.Enabled
}
