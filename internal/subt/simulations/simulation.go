package simulations

import "gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"

type Simulation interface {
	simulations.Simulation
	Track() string
}
