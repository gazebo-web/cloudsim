package factory

import (
	"github.com/gazebo-web/cloudsim/pkg/factory"
	"github.com/gazebo-web/cloudsim/pkg/orchestrator/components/services"
	"github.com/gazebo-web/gz-go/v7"
	"github.com/stretchr/testify/suite"
	"k8s.io/client-go/kubernetes"
	"testing"
)

func TestKubernetesServicesFactorySuite(t *testing.T) {
	suite.Run(t, new(testKubernetesServicesFactorySuite))
}

type testKubernetesServicesFactorySuite struct {
	suite.Suite
}

func (s *testKubernetesServicesFactorySuite) TestNewFunc() {
	// Prepare dependencies
	logger := gz.NewLoggerNoRollbar("test", gz.VerbosityWarning)
	kubernetesAPI := struct {
		kubernetes.Interface
	}{}
	dependencies := factory.Dependencies{
		"Logger": logger,
		"API":    kubernetesAPI,
	}

	var out services.Services
	s.Nil(NewFunc(nil, dependencies, &out))
	s.NotNil(out)
}
