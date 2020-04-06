package simulator

import (
	"github.com/jinzhu/gorm"
)

const (
	NODE_NAME_SERVER         = "gzserver-container"
	NODE_NAME_BRIDGE         = "comms-bridge"
	NODE_NAME_FIELD_COMPUTER = "field-computer"
	NODE_NAME_SIDECAR        = "copy-to-s3"
)

// Node represents a machine that will be used to run a simulation
type Node struct {
	gorm.Model
	InstanceID     string   `json:"instance_id" gorm:"not null;unique"`
	Type           string   `json:"type"`
	Subnet         string   `json:"subnet"`
	SecurityGroups []string `json:"security_groups"`
	Status         string   `json:"status"`
}

func (Node) TableName() string {
	return "cloudsim_simulator_nodes"
}
