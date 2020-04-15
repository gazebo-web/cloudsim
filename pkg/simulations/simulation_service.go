package simulations

import (
	"context"
	"fmt"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/logger"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

type IService interface {
	GetRepository() IRepository
	SetRepository(repository IRepository)
	Get(groupID string) (*Simulation, error)
	GetAll() []Simulation
	Launch(ctx context.Context, simulation *Simulation) *ign.ErrMsg
	ValidateLaunch(ctx context.Context, simulation *Simulation) *ign.ErrMsg
	GetAllByOwner(owner string, application string, statusFrom, statusTo Status) (*Simulations, error)
	GetChildren(groupID string, application string, statusFrom, statusTo Status) (*Simulations, error)
	GetAllParents(application string, statusFrom, statusTo Status) (*Simulations, error)
	Update(groupID string, simulation Simulation) (*Simulation, error)
}

type Service struct {
	repository IRepository
}

func NewService(repository IRepository) IService {
	var s IService
	s = &Service{repository: repository}
	return s
}

func (s *Service) GetRepository() IRepository {
	return s.repository
}

func (s *Service) SetRepository(repository IRepository) {
	s.repository = repository
}

// LaunchSimulation -- sim_service.go:763
func (s *Service) Launch(ctx context.Context, simulation *Simulation) *ign.ErrMsg {
	if err := s.ValidateLaunch(ctx, simulation); err != nil {
		return err
	}
	return nil
}

// ValidateLaunch -- subt_specifics.go:2089
func (s *Service) ValidateLaunch(ctx context.Context, simulation *Simulation) *ign.ErrMsg {
	if err := s.isSimulationHeld(ctx, simulation); err != nil {
		logger.Logger(ctx).Warning(fmt.Sprintf("[LAUNCH|VALIDATE] Cannot run a held simulation. Group ID: [%s]", *simulation.GroupID))
		return err
	}
	return nil
}

func (s *Service) isSimulationHeld(ctx context.Context, simulation *Simulation) *ign.ErrMsg {
	if simulation.Held {
		// TODO: Replace ign.NewErrorMessage(ign.ErrorInvalidSimulationStatus) with ErrorLaunchHeld error.
		return ign.NewErrorMessage(ign.ErrorInvalidSimulationStatus)
	}
	return nil
}

func (s *Service) Update(groupID string, simulation Simulation) (*Simulation, error) {
	s.repository.Update(groupID, simulation)
}

func (s *Service) Get(groupID string) (*Simulation, error) {

}

func (s *Service) GetAll() []Simulation {

}

func (s *Service) GetAllByOwner(owner string, application string, statusFrom, statusTo Status) (*Simulations, error) {

}

func (s *Service) GetChildren(groupID string, application string, statusFrom, statusTo Status) (*Simulations, error) {

}

func (s *Service) GetAllParents(application string, statusFrom, statusTo Status) (*Simulations, error) {

}