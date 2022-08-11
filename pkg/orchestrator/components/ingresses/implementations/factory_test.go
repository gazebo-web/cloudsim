package implementations

import (
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/factory"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/ingresses"
	"gitlab.com/ignitionrobotics/web/ign-go/v5"
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
	logger := ign.NewLoggerNoRollbar("test", ign.VerbosityWarning)
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
	logger := ign.NewLoggerNoRollbar("test", ign.VerbosityWarning)
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
