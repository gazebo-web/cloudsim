package groups

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator/nodes"
)

const (
	StatusGroupInitializing = iota + 1000
	StatusGroupInitialized
)

const (
	StatusGroupLaunching = iota + 2000
	StatusGroupLaunched
)

const (
	StatusGroupStopping = iota + 3000
	StatusGroupStopped
)

const (
	StatusGroupRestarting = iota + 4000
	StatusGroupRestarted
)

const (
	StatusGroupDeleting = iota + 5000
	StatusGroupDeleted
)

// Group represents a set of nodes working together to run a simulation.
type Group struct {
	gorm.Model
	Name        string `json:"name"`
	GroupID     string `json:"group_id"`
	Platform    string `json:"platform"`
	Application string `json:"application"`
	Status      int64  `json:"status"`
	PrivateKey  string `json:"private_key"`
	IAM         string `json:"iam"`
	Region      string `json:"region"`
	Zone        string `json:"zone"`
	Nodes       []nodes.Node
}

func (Group) TableName() string {
	return "cloudsim_simulator_groups"
}
