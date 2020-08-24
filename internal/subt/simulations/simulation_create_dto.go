package simulations

import "gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"

type ServiceCreateInput interface {
	simulations.ServiceCreateInput
	ChildInput() SimulationCreate
}

type SimulationCreate struct {
	*simulations.SimulationCreate
	Score *float64 `json:"score"`
	// Simulation run info
	SimTimeDurationSec  int `json:"sim_time_duration_sec"`
	RealTimeDurationSec int `json:"real_time_duration_sec"`
	ModelCount          int `json:"model_count"`
}

func (sc *SimulationCreate) Input() *simulations.SimulationCreate {
	return sc.SimulationCreate
}

func (sc *SimulationCreate) ChildInput() *SimulationCreate {
	return sc
}

type RepositoryCreateInput interface {
	simulations.RepositoryCreateInput
	ChildInput() *Simulation
}
