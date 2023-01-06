package factory

import (
	"github.com/gazebo-web/cloudsim/v4/pkg/factory"
	"github.com/stretchr/testify/suite"
	"k8s.io/client-go/kubernetes"
	"testing"
)

func TestKubernetesFactorySuite(t *testing.T) {
	suite.Run(t, new(testKubernetesFactorySuite))
}

type testKubernetesFactorySuite struct {
	suite.Suite
}

func (s *testKubernetesFactorySuite) TestInitializeAPIDependencyIsNil() {
	config := Config{
		API: APIConfig{
			KubeConfig: "",
		},
	}
	dependencies := Dependencies{}
	s.Nil(initializeAPI(&config, factory.Dependencies{}, &dependencies))
	s.NotNil(dependencies.API)
}

func (s *testKubernetesFactorySuite) TestInitializeAPIDependencyIsNotNil() {
	// Prepare dependencies
	kubernetesAPI := struct {
		kubernetes.Interface
	}{}
	dependencies := Dependencies{
		API: kubernetesAPI,
	}

	s.Nil(initializeAPI(nil, factory.Dependencies{}, &dependencies))
	s.Exactly(kubernetesAPI, dependencies.API)
}

func (s *testKubernetesFactorySuite) TestInitializeAPIDependencyAndConfigAreNil() {
	dependencies := Dependencies{}

	s.NotNil(initializeAPI(nil, factory.Dependencies{}, &dependencies))
}

func (s *testKubernetesFactorySuite) TestInitializeComponentConfig() {
	orchestratorConfig := APIConfig{
		KubeConfig: "orchestrator",
	}
	componentConfig := APIConfig{
		KubeConfig: "component",
	}
	typeConfig := &Config{
		API: orchestratorConfig,
		Components: Components{
			Nodes:           &factory.Config{},
			Pods:            &factory.Config{},
			Ingresses:       &factory.Config{Config: factory.ConfigValues{"api": componentConfig}},
			IngressRules:    &factory.Config{Config: factory.ConfigValues{"api": componentConfig}},
			Services:        &factory.Config{},
			NetworkPolicies: &factory.Config{},
			Configurations:  &factory.Config{},
		},
	}

	initializeComponentConfig(typeConfig)
	// Verify that the nil configs have been updated
	s.Require().Equal(orchestratorConfig, typeConfig.Components.Nodes.Config["api"])
	s.Require().Equal(orchestratorConfig, typeConfig.Components.Pods.Config["api"])
	s.Require().Equal(orchestratorConfig, typeConfig.Components.Services.Config["api"])
	s.Require().Equal(orchestratorConfig, typeConfig.Components.NetworkPolicies.Config["api"])
	s.Require().Equal(orchestratorConfig, typeConfig.Components.Configurations.Config["api"])
	// Verify that the pre-existing configs were not updated
	s.Require().Equal(componentConfig, typeConfig.Components.Ingresses.Config["api"])
	s.Require().Equal(componentConfig, typeConfig.Components.IngressRules.Config["api"])
}
