package simulations

import "gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"

// Simulation groups a set of methods to identify a SubT simulation.
type Simulation interface {
	simulations.Simulation

	// GetTrack returns the track name of the simulation that will be used as the simulation world.
	GetTrack() string

	// GetToken returns the websocket access token for users to connect to a gazebo instance through GZ3D.
	GetToken() *string

	// GetRobots returns the list of robots from a certain simulation.
	GetRobots() []simulations.Robot

	// GetMarsupials returns the list of marsupials from a certain simulation.
	GetMarsupials() []simulations.Marsupial

	// GetWorldIndex returns the WorldIndex id.
	GetWorldIndex() int

	// GetName returns the simulation's name.
	GetName() string
}

// IsRobotChildMarsupial checks if the given robot is a child from the given slice of marsupials.
func IsRobotChildMarsupial(marsupials []simulations.Marsupial, robot simulations.Robot) bool {
	for _, m := range marsupials {
		if robot.IsEqual(m.GetChild()) {
			return true
		}
	}
	return false
}
