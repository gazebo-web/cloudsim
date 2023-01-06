package factory

import (
	"github.com/gazebo-web/cloudsim/pkg/orchestrator/components/spdy"
	"github.com/stretchr/testify/suite"
	"testing"
)

func TestKubernetesPodsFactorySuite(t *testing.T) {
	suite.Run(t, new(testKubernetesPodsFactorySuite))
}

type testKubernetesPodsFactorySuite struct {
	suite.Suite
}

func (s *testKubernetesPodsFactorySuite) TestInitializeSPDYDependencyIsNil() {
	config := Config{
		API: APIConfig{
			KubeConfig: "",
		},
	}
	dependencies := Dependencies{}
	s.Nil(initializeSPDY(&config, &dependencies))
	s.NotNil(dependencies.SPDY)
}

func (s *testKubernetesPodsFactorySuite) TestInitializeSPDYDependencyIsNotNil() {
	// Prepare dependencies
	spdy := spdy.NewSPDYFakeInitializer()
	dependencies := Dependencies{
		SPDY: spdy,
	}

	s.Nil(initializeSPDY(nil, &dependencies))
	s.Exactly(spdy, dependencies.SPDY)
}
