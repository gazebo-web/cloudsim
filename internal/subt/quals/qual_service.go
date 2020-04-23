package quals

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/circuits"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/users"
)

type IService interface {
	IsQualified(owner, circuit, username string) bool
}

type Service struct {
	services services
	repository IRepository
}

type services struct {
	User *users.Service
	Circuit circuits.IService
}

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