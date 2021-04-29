package factory

import (
	"github.com/stretchr/testify/suite"
	"k8s.io/client-go/kubernetes"
	"testing"
)

func TestKubernetesSecretsFactorySuite(t *testing.T) {
	suite.Run(t, new(testKubernetesSecretsFactorySuite))
}

type testKubernetesSecretsFactorySuite struct {
	suite.Suite
}

func (s *testKubernetesSecretsFactorySuite) TestInitializeAPIDependencyIsNil() {
	// Prepare config
	config := Config{
		API: APIConfig{
			KubeConfig: "",
		},
	}

	// Prepare dependencies
	dependencies := Dependencies{}

	s.Nil(initializeAPI(&config, &dependencies))
	s.NotNil(dependencies.API)
}

func (s *testKubernetesSecretsFactorySuite) TestInitializeAPIDependencyIsNotNil() {
	// Prepare dependencies
	kubernetesAPI := struct {
		kubernetes.Interface
	}{}
	dependencies := Dependencies{
		API: kubernetesAPI,
	}

	s.Nil(initializeAPI(nil, &dependencies))
	s.Exactly(kubernetesAPI, dependencies.API)
}
