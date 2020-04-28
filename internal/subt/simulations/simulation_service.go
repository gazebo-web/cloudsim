package simulations

import (
	"context"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/users"
	fuel "gitlab.com/ignitionrobotics/web/fuelserver/bundles/users"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

type IService interface {
	simulations.IService
	CountByOwnerAndCircuit(owner, circuit string) (*int, error)
	simulationIsHeld(ctx context.Context, sim *simulations.Simulation) bool
	checkValidNumberOfSimulations(ctx context.Context, sim *simulations.Simulation) *ign.ErrMsg
}

type Service struct {
	*simulations.Service
	userService users.IService
}

func NewService(repository IRepository) IService {
	var s IService
	parent := simulations.NewService(repository)
	service := parent.(*simulations.Service)
	s = &Service{
		Service: service,
	}
	return s
}

func (s *Service) CountByOwnerAndCircuit(owner, circuit string) (*int, error) {
	panic("Not implemented")
}

func (s *Service) simulationIsHeld(ctx context.Context, sim *simulations.Simulation) bool {
	panic("implement me")
}


func (s *Service) Create(ctx context.Context, createSimulation *simulations.SimulationCreate, user *fuel.User) (*simulations.Simulation, *ign.ErrMsg) {


	sim, err := s.Service.Create(ctx, createSimulation, user)

	// Set held state if the user is not a sysadmin and the simulations needs to be held
	if !s.userService.IsSystemAdmin(*user.Username) && s.simulationIsHeld(ctx, sim) {
		held := true
		simUpdate := simulations.SimulationUpdate{
			Held: &held,
		}
		sim, err = s.Update(ctx, *sim.GroupID, simUpdate)
		if err != nil {
			return nil, err
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

func (s *Service) checkValidNumberOfSimulations(ctx context.Context, sim *simulations.Simulation) *ign.ErrMsg {
	panic("implement me")
}