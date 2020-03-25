package simulation

import "time"

type Simulation struct {
	// Override default GORM Model fields
	ID        uint      `gorm:"primary_key" json:"-"`
	CreatedAt time.Time `gorm:"type:timestamp(3) NULL" json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// Added 2 milliseconds to DeletedAt field
	DeletedAt *time.Time `gorm:"type:timestamp(2) NULL" sql:"index" json:"-"`
	// Timestamp in which this simulation was stopped/terminated.
	StoppedAt *time.Time `gorm:"type:timestamp(3) NULL" json:"stopped_at,omitempty"`
	// Represents the maximum time this simulation should live. After that time
	// it will be eligible for automatic termination.
	// It is a time.Duration (stored as its string representation).
	ValidFor *string `json:"valid_for,omitempty"`
	// The owner of this deployment (must exist in UniqueOwners). Can be user or org.
	// Also added to the name_owner unique index
	Owner *string `json:"owner,omitempty"`
	// The username of the User that created this resource (usually got from the JWT)
	Creator *string `json:"creator,omitempty"`
	// Private - True to make this a private resource
	Private *bool `json:"private,omitempty"`
	// When shutting down simulations, stop EC2 instances instead of terminating them. Requires admin privileges.
	StopOnEnd *bool `json:"stop_on_end"`
	// The user defined Name for the simulation.
	Name *string `json:"name,omitempty"`
	// The docker image url to use for the simulation (usually for the Field Computer)
	Image *string `json:"image,omitempty" form:"image"`
	// GroupID - Simulation Unique identifier
	// All k8 pods and services (or other created resources) will share this groupID
	GroupID *string `gorm:"not null;unique" json:"group_id"`
	// ParentGroupID (optional) holds the GroupID of the parent simulation record.
	// It is used with requests for multi simulations (multiSims), where a single
	// user request spawns multiple simulation runs based on a single template.
	ParentGroupID *string `json:"parent"`
	// MultiSim holds which role this simulation plays within a multiSim deployment.
	// Values should be of type MultiSimType.
	MultiSim int
	// A value from DeploymentStatus constants
	DeploymentStatus *int `json:"status,omitempty"`
	// A value from ErrorStatus constants
	ErrorStatus *string `json:"error_status,omitempty"`
	// NOTE: statuses should be updated in sequential DB Transactions. ie. one status per TX.
	Platform    *string `json:"platform,omitempty" form:"platform"`
	Application *string `json:"application,omitempty" form:"application"`
	// TODO: both fields Extra and ExtraSelector should be a separate table, specific to
	//   each Application, and with a reference back to this simulation ID.
	// A free form string field to store extra details, usually associated to the
	// chosen Application. Eg. SubT would store here the different robot names, types
	// and images.
	Extra *string `gorm:"size:999999" json:"extra,omitempty"`
	// A extra string field to store a selector that can help specific Applications
	// to filter simulations (eg. SQL WHERE). SubT could store the circuit here.
	ExtraSelector *string `json:"-"`
	// TODO: This is a field specific to SubT. As such this is a temporary field
	//  that should be included in the same separate table where Extra and
	//  ExtraSelector should reside.
	// Contains the names of all robots in the simulation in a comma-separated list.
	Robots *string `gorm:"size:1000" json:"robots"`
	// TODO: This is a field specific to SubT. This is a temporary field that should be
	//  extracted from the Simulation struct.
	Held bool `json:"held"`
}