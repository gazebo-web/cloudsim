package factory

import (
	"github.com/gazebo-web/cloudsim/pkg/factory"
	"github.com/gazebo-web/cloudsim/pkg/validate"
)

// APIConfig contains configuration values used to initialize a Kubernetes API.
type APIConfig struct {
	// KubeConfig contains the path to the target cluster's kubeconfig file.
	KubeConfig string `yaml:"kubeconfig"`
}

// Components contains factory configs used to create cluster components.
type Components struct {
	// Nodes contains configuration to instance a Nodes implementation using a factory.
	Nodes *factory.Config `yaml:"nodes"`

	// Pods is a configuration to instance a Pods implementation.
	Pods *factory.Config `yaml:"pods"`

	// Ingresses is a configuration to instance a Ingresses implementation.
	Ingresses *factory.Config `yaml:"ingresses"`

	// IngressRules is a configuration to instance a IngressRules implementation.
	IngressRules *factory.Config `yaml:"ingressRules"`

	// Services is a configuration to instance a Services implementation.
	Services *factory.Config `yaml:"services"`

	// NetworkPolicies is a configuration to instance a NetworkPolicies implementation.
	NetworkPolicies *factory.Config `yaml:"networkPolicies"`

	// Configurations is a configuration to instance a Configurations implementation.
	Configurations *factory.Config `yaml:"configurations"`
}

// Config is used to create a Kubernetes cluster component.
type Config struct {
	// API contains configuration values used to initialize a Kubernetes API.
	API APIConfig

	// Components contains configuration information for different Cluster components.
	Components Components
}

// Validate validates that the config values are valid.
func (c *Config) Validate() error {
	return validate.DefaultStructValidator(c)
}
