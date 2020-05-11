package simulations

import "gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"

type SimulationCreateInput interface {
	simulations.SimulationCreateInput
	ChildInput() SimulationCreate
}

type SimulationCreate struct {
	*simulations.SimulationCreate
}

func (sc *SimulationCreate) Input() *simulations.SimulationCreate {
	return sc.SimulationCreate
}

func (sc *SimulationCreate) ChildInput() *SimulationCreate {
	return sc
}

type SimulationCreatePersistentInput interface {
	simulations.SimulationCreatePersistentInput
	ChildInput()
}
