package simulator

import (
	"context"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator"
)

type subTSimulator struct {
	platform          platform.Platform
	simulationService simulations.Service
}

// Start returns an action that will be in charge of launching a simulation with the given Group ID.
func (s *subTSimulator) Start(ctx context.Context) (*actions.Action, error) {
	jobs := actions.Jobs{
		JobCheckPendingStatus,
		JobCheckSimulationParenthood,
		JobCheckParentSimulationWithError,
	}
	action, err := actions.NewAction(jobs)
	if err != nil {
		return nil, err
	}
	return action, nil
}

// Stop returns an action that will be in charge of stopping a simulation with the given Group ID.
func (s *subTSimulator) Stop(ctx context.Context) (*actions.Action, error) {
	panic("implement me")
}

// Restart returns an action that will be in charge of restarting a simulation with the given Group ID.
func (s *subTSimulator) Restart(ctx context.Context) (*actions.Action, error) {
	panic("implement me")
}

func NewSimulator(platform platform.Platform) simulator.Simulator {
	return &subTSimulator{
		platform: platform,
	}
}
