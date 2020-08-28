package store

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"time"
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
	SubnetAndZone() (string, string)
	Tags(simulation simulations.Simulation, nodeType string, nameSuffix string) []cloud.Tag
	InitScript() *string
	BaseImage() string
	Timeout() time.Duration
	PollFrequency() time.Duration
	Limit() int
	NamePrefix() string
}