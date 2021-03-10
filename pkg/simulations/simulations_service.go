package simulations

// CreateSimulationInput contains all the information needed to create a simulation.
type CreateSimulationInput struct {
	Name      string
	Owner     string
	Creator   string
	Image     []string
	Private   bool
	StopOnEnd bool
	Extra     string
	Track     string
	Robots    string
}

// Service is a generic simulation service interface.
type Service interface {
	Create(input CreateSimulationInput) (Simulation, error)

	// Get returns a simulation with the given GroupID.
	Get(groupID GroupID) (Simulation, error)

	// GetParent returns the child simulation's parent with the given GroupID.
	GetParent(groupID GroupID) (Simulation, error)

	// UpdateStatus updates the simulation status with the given groupID.
	UpdateStatus(groupID GroupID, status Status) error

	// UpdateScore updates the simulation score.
	UpdateScore(groupID GroupID, score *float64) error

	// Update updates the simulation matching the given groupID with the given simulation data.
	Update(groupID GroupID, simulation Simulation) error

	// GetRobots returns the robot list of the simulation with the given GroupID.
	GetRobots(groupID GroupID) ([]Robot, error)

	// MarkStopped marks a simulation identified with the given Group ID as stopped.
	MarkStopped(groupID GroupID) error

	// GetWebsocketToken returns a websocket token for a certain simulation with the given GroupID.
	GetWebsocketToken(groupID GroupID) (string, error)
}
