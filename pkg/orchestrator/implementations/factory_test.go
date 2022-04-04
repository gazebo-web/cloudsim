package implementations

import (
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/factory"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	configurationsImpl "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/configurations/implementations"
	ingressesImpl "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/ingresses/implementations"
	networkImpl "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/network/implementations"
	nodesImpl "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/nodes/implementations"
	podsImpl "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/pods/implementations"
	servicesImpl "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/services/implementations"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"k8s.io/client-go/kubernetes"
	"reflect"
	"testing"
)

func TestKubernetesFactorySuite(t *testing.T) {
	suite.Run(t, new(testKubernetesFactorySuite))
}

type testKubernetesFactorySuite struct {
	suite.Suite
}

func (s *testKubernetesFactorySuite) createFactoryConfig(objectType string,
	config factory.ConfigValues) *factory.Config {

	return &factory.Config{
		Type:   objectType,
		Config: config,
	}
}

func (s *testKubernetesFactorySuite) callFactory(config *factory.Config, out interface{}) {
	// Prepare dependencies
	logger := ign.NewLoggerNoRollbar("test", ign.VerbosityWarning)
	kubernetesAPI := struct {
		kubernetes.Interface
	}{}
	dependencies := factory.Dependencies{
		"Logger": logger,
		"API":    kubernetesAPI,
	}

	s.Nil(Factory.New(config, dependencies, out))
	s.NotNil(out)
}

func (s *testKubernetesFactorySuite) TestIngressesNewKubernetes() {
	// Prepare config
	config := factory.Config{
		Type: Kubernetes,
		Config: factory.ConfigValues{
			"api": factory.ConfigValues{
				"kubeconfig": "",
			},
			"components": factory.ConfigValues{
				"nodes":           factory.ConfigValues{"type": nodesImpl.Kubernetes},
				"pods":            factory.ConfigValues{"type": podsImpl.Kubernetes},
				"services":        factory.ConfigValues{"type": servicesImpl.Kubernetes},
				"ingresses":       factory.ConfigValues{"type": ingressesImpl.Kubernetes},
				"ingressRules":    factory.ConfigValues{"type": ingressesImpl.Kubernetes},
				"networkPolicies": factory.ConfigValues{"type": networkImpl.Kubernetes},
				"configurations":  factory.ConfigValues{"type": configurationsImpl.Kubernetes},
			},
		},
	}

	var cluster orchestrator.Cluster
	s.callFactory(&config, &cluster)

	// Validate the type of the returned object
	s.Equal("*kubernetes.k8s", reflect.TypeOf(cluster).String())
}

func (s *testKubernetesFactorySuite) TestIngressesNewKubernetesGloo() {
	config := factory.Config{
		Type: Kubernetes,
		Config: factory.ConfigValues{
			"api": factory.ConfigValues{
				"kubeconfig": "",
			},
			"components": factory.ConfigValues{
				"nodes":           factory.ConfigValues{"type": nodesImpl.Kubernetes},
				"pods":            factory.ConfigValues{"type": podsImpl.Kubernetes},
				"services":        factory.ConfigValues{"type": servicesImpl.Kubernetes},
				"ingresses":       factory.ConfigValues{"type": ingressesImpl.Gloo},
				"ingressRules":    factory.ConfigValues{"type": ingressesImpl.Gloo},
				"networkPolicies": factory.ConfigValues{"type": networkImpl.Kubernetes},
				"configurations":  factory.ConfigValues{"type": configurationsImpl.Kubernetes},
			},
		},
	}

	var cluster orchestrator.Cluster
	s.callFactory(&config, &cluster)

	// Validate the type of the returned object
	s.Equal("*kubernetes.k8s", reflect.TypeOf(cluster).String())
	s.Equal("*gloo.VirtualServices", reflect.TypeOf(cluster.Ingresses()).String())
	s.Equal("*gloo.virtualHosts", reflect.TypeOf(cluster.IngressRules()).String())
}
