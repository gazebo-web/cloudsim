package simulations

// GroupID is an universally unique identifier that helps to identify a Simulation.
type GroupID string

// Status defines the latest stage that a Simulation has reached.
type Status string

var (
	// StatusPending is used when a simulation is pending to be scheduled.
	StatusPending Status = "pending"

	// StatusRunning is used when a simulation is running.
	StatusRunning Status = "running"
)

// Kind is used to identify if a Simulation is a single simulation or a multisim.
// If a simulation is a multisim, different Kind values are used
// to identify if the simulations is a parent simulation or a child simulation.
type Kind uint

var (
	// SimSingle represents a single simulation.
	SimSingle Kind = 0
	// SimParent represents a parent simulation.
	SimParent Kind = 1
	// SimChild represents a child simulation.
	SimChild Kind = 2
)

// Simulation groups a set of methods to identify a simulation.
type Simulation interface {
	GroupID() GroupID
	// Status returns the current simulation's status.
	Status() Status
	// Kind returns the simulation's kind.
	Kind() Kind
}
