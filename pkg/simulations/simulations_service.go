package simulations

type Service interface {
	Get(groupID GroupID) (Simulation, error)
}
