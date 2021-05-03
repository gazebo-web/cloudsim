package factory

import (
	gatewayV1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/client/clientset/versioned/typed/gateway.solo.io/v1"
	glooV1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/kube/client/clientset/versioned/typed/gloo.solo.io/v1"

	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/factory"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/ingresses/implementations/gloo"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/implementations/kubernetes/client"
)

// IngressesNewFunc is the factory creation function for the Gloo ingresses.Ingresses implementation.
func IngressesNewFunc(config interface{}, dependencies factory.Dependencies, out interface{}) error {
	// Parse config
	if config == nil {
		return factory.ErrNilConfig
	}
	var typeConfig Config
	if err := factory.SetValueAndValidate(&typeConfig, config); err != nil {
		return err
	}

	// Parse dependencies
	var typeDependencies Dependencies
	if err := dependencies.ToStruct(&typeDependencies); err != nil {
		return err
	}

	// Initialize dependencies
	dependenciesInitFns := []func(config *Config, dependencies *Dependencies) error{
		initializeGloo,
		initializeGlooGateway,
	}
	for _, initFn := range dependenciesInitFns {
		if err := initFn(&typeConfig, &typeDependencies); err != nil {
			return err
		}
	}

	// Create instance
	vs := gloo.NewVirtualServices(typeDependencies.GlooGateway, typeDependencies.Logger, typeDependencies.Gloo)
	if err := factory.SetValue(out, vs); err != nil {
		return err
	}

	return nil
}

// IngressRulesNewFunc is the factory creation function for the Gloo ingresses.IngressRules implementation.
func IngressRulesNewFunc(config interface{}, dependencies factory.Dependencies, out interface{}) error {
	// Parse config
	if config == nil {
		return factory.ErrNilConfig
	}
	var typeConfig Config
	if err := factory.SetValueAndValidate(&typeConfig, config); err != nil {
		return err
	}

	// Parse dependencies
	var typeDependencies Dependencies
	if err := dependencies.ToStruct(&typeDependencies); err != nil {
		return err
	}

	// Initialize dependencies
	dependenciesInitFns := []func(config *Config, dependencies *Dependencies) error{
		initializeGlooGateway,
	}
	for _, initFn := range dependenciesInitFns {
		if err := initFn(&typeConfig, &typeDependencies); err != nil {
			return err
		}
	}

	// Create instance
	vh := gloo.NewVirtualHosts(typeDependencies.GlooGateway, typeDependencies.Logger)
	if err := factory.SetValue(out, vh); err != nil {
		return err
	}

	return nil
}

// initializeGloo initializes creates a Gloo Clientset if it wasn't passed as a dependency.
func initializeGloo(config *Config, dependencies *Dependencies) error {
	// Validate inputs
	if dependencies.Gloo != nil {
		return nil
	}

	kubeconfig, err := client.GetConfig(config.API.KubeConfig)
	if err != nil {
		return err
	}

	gloo, err := glooV1.NewForConfig(kubeconfig)
	if err != nil {
		return err
	}

	if err = factory.SetValue(&dependencies.Gloo, gloo); err != nil {
		return err
	}

	return nil
}

// initializeGloo initializes creates a Gloo Gateway Clientset if it wasn't passed as a dependency.
func initializeGlooGateway(config *Config, dependencies *Dependencies) error {
	// Validate inputs
	if dependencies.GlooGateway != nil {
		return nil
	}

	kubeconfig, err := client.GetConfig(config.API.KubeConfig)
	if err != nil {
		return err
	}

	gw, err := gatewayV1.NewForConfig(kubeconfig)
	if err != nil {
		return err
	}

	if err = factory.SetValue(&dependencies.GlooGateway, gw); err != nil {
		return err
	}

	return nil
}
