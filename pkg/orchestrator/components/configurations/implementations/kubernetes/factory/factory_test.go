package factory

import (
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/factory"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/configurations"
	"gitlab.com/ignitionrobotics/web/ign-go/v6"
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
	logger := ign.NewLoggerNoRollbar("test", ign.VerbosityWarning)
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
