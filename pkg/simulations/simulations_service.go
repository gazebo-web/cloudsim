package simulations

// Service is a generic simulation service interface.
type Service interface {
	// Get returns a simulation with the given groupID.
	Get(groupID GroupID) (Simulation, error)

	// GetParent returns the child simulation's parent with the given groupID.
	GetParent(groupID GroupID) (Simulation, error)

	// UpdateStatus updates the simulation status with the given groupID.
	UpdateStatus(groupID GroupID, status Status) error

	// Update updates the simulation matching the given groupID with the given simulation data.
	Update(groupID GroupID, simulation Simulation) error

	// GetRobots returns the robot list of the simulation with the given groupID.
	GetRobots(groupID GroupID) ([]Robot, error)
}
