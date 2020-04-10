package simulations

import (
	"context"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

type IService interface {
	Get(groupID string) (*Simulation, error)
	GetAll() []Simulation
	Launch(ctx context.Context, simulation *Simulation) *ign.ErrMsg
	ValidateLaunch(ctx context.Context, simulation *Simulation) *ign.ErrMsg
}

type Service struct {}

// LaunchSimulation -- sim_service.go:763
func (s *Service) Launch(ctx context.Context, simulation *Simulation) *ign.ErrMsg {
	if err := s.ValidateLaunch(ctx, simulation); err != nil {
		return err
	}
	return nil
}

// ValidateLaunch -- subt_specifics.go:2089
func (s *Service) ValidateLaunch(ctx context.Context, simulation *Simulation) *ign.ErrMsg {
	return nil
}