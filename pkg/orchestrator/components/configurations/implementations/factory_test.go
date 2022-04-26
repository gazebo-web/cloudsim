package implementations

import (
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/factory"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/configurations"
	"gitlab.com/ignitionrobotics/web/ign-go/v5"
	"k8s.io/client-go/kubernetes"
	"reflect"
	"testing"
)

func TestKubernetesConfigurationsFactorySuite(t *testing.T) {
	suite.Run(t, new(testKubernetesConfigurationsFactorySuite))
}

type testKubernetesConfigurationsFactorySuite struct {
	suite.Suite
}

func (s *testKubernetesConfigurationsFactorySuite) TestNewKubernetes() {
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

	var out configurations.Configurations
	s.Nil(Factory.New(config, dependencies, &out))
	s.NotNil(out)

	// Validate that the returned object is the correct implementation
	s.Equal("*kubernetes.configMaps", reflect.TypeOf(out).String())
}
