package simulations

import "gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"

// Simulation groups a set of methods to identify a SubT simulation.
type Simulation interface {
	simulations.Simulation

	// Track returns the track name of the simulation that will be used as the simulation world.
	Track() string
	Token() *string
	Robots() []simulations.Robot
	Marsupials() []simulations.Marsupial
}
