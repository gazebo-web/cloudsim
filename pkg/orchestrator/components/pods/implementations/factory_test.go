package implementations

import (
	"github.com/gazebo-web/cloudsim/v4/pkg/factory"
	"github.com/gazebo-web/cloudsim/v4/pkg/orchestrator/components/pods"
	"github.com/gazebo-web/gz-go/v7"
	"github.com/stretchr/testify/suite"
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
	logger := gz.NewLoggerNoRollbar("test", gz.VerbosityWarning)
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
