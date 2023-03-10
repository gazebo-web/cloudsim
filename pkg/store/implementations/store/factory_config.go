package store

import (
	"github.com/gazebo-web/gz-go/v7/defaults"
	"github.com/gazebo-web/gz-go/v7/validate"
)

// Config is used to create a store component.
type Config struct {
	MachinesStore     machinesStore
	IgnitionStore     ignitionStore
	OrchestratorStore orchestratorStore
}

// Validate validates that the config values are valid.
func (c *Config) Validate() error {
	return validate.DefaultStructValidator(c)
}

// SetDefaults sets defaults values for the config.
func (c *Config) SetDefaults() error {
	return defaults.SetStructValues(c)
}
