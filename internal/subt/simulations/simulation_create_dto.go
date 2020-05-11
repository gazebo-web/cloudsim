package simulations

import "gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"

type SimulationCreateInput interface {
	simulations.SimulationCreateInput
	Child() Simulation

}

type SimulationCreate struct {
	*simulations.SimulationCreate
}