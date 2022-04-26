package implementations

import (
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/factory"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/secrets"
	"gitlab.com/ignitionrobotics/web/ign-go/v5"
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

func (s *testSecretsFactorySuite) createFactoryConfig(objectType string,
	config factory.ConfigValues) *factory.Config {

	return &factory.Config{
		Type:   objectType,
		Config: config,
	}
}

func (s *testSecretsFactorySuite) TestIngressesNewKubernetes() {
	// Prepare config
	config := &factory.Config{
		Type: Kubernetes,
	}

	// Prepare dependencies
	logger := ign.NewLoggerNoRollbar("test", ign.VerbosityWarning)
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
