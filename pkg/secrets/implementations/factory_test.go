package implementations

import (
	"github.com/gazebo-web/cloudsim/pkg/factory"
	"github.com/gazebo-web/cloudsim/pkg/secrets"
	"github.com/gazebo-web/gz-go/v7"
	"github.com/stretchr/testify/suite"
	"k8s.io/client-go/kubernetes"
	kfake "k8s.io/client-go/kubernetes/fake"
	"reflect"
	"testing"
)

func TestSecretsFactorySuite(t *testing.T) {
	suite.Run(t, new(testSecretsFactorySuite))
}

type testSecretsFactorySuite struct {
	suite.Suite
}

func (s *testSecretsFactorySuite) TestIngressesNewKubernetes() {
	// Prepare config
	config := &factory.Config{
		Type: Kubernetes,
	}

	// Prepare dependencies
	logger := gz.NewLoggerNoRollbar("test", gz.VerbosityWarning)
	kubernetesAPI := struct {
		kubernetes.Interface
	}{
		kfake.NewSimpleClientset(),
	}
	dependencies := factory.Dependencies{
		"Logger": logger,
		"API":    kubernetesAPI,
	}

	var secrets secrets.Secrets
	s.Nil(Factory.New(config, dependencies, &secrets))
	s.NotNil(&secrets)

	// Validate the type of the returned object
	s.Equal("*kubernetes.kubernetesSecrets", reflect.TypeOf(secrets).String())
}
