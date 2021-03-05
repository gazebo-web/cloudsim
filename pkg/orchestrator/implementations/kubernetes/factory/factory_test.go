package factory

import (
	"github.com/stretchr/testify/suite"
	"k8s.io/client-go/kubernetes"
	"testing"
)

func TestKubernetesFactorySuite(t *testing.T) {
	suite.Run(t, new(testKubernetesFactorySuite))
}

type testKubernetesFactorySuite struct {
	suite.Suite
}

func (s *testKubernetesFactorySuite) TestInitializeAPIDependencyIsNil() {
	config := Config{
		API: APIConfig{
			KubeConfig: "",
		},
	}
	dependencies := Dependencies{}
	s.Nil(initializeAPI(&config, &dependencies))
	s.NotNil(dependencies.API)
}

func (s *testKubernetesFactorySuite) TestInitializeAPIDependencyIsNotNil() {
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

func (s *testKubernetesFactorySuite) TestInitializeAPIDependencyAndConfigAreNil() {
	dependencies := Dependencies{}

	s.NotNil(initializeAPI(nil, &dependencies))
}
