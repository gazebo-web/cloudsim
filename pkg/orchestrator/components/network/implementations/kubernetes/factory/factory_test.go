package factory

import (
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/factory"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/network"
	"gitlab.com/ignitionrobotics/web/ign-go/v5"
	"k8s.io/client-go/kubernetes"
	"testing"
)

func TestKubernetesNetPolFactorySuite(t *testing.T) {
	suite.Run(t, new(testKubernetesNetPolFactorySuite))
}

type testKubernetesNetPolFactorySuite struct {
	suite.Suite
}

func (s *testKubernetesNetPolFactorySuite) TestNewFunc() {
	// Prepare dependencies
	logger := ign.NewLoggerNoRollbar("test", ign.VerbosityWarning)
	kubernetesAPI := struct {
		kubernetes.Interface
	}{}
	dependencies := factory.Dependencies{
		"Logger": logger,
		"API":    kubernetesAPI,
	}

	var out network.Policies
	s.Nil(NewFunc(nil, dependencies, &out))
	s.NotNil(out)
}
