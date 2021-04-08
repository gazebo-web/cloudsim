package factory

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/factory"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/validate"
)

// Config is used to create a platform.
type Config struct {
	// Machines contains configuration to instance a Machines implementation using a factory.
	Machines *factory.Config `yaml:"machines"`
	// Storage contains configuration to instance a Storage implementation using a factory.
	Storage *factory.Config `yaml:"storage"`
	// Orchestrator contains configuration to instance an Orchestrator implementation using a factory.
	Orchestrator *factory.Config `yaml:"orchestrator"`
	// Store contains configuration to instance a Store implementation using a factory.
	Store *factory.Config `yaml:"store"`
	// Secrets contains configuration to instance a Secrets implementation using a factory.
	Secrets *factory.Config `yaml:"secrets"`
	// EmailSender contains configuration to instance a EmailSender implementation using a factory.
	EmailSender *factory.Config `yaml:"emailSender"`
	// RunningSimulations contains configuration to instance a RunningSimulations implementation using a factory.
	RunningSimulations *factory.Config `yaml:"runningSimulations"`
}

// Validate validates that the config values are valid.
func (c *Config) Validate() error {
	return validate.DefaultStructValidator(c)
}
