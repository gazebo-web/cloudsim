package simulations

import (
	"context"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/metadata"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/users"
	fuel "gitlab.com/ignitionrobotics/web/fuelserver/bundles/users"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

type IService interface {
	simulations.Service
	CountByOwnerAndCircuit(owner, circuit string) (*int, error)
	simulationIsHeld(ctx context.Context, sim *simulations.Simulation) bool
	checkValidNumberOfSimulations(ctx context.Context, sim *simulations.Simulation) *ign.ErrMsg
}

type Service struct {
	simulations.Service
	userService users.Service
	repository  Repository
}

func NewService(repository Repository) IService {
	var s IService
	parent := simulations.NewService(simulations.NewServiceInput{
		Repository: repository,
		Config:     simulations.ServiceConfig{},
	})
	s = &Service{
		Service:    parent,
		repository: repository,
	}
	return s
}

func (s *Service) CountByOwnerAndCircuit(owner, circuit string) (*int, error) {
	panic("Not implemented")
}

func (s *Service) simulationIsHeld(ctx context.Context, sim *simulations.Simulation) bool {
	return false
}

func (s *Service) Create(ctx context.Context, input simulations.ServiceCreateInput, user *fuel.User) (simulations.ServiceCreateOutput, *ign.ErrMsg) {
	createSim := input.Input()

	var output simulations.ServiceCreateOutput
	var em *ign.ErrMsg

	createSubTSim := &SimulationCreate{
		SimulationCreate:    createSim,
		Score:               nil,
		SimTimeDurationSec:  0,
		RealTimeDurationSec: 0,
		ModelCount:          0,
	}

	if output, em = s.Service.Create(ctx, createSubTSim, user); em != nil {
		return nil, em
	}

	sim := output.Output()

	// Set held state if the user is not a sysadmin and the simulations needs to be held
	if !s.userService.IsSystemAdmin(*user.Username) && s.simulationIsHeld(ctx, sim) {
		held := true
		simUpdate := simulations.SimulationUpdate{
			Held: &held,
		}
		sim, em = s.Update(ctx, *sim.GroupID, simUpdate)
		if em != nil {
			return nil, em
		}
	}
	// Sanity check: check for maximum number of allowed simultaneous simulations per Owner.
	// Also allow Applications to provide custom validations.
	// Dev note: in this case we check 'after' creating the record in the DB to make
	// sure that in case of a race condition then both records are added with pending state
	// and one of those (or both) can be rejected immediately.
	if em := s.checkValidNumberOfSimulations(ctx, sim); em != nil {
		s.Reject(ctx, sim)
		return nil, em
	}

	return sim, nil
}

func (s *Service) Get(groupID string, user *fuel.User) (*simulations.Simulation, *ign.ErrMsg) {
	sim, em := s.Service.Get(groupID, user)
	if em != nil {
		return nil, em
	}

	extra, err := metadata.Read(sim)
	if err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorUnexpected, err)
	}

	// If the user is not a system admin, remove the RunIndex and WorldIndex fields.
	if ok := s.userService.IsSystemAdmin(*user.Username); !ok {
		extra.RunIndex = nil
		extra.WorldIndex = nil
	}

	sim.Extra, err = extra.ToJSON()
	if err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorUnexpected, err)
	}
	return sim, nil
}

func (s *Service) checkValidNumberOfSimulations(ctx context.Context, sim *simulations.Simulation) *ign.ErrMsg {
	panic("implement me")
}
