package simulator

import "github.com/jinzhu/gorm"

const (
	//////////////////////////////////////////////////////////////////
	// INIT
	CLUSTER_STATUS_INITIALIZING = iota + 1000
	CLUSTER_STATUS_INITIALIZED
	//////////////////////////////////////////////////////////////////
	// LAUNCH
	CLUSTER_STATUS_LAUNCHING = iota + 2000
	CLUSTER_STATUS_LAUNCHED
	//////////////////////////////////////////////////////////////////
	// STOP
	CLUSTER_STATUS_STOPPING = iota + 3000
	CLUSTER_STATUS_STOPPED
	//////////////////////////////////////////////////////////////////
	// RESTART
	CLUSTER_STATUS_RESTARTING = iota + 4000
	CLUSTER_STATUS_RESTARTED
	//////////////////////////////////////////////////////////////////
	// DELETE
	CLUSTER_STATUS_DELETING = iota + 5000
	CLUSTER_STATUS_DELETED
)

// Cluster represents a set of nodes working together to run a simulation.
type Cluster struct {
	gorm.Model
	Name string `json:"name"`
	GroupID string `json:"group_id"`
	Platform string `json:"platform"`
	Application string `json:"application"`
	Status int64 `json:"status"`
	PrivateKey string `json:"private_key"`
	Iam string `json:"iam"`
	Region string `json:"region"`
	Zone string `json:"zone"`
	Nodes []Node
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