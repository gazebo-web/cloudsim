package factory

import (
	gatewayFake "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/client/clientset/versioned/fake"
	glooFake "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/kube/client/clientset/versioned/fake"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/factory"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/ingresses"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"k8s.io/client-go/kubernetes"
	"testing"
)

func TestGlooIngressesFactorySuite(t *testing.T) {
	suite.Run(t, new(testGlooIngressesFactorySuite))
}

type testGlooIngressesFactorySuite struct {
	suite.Suite
	config       Config
	dependencies factory.Dependencies
}

func (s *testGlooIngressesFactorySuite) SetupTest() {
	// Prepare config
	s.config = Config{
		API: APIConfig{
			KubeConfig: "",
		},
	}

	// Prepare dependencies
	logger := ign.NewLoggerNoRollbar("test", ign.VerbosityWarning)
	kubernetesAPI := struct {
		kubernetes.Interface
	}{}
	s.dependencies = factory.Dependencies{
		"Logger": logger,
		"API":    kubernetesAPI,
	}
}

func (s *testGlooIngressesFactorySuite) TestIngressesNewFunc() {
	var out ingresses.Ingresses
	s.Nil(IngressesNewFunc(s.config, s.dependencies, &out))
	s.NotNil(out)
}

func (s *testGlooIngressesFactorySuite) TestIngressRulesNewFunc() {
	var out ingresses.IngressRules
	s.Nil(IngressRulesNewFunc(s.config, s.dependencies, &out))
	s.NotNil(out)
}

func (s *testGlooIngressesFactorySuite) TestInitializeGlooDependencyIsNil() {
	dependencies := Dependencies{}
	s.Nil(initializeGloo(&s.config, &dependencies))
	s.NotNil(dependencies.Gloo)
}

func (s *testGlooIngressesFactorySuite) TestInitializeGlooDependencyIsNotNil() {
	fake := glooFake.NewSimpleClientset().GlooV1()
	dependencies := Dependencies{
		Gloo: fake,
	}

	s.Nil(initializeGloo(&s.config, &dependencies))
	s.Exactly(fake, dependencies.Gloo)
}

func (s *testGlooIngressesFactorySuite) TestInitializeGlooGatewayDependencyIsNil() {
	dependencies := Dependencies{}
	s.Nil(initializeGlooGateway(&s.config, &dependencies))
	s.NotNil(dependencies.GlooGateway)
}

func (s *testGlooIngressesFactorySuite) TestInitializeGlooGatewayDependencyIsNotNil() {
	fake := gatewayFake.NewSimpleClientset().GatewayV1()
	dependencies := Dependencies{
		GlooGateway: fake,
	}

	s.Nil(initializeGloo(&s.config, &dependencies))
	s.Exactly(fake, dependencies.GlooGateway)
}
