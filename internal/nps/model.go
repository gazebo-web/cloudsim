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
type Simulations []Simulation

func (s *Simulation) TableName() string {
	return "simulations"
}

func (s *Simulation) SingularName() string {
	return "simulation"
}

func (s *Simulation) PluralName() string {
	return "simulations"
}

func (s *Simulation) GetGroupID() simulations.GroupID {
	panic("implement me")
}

func (s *Simulation) GetStatus() simulations.Status {
	panic("implement me")
}

func (s *Simulation) HasStatus(status simulations.Status) bool {
	panic("implement me")
}

func (s *Simulation) SetStatus(status simulations.Status) {
	panic("implement me")
}

func (s *Simulation) GetKind() simulations.Kind {
	panic("implement me")
}

func (s *Simulation) IsKind(kind simulations.Kind) bool {
	panic("implement me")
}

func (s *Simulation) GetError() *simulations.Error {
	panic("implement me")
}

func (s *Simulation) GetImage() string {
	panic("implement me")
}

func (s *Simulation) GetValidFor() time.Duration {
	panic("implement me")
}

func (s *Simulation) GetCreator() string {
	panic("implement me")
}

func (s *Simulation) GetOwner() *string {
	panic("implement me")
}

func (s *Simulation) IsProcessed() bool {
	panic("implement me")
}
func NewSimulation() simulations.Simulation {
	return &Simulation{}
}

// Simulation represents the simulation that will be launched in the cloud.
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
type RegisteredUsers []RegisteredUser

func (r *RegisteredUser) TableName() string {
	return "registered_users"
}

func (r *RegisteredUser) SingularName() string {
	return "registered_user"
}

func (r *RegisteredUser) PluralName() string {
	return "registered_users"
}