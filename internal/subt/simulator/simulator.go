package simulator

import (
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
func (s *subTSimulator) Start(groupID simulations.GroupID) (*actions.Action, error) {
	var jobs actions.Jobs
	jobs = append(jobs, s.createCheckPendingStatusJob(groupID))
	jobs = append(jobs, s.createCheckSimulationIsParentJob(groupID))
	jobs = append(jobs, s.createCheckParentSimulationWithErrorJob(groupID))
	action, err := actions.NewAction(jobs)
	if err != nil {
		return nil, err
	}
	return action, nil
}

// Stop returns an action that will be in charge of stopping a simulation with the given Group ID.
func (s *subTSimulator) Stop(groupID simulations.GroupID) (*actions.Action, error) {
	panic("implement me")
}

// Restart returns an action that will be in charge of restarting a simulation with the given Group ID.
func (s *subTSimulator) Restart(groupID simulations.GroupID) (*actions.Action, error) {
	panic("implement me")
}

func NewSimulator(platform platform.Platform) simulator.Simulator {
	return &subTSimulator{
		platform: platform,
	}
}
