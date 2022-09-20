package factory

import (
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/factory"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/ingresses"
	"gitlab.com/ignitionrobotics/web/ign-go/v6"
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
	logger := ign.NewLoggerNoRollbar("test", ign.VerbosityWarning)
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
	logger := ign.NewLoggerNoRollbar("test", ign.VerbosityWarning)
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
