package store

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/defaults"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/validate"
)

// Config is used to create a store component.
type Config struct {
	MachinesStore     machinesStore
	IgnitionStore     ignitionStore
	MoleStore		  moleStore
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
