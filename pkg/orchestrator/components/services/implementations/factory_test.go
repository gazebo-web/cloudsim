package implementations

import (
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/factory"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/services"
	"gitlab.com/ignitionrobotics/web/ign-go/v5"
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
	logger := ign.NewLoggerNoRollbar("test", ign.VerbosityWarning)
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
