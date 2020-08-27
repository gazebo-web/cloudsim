package simulations

// Robot represents a generic robot used in a Simulation.
type Robot interface {
	// Name returns the robot's name. It's usually provided by the user.
	Name() string
}
