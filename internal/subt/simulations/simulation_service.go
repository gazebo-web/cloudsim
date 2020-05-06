package simulations

import (
	"context"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/users"
	fuel "gitlab.com/ignitionrobotics/web/fuelserver/bundles/users"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

type IService interface {
	Parent() simulations.IService
	Create(ctx context.Context, createSimulation *SimulationCreate, user *fuel.User) (*Simulation, *ign.ErrMsg)
	CountByOwnerAndCircuit(owner, circuit string) (*int, error)
	simulationIsHeld(ctx context.Context, sim *simulations.Simulation) bool
	checkValidNumberOfSimulations(ctx context.Context, sim *simulations.Simulation) *ign.ErrMsg
}

type Service struct {
	parent simulations.IService
	userService users.IService
	repository IRepository
}

func NewService(repository IRepository) IService {
	var s IService
	parent := simulations.NewService(repository)
	s = &Service{
		parent: parent,
		repository: repository,
	}
	return s
}

func (s *Service) Parent() simulations.IService {
	return s.parent
}

func (s *Service) CountByOwnerAndCircuit(owner, circuit string) (*int, error) {
	panic("Not implemented")
}

func (s *Service) simulationIsHeld(ctx context.Context, sim *simulations.Simulation) bool {
}


func (s *Service) Create(ctx context.Context, createSimulation *SimulationCreate, user *fuel.User) (*Simulation, *ign.ErrMsg) {
	var sim *simulations.Simulation
	var em *ign.ErrMsg

	if sim, em = s.parent.Create(ctx, &createSimulation.SimulationCreate, user); em != nil {
		return nil, em
	}

	// Set held state if the user is not a sysadmin and the simulations needs to be held
	if !s.userService.IsSystemAdmin(*user.Username) && s.simulationIsHeld(ctx, sim) {
		held := true
		simUpdate := simulations.SimulationUpdate{
			Held: &held,
		}
		sim, em = s.parent.Update(ctx, *sim.GroupID, simUpdate)
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
		s.parent.Reject(ctx, sim)
		return nil, em
	}
	subtSim, err := s.repository.CreateAggregated(sim)
	if err != nil {
		return nil, ign.NewErrorMessageWithBase(ign.ErrorDbSave, err)
	}
	return subtSim, nil
}

func (s *Service) checkValidNumberOfSimulations(ctx context.Context, sim *simulations.Simulation) *ign.ErrMsg {
	panic("implement me")
}