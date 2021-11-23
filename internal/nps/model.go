package nps

// This file contains database tables used to by the application to manage
// simulations.

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"time"
)

// Simulation represents the simulation that will be launched in the cloud.
type Simulation struct {
	// Override default GORM Model fields
	ID        uint       `gorm:"primary_key" json:"-"`
	CreatedAt time.Time  `gorm:"type:timestamp(3) NULL" json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"type:timestamp(2) NULL" sql:"index" json:"-"`
	// Timestamp in which this simulation was stopped/terminated.
	StoppedAt *time.Time `gorm:"type:timestamp(3) NULL" json:"stopped_at,omitempty"`
	Owner     string     `json:"owner"`

	Name    string `json:"name"`
	GroupID string `json:"groupid"`
	Status  string `json:"status"`

	// The docker to run
	Image string `json:"image"`

	// Comma separated list of arguments to pass into the docker image
	Args string `json:"args"`
	URI  string `json:"uri"`
	IP   string `json:"ip"`
}

// Simulations is a slice of Simulation
type Simulations []Simulation

// TableName returns the Simulation database table name.
func (s *Simulation) TableName() string {
	return "simulations"
}

// SingularName defines the singular name of the entity represented by Simulation.
func (s *Simulation) SingularName() string {
	return "simulation"
}

// PluralName defines the plural name of the entity represented by Simulation.
func (s *Simulation) PluralName() string {
	return "simulations"
}

// GetGroupID returns the simulation's group ID.
func (s *Simulation) GetGroupID() simulations.GroupID {
	panic("implement me")
}

// GetStatus returns the simulation's current status.
func (s *Simulation) GetStatus() simulations.Status {
	panic("implement me")
}

// HasStatus checks if the simulation is in a specific status.
func (s *Simulation) HasStatus(status simulations.Status) bool {
	panic("implement me")
}

// SetStatus set the simulation status.
func (s *Simulation) SetStatus(status simulations.Status) {
	panic("implement me")
}

// GetKind returns the simulation kind.
// Currently the following kinds are available:
// * Single simulation.
// * Multisimulation parent.
// * Multisimulation child.
func (s *Simulation) GetKind() simulations.Kind {
	panic("implement me")
}

// IsKind checks that the simulation is of a specific kind.
func (s *Simulation) IsKind(kind simulations.Kind) bool {
	panic("implement me")
}

// GetError returns the simulation's registered error.
// Is returns `nil` if the simulation has no error.
func (s *Simulation) GetError() *simulations.Error {
	panic("implement me")
}

// GetImage returns the simulator image.
func (s *Simulation) GetImage() string {
	panic("implement me")
}

// GetValidFor returns amount of wall-clock time a simulation can run for.
// This value is used to verify that a simulation has expired.
func (s *Simulation) GetValidFor() time.Duration {
	panic("implement me")
}

// GetCreator returns the creater (typically a user) that requested the simulation.
func (s *Simulation) GetCreator() string {
	panic("implement me")
}

// GetOwner returns the owner (typically an organization) that request the simulation.
func (s *Simulation) GetOwner() *string {
	panic("implement me")
}

// IsProcessed indicates if the simulation has been post-processed after being marked as finished.
// This value is used to prevent simulations from being processed multiple times.
func (s *Simulation) IsProcessed() bool {
	panic("implement me")
}

// NewSimulation returns a Simulation type object.
func NewSimulation() simulations.Simulation {
	return &Simulation{}
}

// RegisteredUser represents a user that can launch simulations.
type RegisteredUser struct {
	// Override default GORM Model fields
	ID        uint       `gorm:"primary_key" json:"-"`
	CreatedAt time.Time  `gorm:"type:timestamp(3) NULL" json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"type:timestamp(2) NULL" sql:"index" json:"-"`

	Username *string `gorm:"not null;unique" json:"username,omitempty" validate:"required,min=3,alphanum,notinblacklist"`

	// SimulationLimit is the number of allowed simultaneous simulations.
	// A negative number indicates unilimited simulations.
	SimulationLimit int `json:"simulation_limit"`
}

// RegisteredUsers is a slice of RegisteredUser.
type RegisteredUsers []RegisteredUser

// TableName returns the RegisteredUser database table name.
func (r *RegisteredUser) TableName() string {
	return "registered_users"
}

// SingularName defines the singular name of the entity represented by RegisteredUser.
func (r *RegisteredUser) SingularName() string {
	return "registered_user"
}

// PluralName defines the plural name of the entity represented by RegisteredUser.
func (r *RegisteredUser) PluralName() string {
	return "registered_users"
}
