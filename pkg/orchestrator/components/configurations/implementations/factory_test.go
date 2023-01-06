package implementations

import (
	"github.com/gazebo-web/cloudsim/pkg/factory"
	"github.com/gazebo-web/cloudsim/pkg/orchestrator/components/configurations"
	"github.com/gazebo-web/gz-go/v7"
	"github.com/stretchr/testify/suite"
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
	logger := gz.NewLoggerNoRollbar("test", gz.VerbosityWarning)
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
