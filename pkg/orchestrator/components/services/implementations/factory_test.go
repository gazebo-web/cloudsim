package implementations

import (
	"github.com/gazebo-web/cloudsim/pkg/factory"
	"github.com/gazebo-web/cloudsim/pkg/orchestrator/components/services"
	"github.com/gazebo-web/gz-go/v7"
	"github.com/stretchr/testify/suite"
	"k8s.io/client-go/kubernetes"
	"reflect"
	"testing"
)

func TestKubernetesServicesFactorySuite(t *testing.T) {
	suite.Run(t, new(testKubernetesServicesFactorySuite))
}

type testKubernetesServicesFactorySuite struct {
	suite.Suite
}

func (s *testKubernetesServicesFactorySuite) TestNewKubernetes() {
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

	var out services.Services
	s.Nil(Factory.New(config, dependencies, &out))
	s.NotNil(out)

	// Validate that the returned object is the correct implementation
	s.Equal("*kubernetes.kubernetesServices", reflect.TypeOf(out).String())
}
