package simulations

import (
	"context"
	"fmt"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/logger"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

type IService interface {
	simulations.IService
}

type Service struct {
	*simulations.Service
}

func (s *Service) ValidateLaunch(ctx context.Context, simulation *simulations.Simulation) *ign.ErrMsg {
	if err := s.isHeld(ctx, simulation); err != nil {
		return err
	}
	return nil
}

func (s *Service) isHeld(ctx context.Context, simulation *simulations.Simulation) *ign.ErrMsg {
	if simulation.Held {
		logger.Logger(ctx).Warning(fmt.Sprintf("[LAUNCH|VALIDATE] Cannot run a held simulation. Group ID: [%s]", *simulation.GroupID))
		return nil
	}
	return nil
}