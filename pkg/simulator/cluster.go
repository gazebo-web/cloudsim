package simulator

type Status int64

const (
	CLUSTER_STATUS_INITIALIZING = iota + 1000
	CLUSTER_STATUS_INITIALIZED
	CLUSTER_STATUS_LAUNCHING
	CLUSTER_STATUS_LAUNCHED
	CLUSTER_STATUS_STOPPING
	CLUSTER_STATUS_STOPPED
	CLUSTER_STATUS_RESTARTING
	CLUSTER_STATUS_RESTARTED
	CLUSTER_STATUS_DELETING
	CLUSTER_STATUS_DELETED
)

// Cluster represents a set of nodes working together to run a simulation.
type Cluster struct {
	Name string
	Status Status
	Nodes []Node
}

func (c Cluster) ListNodes() []Node {
	return c.Nodes
}

func (c *Cluster) Launch() {
	c.Status = CLUSTER_STATUS_LAUNCHING
}

func (c *Cluster) Stop() {
	c.Status = CLUSTER_STATUS_STOPPING
}

func (c *Cluster) Delete() {
	c.Status = CLUSTER_STATUS_DELETING
}