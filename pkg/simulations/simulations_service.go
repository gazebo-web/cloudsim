package simulations

// Service is a generic simulation service interface.
type Service interface {
	// Get returns a simulation with the given groupID.
	Get(groupID GroupID) (Simulation, error)
	// Reject rejects a simulation with the given groupID.
	Reject(groupID GroupID) (Simulation, error)
	// GetParent returns the child simulation's parent with the given groupID.
	GetParent(gid GroupID) (Simulation, error)
}
