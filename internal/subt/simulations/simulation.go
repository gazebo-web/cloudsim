package simulations

import "gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"

// Simulation groups a set of methods to identify a SubT simulation.
type Simulation interface {
	simulations.Simulation

	// Track returns the track name of the simulation that will be used as the simulation world.
	Track() string

	// Token returns the websocket access token for users to connect to a gazebo instance through GZ3D.
	Token() *string

	// Robots returns the list of robots from a certain simulation.
	Robots() []simulations.Robot

	// Marsupials returns the list of marsupials from a certain simulation.
	Marsupials() []simulations.Marsupial
}

func IsRobotChildMarsupial(marsupials []simulations.Marsupial, robot simulations.Robot) bool {
	for _, m := range marsupials {
		if robot.IsEqual(m.Child()) {
			return true
		}
	}
	return false
}
