package factory

import (
	"github.com/gazebo-web/cloudsim/v4/pkg/factory"
	configMapsImpl "github.com/gazebo-web/cloudsim/v4/pkg/orchestrator/components/configurations/implementations"
	ingressesImpl "github.com/gazebo-web/cloudsim/v4/pkg/orchestrator/components/ingresses/implementations"
	networkImpl "github.com/gazebo-web/cloudsim/v4/pkg/orchestrator/components/network/implementations"
	nodesImpl "github.com/gazebo-web/cloudsim/v4/pkg/orchestrator/components/nodes/implementations"
	podsImpl "github.com/gazebo-web/cloudsim/v4/pkg/orchestrator/components/pods/implementations"
	servicesImpl "github.com/gazebo-web/cloudsim/v4/pkg/orchestrator/components/services/implementations"
	kubernetesCluster "github.com/gazebo-web/cloudsim/v4/pkg/orchestrator/implementations/kubernetes"
	"github.com/gazebo-web/cloudsim/v4/pkg/orchestrator/implementations/kubernetes/client"
	"reflect"
)

// NewFunc is the factory creation function for the Kubernetes orchestrator.Cluster implementation.
func NewFunc(config interface{}, dependencies factory.Dependencies, out interface{}) error {
	// Parse config
	var typeConfig Config
	if err := factory.SetValueAndValidate(&typeConfig, config); err != nil {
		return factory.ErrorWithContext(err)
	}

	// Initialize orchestrator component configs
	initializeComponentConfig(&typeConfig)

	// Parse dependencies
	var err error
	if dependencies, err = dependencies.DeepCopy(); err != nil {
		return factory.ErrorWithContext(err)
	}
	var typeDependencies Dependencies
	if err := dependencies.ToStruct(&typeDependencies); err != nil {
		return factory.ErrorWithContext(err)
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
		// ConfigMaps
		{
			Factory:      configMapsImpl.Factory,
			Config:       typeConfig.Components.Configurations,
			Dependencies: dependencies,
			Out:          &components.Configurations,
		},
	}
	if err := factory.CallFactories(factoryCalls); err != nil {
		return err
	}

	// Set output value
	cluster := kubernetesCluster.NewCustomKubernetes(components)
	if err := factory.SetValue(out, cluster); err != nil {
		return factory.ErrorWithContext(err)
	}

	return nil
}

// initializeComponentConfig initializes a component config valies.
// Orchestrator's API config values will be used to set the component config values if it does not have any component
// values defined.
func initializeComponentConfig(typeConfig *Config) {
	v := reflect.ValueOf(typeConfig.Components)

	// Replace every `nil` component config with the orchestrator API config
	for i := 0; i < v.NumField(); i++ {
		config := v.Field(i).Interface().(*factory.Config)
		if config.Config == nil {
			config.Config = factory.ConfigValues{
				"api": typeConfig.API,
			}
		}
	}
}

// initializeAPI initializes the API dependency.
func initializeAPI(config *Config, dependencies factory.Dependencies, typeDependencies *Dependencies) error {
	if typeDependencies.API != nil {
		return nil
	}
	if config == nil {
		return factory.ErrorWithContext(factory.ErrNilConfig)
	}

	// Get the Kubernetes config
	kubeconfig, err := client.GetConfig(config.API.KubeConfig)
	if err != nil {
		return factory.ErrorWithContext(err)
	}

	// Create the API
	api, err := client.NewAPI(kubeconfig)
	if err != nil {
		return factory.ErrorWithContext(err)
	}

	if err = factory.SetValue(&typeDependencies.API, api); err != nil {
		return factory.ErrorWithContext(err)
	}
	dependencies.Set("API", api)

	return nil
}
