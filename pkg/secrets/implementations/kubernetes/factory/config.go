package factory

import (
	"github.com/gazebo-web/gz-go/v7/validate"
)

// APIConfig contains configuration values used to initialize a Kubernetes API.
type APIConfig struct {
	// KubeConfig contains the path to the target cluster's kubeconfig file.
	KubeConfig string `yaml:"kubeconfig"`
}

// Config is used to create a Kubernetes Secrets component.
type Config struct {
	// API contains config
	API APIConfig
}

// Validate validates that the config values are valid.
func (c *Config) Validate() error {
	return validate.DefaultStructValidator(c)
}
