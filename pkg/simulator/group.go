package simulator

import "github.com/jinzhu/gorm"

const (
	GROUP_STATUS_INITIALIZING = iota + 1000
	GROUP_STATUS_INITIALIZED
)

const (
	GROUP_STATUS_LAUNCHING = iota + 2000
	GROUP_STATUS_LAUNCHED
)

const (
	GROUP_STATUS_STOPPING = iota + 3000
	GROUP_STATUS_STOPPED
)

const (
	GROUP_STATUS_RESTARTING = iota + 4000
	GROUP_STATUS_RESTARTED
)

const (
	GROUP_STATUS_DELETING = iota + 5000
	GROUP_STATUS_DELETED
)

// Group represents a set of nodes working together to run a simulation.
type Group struct {
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

func (c *Group) Launch() {
	c.Status = GROUP_STATUS_LAUNCHING
}

func (c *Group) Stop() {
	c.Status = GROUP_STATUS_STOPPING
}

func (c *Group) Delete() {
	c.Status = GROUP_STATUS_DELETING
}