package simulator

import "time"

const (
	NODE_NAME_SERVER = "gzserver-container"
	NODE_NAME_BRIDGE = "comms-bridge"
	NODE_NAME_FIELD_COMPUTER = "field-computer"
	NODE_NAME_SIDECAR = "copy-to-s3"
)

// Node represents a machine that will be used to run a simulation
type Node struct {
	ID        uint       `gorm:"primary_key" json:"-"`
	CreatedAt time.Time  `gorm:"type:timestamp(3) NULL" json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"type:timestamp(2) NULL" sql:"index" json:"-"`
	InstanceID      *string `json:"instance_id" gorm:"not null;unique"`
	LastKnownStatus *string `json:"status,omitempty"`
	GroupID *string `json:"group_id"`
	Application *string `json:"application,omitempty"`
}
