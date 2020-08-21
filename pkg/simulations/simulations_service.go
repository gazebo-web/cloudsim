package simulations

type Service interface {
	Get(groupID GroupID) (Simulation, error)
	Reject(groupID GroupID) (Simulation, error)
	GetParent(gid GroupID) (Simulation, error)
}
