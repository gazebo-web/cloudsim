package factory

import (
	"github.com/gazebo-web/cloudsim/v4/pkg/factory"
	"github.com/gazebo-web/cloudsim/v4/pkg/orchestrator/components/configurations"
	"github.com/gazebo-web/gz-go/v7"
	"github.com/stretchr/testify/suite"
	"k8s.io/client-go/kubernetes"
	"testing"
)

func TestKubernetesConfigurationsFactorySuite(t *testing.T) {
	suite.Run(t, new(testKubernetesConfigurationsFactorySuite))
}

type testKubernetesConfigurationsFactorySuite struct {
	suite.Suite
}

func (s *testKubernetesConfigurationsFactorySuite) TestNewFunc() {
	// Prepare dependencies
	logger := gz.NewLoggerNoRollbar("test", gz.VerbosityWarning)
	kubernetesAPI := struct {
		kubernetes.Interface
	}{}
	dependencies := factory.Dependencies{
		"Logger": logger,
		"API":    kubernetesAPI,
	}

	var out configurations.Configurations
	s.Nil(NewFunc(nil, dependencies, &out))
	s.NotNil(out)
}
