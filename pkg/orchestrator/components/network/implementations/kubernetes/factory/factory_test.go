package factory

import (
	"github.com/gazebo-web/cloudsim/pkg/factory"
	"github.com/gazebo-web/cloudsim/pkg/orchestrator/components/network"
	"github.com/gazebo-web/gz-go/v7"
	"github.com/stretchr/testify/suite"
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
	logger := gz.NewLoggerNoRollbar("test", gz.VerbosityWarning)
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
