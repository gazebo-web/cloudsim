package rules

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/users"
	"gitlab.com/ignitionrobotics/web/cloudsim/tools"
	fuel "gitlab.com/ignitionrobotics/web/fuelserver/bundles/users"
	per "gitlab.com/ignitionrobotics/web/fuelserver/permissions"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"strconv"
)

type Service interface {
	GetByCircuitAndOwner(ruleType Type, circuit, owner string) (*Rule, error)
	GetRemainingSubmissions(owner, circuit string, user fuel.User) (*int, *ign.ErrMsg)
}

type service struct {
	repository Repository
	users users.Service
	simulation simulations.IService
}

func (s *service) GetByCircuitAndOwner(ruleType Type, circuit, owner string) (*Rule, error) {
	return s.repository.GetByCircuitAndOwner(ruleType, circuit, owner)
}

func (s *service) GetRemainingSubmissions(owner, circuit string, user fuel.User) (*int, *ign.ErrMsg) {
	if owner != "" {
		ok, em := s.users.VerifyOwner(owner, *user.Username, per.Read)
		if !ok {
			return nil, em
		}
	}

	if owner == "" {
		owner = *user.Username
	}

	rule, err := s.GetByCircuitAndOwner(MaxSubmissions, circuit, owner)
	if err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorUnexpected, err)
	}

	limit, err := strconv.Atoi(rule.Value)
	if err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorUnexpected, err)
	}

	if limit == 0 {
		return nil, nil
	}

	count, err := s.simulation.CountByOwnerAndCircuit(owner, circuit)
	if err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorUnexpected, err)
	}

	remaining := tools.Max(0, limit - *count)

	return &remaining, nil
}

func NewService(repository Repository, user users.Service) Service {
	return &service{
		repository: repository,
		users: user,
	}
}