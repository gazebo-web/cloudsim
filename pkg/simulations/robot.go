package simulations

// Robot represents a generic robot used in a Simulation.
type Robot interface {
	// Name returns the robot's name. It's usually provided by the user.
	Name() string
	// Type returns the robot's type. It's the robot config type.
	Type() string
}
