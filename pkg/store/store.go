package store

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
)

// Store provides a set of components to store data that needs to be accessed by different services.
type Store interface {
	Machines() Machines
}

// Machines provides different information for creating machines.
type Machines interface {
	InstanceProfile() *string
	KeyName() string
	Type() string
	FirewallRules() []string
	Subnet() string
	Zone() string
	Tags(simulation simulations.Simulation, nodeType string) []cloud.Tag
	InitScript() *string
	BaseImage() string
}
