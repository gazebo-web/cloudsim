package starter

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator"
)

// start is a Starter implementation.
type start struct {
	Orchestrator orchestrator.Cluster
	Machines     cloud.Machines
	Storage      cloud.Storage
	Simulations  simulations.Service
}

// HasStatus returns true if the simulation with the given group id has the given status.
func (s *start) HasStatus(groupID simulations.GroupID, status simulations.Status) (bool, error) {
	sim, err := s.Simulations.Get(groupID)
	if err != nil {
		return false, err
	}
	return sim.Status() == status, nil
}

func (s *start) IsParent(groupID simulations.GroupID) (bool, error) {
	sim, err := s.Simulations.Get(groupID)
	if err != nil {
		return false, err
	}
	return sim.Kind() == simulations.SimParent, nil
}

func (s *start) IsChild(groupID simulations.GroupID) (bool, error) {
	panic("implement me")
}

func (s *start) SetStatus(groupID simulations.GroupID, status simulations.Status) error {
	panic("implement me")
}

func (s *start) CreateNodes(groupID simulations.GroupID) error {
	panic("implement me")
}

func (s *start) LaunchFieldComputer(groupID simulations.GroupID) error {
	panic("implement me")
}

func (s *start) LaunchCommsBridge(groupID simulations.GroupID) error {
	panic("implement me")
}

func (s *start) WaitForSimulationPods(groupID simulations.GroupID) error {
	panic("implement me")
}

func NewSimulatorStarter(cluster orchestrator.Cluster, machines cloud.Machines, storage cloud.Storage, simulationService simulations.Service) simulator.Starter {
	return &start{
		Orchestrator: cluster,
		Machines:     machines,
		Storage:      storage,
		Simulations:  simulationService,
	}
}
