package simulator

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// Starter groups a set of methods of different jobs that will start a simulation.
type Starter interface {
	// HasStatus checks if the given simulation has a certain status.
	HasStatus(groupID simulations.GroupID, status simulations.Status) (bool, error)

	// IsParent checks that the given simulation is a parent simulation.
	IsParent(groupID simulations.GroupID) (bool, error)

	// IsChild checks that the given simulation is a child simulation.
	IsChild(groupID simulations.GroupID) (bool, error)

	// SetStatus sets the given status to the given simulation.
	SetStatus(groupID simulations.GroupID, status simulations.Status) error

	// CreateNodes creates a certain amount of machines for the given simulation.
	// It blocks the execution until the nodes are created.
	CreateNodes(groupID simulations.GroupID) error

	// LaunchFieldComputer launches a simulation's field computer pod.
	LaunchFieldComputer(groupID simulations.GroupID) error

	// LaunchCommsBridge launches a simulation's comms bridge pod.
	LaunchCommsBridge(groupID simulations.GroupID) error

	// WaitForSimulationPods waits for the simulation's pods are ready.
	// It blocks the execution until the pods are ready.
	WaitForSimulationPods(groupID simulations.GroupID) error
}

// start is a Starter implementation.
type start struct {
	Orchestrator orchestrator.Cluster
	Machines     cloud.Machines
	Storage      cloud.Storage
	Simulations  simulations.Service
}

func (s start) HasStatus(groupID simulations.GroupID, status simulations.Status) (bool, error) {
	panic("implement me")
}

func (s start) IsParent(groupID simulations.GroupID) (bool, error) {
	panic("implement me")
}

func (s start) IsChild(groupID simulations.GroupID) (bool, error) {
	panic("implement me")
}

func (s start) SetStatus(groupID simulations.GroupID, status simulations.Status) error {
	panic("implement me")
}

func (s start) CreateNodes(groupID simulations.GroupID) error {
	panic("implement me")
}

func (s start) LaunchFieldComputer(groupID simulations.GroupID) error {
	panic("implement me")
}

func (s start) LaunchCommsBridge(groupID simulations.GroupID) error {
	panic("implement me")
}

func (s start) WaitForSimulationPods(groupID simulations.GroupID) error {
	panic("implement me")
}

func NewSimulatorStarter(cluster orchestrator.Cluster, machines cloud.Machines, storage cloud.Storage, simulationService simulations.Service) Starter {
	return &start{
		Orchestrator: cluster,
		Machines:     machines,
		Storage:      storage,
		Simulations:  simulationService,
	}
}
