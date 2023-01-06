package factory

import (
	"github.com/gazebo-web/cloudsim/v4/pkg/factory"
	"github.com/gazebo-web/cloudsim/v4/pkg/orchestrator/components/ingresses"
	"github.com/gazebo-web/gz-go/v7"
	"github.com/stretchr/testify/suite"
	"k8s.io/client-go/kubernetes"
	"testing"
)

func TestKubernetesIngressesFactorySuite(t *testing.T) {
	suite.Run(t, new(testKubernetesIngressesFactorySuite))
}

type testKubernetesIngressesFactorySuite struct {
	suite.Suite
}

func (s *testKubernetesIngressesFactorySuite) TestNewFunc() {
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
	s.Nil(IngressesNewFunc(nil, dependencies, &out))
	s.NotNil(out)
}

func TestKubernetesIngressRulesFactorySuite(t *testing.T) {
	suite.Run(t, new(testKubernetesIngressRulesFactorySuite))
}

type testKubernetesIngressRulesFactorySuite struct {
	suite.Suite
}

func (s *testKubernetesIngressRulesFactorySuite) TestNewFunc() {
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
	s.Nil(IngressRulesNewFunc(nil, dependencies, &out))
	s.NotNil(out)
}
