package factory

import (
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/factory"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/services"
	"gitlab.com/ignitionrobotics/web/ign-go/v5"
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
	logger := ign.NewLoggerNoRollbar("test", ign.VerbosityWarning)
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
