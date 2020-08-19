package starter

import (
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud/aws/ec2"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud/aws/s3"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/ingresses"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/ingresses/rules"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/nodes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/pods"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/spdy"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	simfake "gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations/fake"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulator"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestSimulatorStarter(t *testing.T) {
	suite.Run(t, new(starterTestSuite))
}

type starterTestSuite struct {
	suite.Suite
	starter           simulator.Starter
	k8s               *fake.Clientset
	nodes             orchestrator.Nodes
	spdy              *spdy.Fake
	pods              orchestrator.Pods
	ingresses         orchestrator.Ingresses
	orchestrator      orchestrator.Cluster
	ingressRules      orchestrator.IngressRules
	logger            ign.Logger
	awsAPI            *awsAPI
	storage           cloud.Storage
	machines          cloud.Machines
	simulationService *simfake.Service
}

func (s *starterTestSuite) initializeKubernetes() {
	s.k8s = fake.NewSimpleClientset()
	s.nodes = nodes.NewNodes(s.k8s)
	s.spdy = spdy.NewSPDYFakeInitializer()
	s.pods = pods.NewPods(s.k8s, s.spdy)
	s.ingresses = ingresses.NewIngresses(s.k8s)
	s.ingressRules = rules.NewIngressRules(s.k8s)
	s.orchestrator = kubernetes.NewKubernetes(s.nodes, s.pods, s.ingresses, s.ingressRules)
}

func (s *starterTestSuite) initializeAWS() {
	s.awsAPI = &awsAPI{}
	s.storage = s3.NewStorage(s.awsAPI, s.logger)
	s.machines = ec2.NewMachines(s.awsAPI, s.logger)
}

func (s *starterTestSuite) SetupTest() {
	s.logger = ign.NewLoggerNoRollbar("SimulatorStarterTestSuite", ign.VerbosityDebug)
	s.initializeKubernetes()
	s.initializeAWS()
	s.simulationService = simfake.NewService()
	s.starter = NewSimulatorStarter(s.orchestrator, s.machines, s.storage, s.simulationService)
}

func (s *starterTestSuite) TestStartSimulation_IncorrectStatusWhenSimulationsIsRunning() {
	s.simulationService.
		On("Get", simulations.GroupID("test-group-id")).
		Return(
			simfake.NewSimulation("test-group-id", simulations.StatusRunning, simulations.SimSingle),
			nil,
		)

	ok, err := s.starter.HasStatus("test-group-id", simulations.StatusPending)
	s.False(ok)
	s.NoError(err)
}

func (s *starterTestSuite) TestStartSimulation_CorrectStatus() {
	s.simulationService.
		On("Get", simulations.GroupID("test-group-id")).
		Return(
			simfake.NewSimulation("test-group-id", simulations.StatusPending, simulations.SimSingle),
			nil,
		)

	ok, err := s.starter.HasStatus("test-group-id", simulations.StatusPending)
	s.True(ok)
	s.NoError(err)
}

func (s *starterTestSuite) TestStartSimulation_IsNotParent() {
	s.simulationService.On("Get", simulations.GroupID("test-group-id")).Return(
		simfake.NewSimulation("test-group-id", simulations.StatusPending, simulations.SimSingle),
		nil,
	)

	ok, err := s.starter.IsParent("test-group-id")
	s.False(ok)
	s.NoError(err)
}

func (s *starterTestSuite) TestStartSimulation_IsParent() {
	s.simulationService.
		On("Get", simulations.GroupID("test-group-id")).
		Return(
			simfake.NewSimulation("test-group-id", simulations.StatusPending, simulations.SimParent),
			nil,
		)

	ok, err := s.starter.IsParent("test-group-id")
	s.True(ok)
	s.NoError(err)
}

type awsAPI struct {
	ec2iface.EC2API
	s3iface.S3API
}
