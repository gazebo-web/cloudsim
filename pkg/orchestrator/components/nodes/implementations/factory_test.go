package implementations

import (
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/factory"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/nodes"
	"gitlab.com/ignitionrobotics/web/ign-go/v6"
	"k8s.io/client-go/kubernetes"
	"reflect"
	"testing"
)

func TestKubernetesNodesFactorySuite(t *testing.T) {
	suite.Run(t, new(testKubernetesNodesFactorySuite))
}

type testKubernetesNodesFactorySuite struct {
	suite.Suite
}

func (s *testKubernetesNodesFactorySuite) TestNewKubernetes() {
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

	var out nodes.Nodes
	s.Nil(Factory.New(config, dependencies, &out))
	s.NotNil(out)

	// Validate that the returned object is the correct implementation
	s.Equal("*kubernetes.kubernetesNodes", reflect.TypeOf(out).String())
}
