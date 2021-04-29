package implementations

import (
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/factory"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/pods"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"k8s.io/client-go/kubernetes"
	"reflect"
	"testing"
)

func TestKubernetesPodsFactorySuite(t *testing.T) {
	suite.Run(t, new(testKubernetesPodsFactorySuite))
}

type testKubernetesPodsFactorySuite struct {
	suite.Suite
}

func (s *testKubernetesPodsFactorySuite) TestNewKubernetes() {
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

	var out pods.Pods
	s.Nil(Factory.New(config, dependencies, &out))
	s.NotNil(out)

	// Validate the type of the returned object
	s.Equal("*kubernetes.kubernetesPods", reflect.TypeOf(out).String())
}
