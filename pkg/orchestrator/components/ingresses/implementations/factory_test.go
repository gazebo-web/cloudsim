package implementations

import (
	"github.com/gazebo-web/cloudsim/v4/pkg/factory"
	"github.com/gazebo-web/cloudsim/v4/pkg/orchestrator/components/ingresses"
	"github.com/gazebo-web/gz-go/v7"
	"github.com/stretchr/testify/suite"
	"k8s.io/client-go/kubernetes"
	"reflect"
	"testing"
)

func TestKubernetesIngressesFactorySuite(t *testing.T) {
	suite.Run(t, new(testKubernetesIngressesFactorySuite))
}

type testKubernetesIngressesFactorySuite struct {
	suite.Suite
}

func (s *testKubernetesIngressesFactorySuite) TestIngressesNewKubernetes() {
	// Prepare config
	config := &factory.Config{
		Type: Kubernetes,
	}

	// Prepare dependencies
	logger := gz.NewLoggerNoRollbar("test", gz.VerbosityWarning)
	kubernetesAPI := struct {
		kubernetes.Interface
	}{}
	dependencies := factory.Dependencies{
		"Logger": logger,
		"API":    kubernetesAPI,
	}

	var out ingresses.Ingresses
	s.Nil(IngressesFactory.New(config, dependencies, &out))
	s.NotNil(out)

	// Validate the type of the returned object
	s.Equal("*kubernetes.kubernetesIngresses", reflect.TypeOf(out).String())
}

func (s *testKubernetesIngressesFactorySuite) TestIngressRulesNewKubernetes() {
	// Prepare config
	config := &factory.Config{
		Type: Kubernetes,
	}

	// Prepare dependencies
	logger := gz.NewLoggerNoRollbar("test", gz.VerbosityWarning)
	kubernetesAPI := struct {
		kubernetes.Interface
	}{}
	dependencies := factory.Dependencies{
		"Logger": logger,
		"API":    kubernetesAPI,
	}

	var out ingresses.IngressRules
	s.Nil(IngressRulesFactory.New(config, dependencies, &out))
	s.NotNil(out)

	// Validate the type of the returned object
	s.Equal("*rules.ingressRules", reflect.TypeOf(out).String())
}
