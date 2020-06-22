package quals

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/circuits"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/users"
)

type IService interface {
	IsQualified(owner, circuit, username string) bool
}

type Service struct {
	services   services
	repository IRepository
}

type NewServiceInput struct {
	Services   services
	Repository IRepository
}

func NewService(input NewServiceInput) IService {
	var s IService
	s = &Service{
		services: services{
			User:    input.Services.User,
			Circuit: input.Services.Circuit,
		},
		repository: input.Repository,
	}
	return s
}

// services represents the imported services used by the Qualification service.
type services struct {
	User    users.Service
	Circuit circuits.Service
}

// IsQualified returns true if the given owner was qualified for the given circuit.
// If the provided username is an admin, it will skip the qualified condition.
func (s *Service) IsQualified(owner, circuit, username string) bool {
	if s.services.User.IsSystemAdmin(username) {
		return true
	}
	var c *circuits.Circuit
	var err error
	if c, err = s.services.Circuit.GetByName(circuit); err != nil {
		return false
	}

	if c.RequiresQualification == nil && !(*c.RequiresQualification) {
		return true
	}

	_, err = s.repository.GetByOwnerAndCircuit(owner, circuit)
	return err == nil
}
