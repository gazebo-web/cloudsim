package implementations

import (
	"github.com/gazebo-web/cloudsim/v4/pkg/factory"
	"github.com/gazebo-web/cloudsim/v4/pkg/orchestrator"
	configurationsImpl "github.com/gazebo-web/cloudsim/v4/pkg/orchestrator/components/configurations/implementations"
	ingressesImpl "github.com/gazebo-web/cloudsim/v4/pkg/orchestrator/components/ingresses/implementations"
	networkImpl "github.com/gazebo-web/cloudsim/v4/pkg/orchestrator/components/network/implementations"
	nodesImpl "github.com/gazebo-web/cloudsim/v4/pkg/orchestrator/components/nodes/implementations"
	podsImpl "github.com/gazebo-web/cloudsim/v4/pkg/orchestrator/components/pods/implementations"
	servicesImpl "github.com/gazebo-web/cloudsim/v4/pkg/orchestrator/components/services/implementations"
	"github.com/gazebo-web/gz-go/v7"
	"github.com/stretchr/testify/suite"
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

func (s *testKubernetesFactorySuite) callFactory(config *factory.Config, out interface{}) {
	// Prepare dependencies
	logger := gz.NewLoggerNoRollbar("test", gz.VerbosityWarning)
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
