package factory

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/factory"
	ingressesImpl "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/ingresses/implementations"
	networkImpl "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/network/implementations"
	nodesImpl "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/nodes/implementations"
	podsImpl "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/pods/implementations"
	servicesImpl "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/services/implementations"
	kubernetesCluster "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/implementations/kubernetes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/implementations/kubernetes/client"
)

// NewFunc is the factory creation function for the Kubernetes orchestrator.Cluster implementation.
func NewFunc(config interface{}, dependencies factory.Dependencies, out interface{}) error {
	// Parse config
	var typeConfig Config
	if err := factory.SetValueAndValidate(&typeConfig, config); err != nil {
		return err
	}

	// Parse dependencies
	var err error
	if dependencies, err = dependencies.DeepCopy(); err != nil {
		return err
	}
	var typeDependencies Dependencies
	if err := dependencies.ToStruct(&typeDependencies); err != nil {
		return err
	}

	// Initialize dependencies
	dependenciesInitFns := []func(config *Config, dependencies factory.Dependencies, typeDependencies *Dependencies) error{
		initializeAPI,
	}
	for _, initFn := range dependenciesInitFns {
		if err := initFn(&typeConfig, dependencies, &typeDependencies); err != nil {
			return err
		}
	}

	// Create components
	components := kubernetesCluster.Config{}
	factoryCalls := factory.Calls{
		// Nodes
		{
			Factory:      nodesImpl.Factory,
			Config:       typeConfig.Components.Nodes,
			Dependencies: dependencies,
			Out:          &components.Nodes,
		},
		// Pods
		{
			Factory:      podsImpl.Factory,
			Config:       typeConfig.Components.Pods,
			Dependencies: dependencies,
			Out:          &components.Pods,
		},
		// Ingresses
		{
			Factory:      ingressesImpl.IngressesFactory,
			Config:       typeConfig.Components.Ingresses,
			Dependencies: dependencies,
			Out:          &components.Ingresses,
		},
		// IngressRules
		{
			Factory:      ingressesImpl.IngressRulesFactory,
			Config:       typeConfig.Components.IngressRules,
			Dependencies: dependencies,
			Out:          &components.IngressRules,
		},
		// Services
		{
			Factory:      servicesImpl.Factory,
			Config:       typeConfig.Components.Services,
			Dependencies: dependencies,
			Out:          &components.Services,
		},
		// NetworkPolicies
		{
			Factory:      networkImpl.Factory,
			Config:       typeConfig.Components.NetworkPolicies,
			Dependencies: dependencies,
			Out:          &components.NetworkPolicies,
		},
	}
	if err := factory.CallFactories(factoryCalls); err != nil {
		return err
	}

	// Set output value
	cluster := kubernetesCluster.NewCustomKubernetes(components)
	if err := factory.SetValue(out, cluster); err != nil {
		return err
	}

	return nil
}

// initializeAPI initializes the API dependency.
func initializeAPI(config *Config, dependencies factory.Dependencies, typeDependencies *Dependencies) error {
	if typeDependencies.API != nil {
		return nil
	}
	if config == nil {
		return factory.ErrNilConfig
	}

	// Get the Kubernetes config
	kubeconfig, err := client.GetConfig(config.API.KubeConfig)
	if err != nil {
		return err
	}

	// Create the API
	api, err := client.NewAPI(kubeconfig)
	if err != nil {
		return err
	}

	if err = factory.SetValue(&typeDependencies.API, api); err != nil {
		return err
	}
	dependencies.Set("API", api)

	return nil
}
