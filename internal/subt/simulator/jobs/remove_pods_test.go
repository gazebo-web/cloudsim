package jobs

import (
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/suite"
	subtapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/pods"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	simfake "gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations/fake"
	sfake "gitlab.com/ignitionrobotics/web/cloudsim/pkg/store/fake"
	gormdb "gitlab.com/ignitionrobotics/web/cloudsim/pkg/utils/db/gorm"
	"gitlab.com/ignitionrobotics/web/ign-go"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kfake "k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestRemovePods(t *testing.T) {
	suite.Run(t, new(removePodsTestSuite))
}

type removePodsTestSuite struct {
	suite.Suite
	Logger              ign.Logger
	DB                  *gorm.DB
	Client              *kfake.Clientset
	Orchestrator        orchestrator.Cluster
	Platform            platform.Platform
	Namespace           string
	Store               *sfake.Fake
	GroupID             simulations.GroupID
	StopSimulationState *state.StopSimulation
	ActionStore         actions.Store
	SimulationService   *simfake.Service
	ApplicationServices subtapp.Services
	Robots              []simulations.Robot
	Pods                []corev1.Pod
}

func (s *removePodsTestSuite) SetupTest() {
	s.Logger = ign.NewLoggerNoRollbar("TestRemovePods", ign.VerbosityDebug)

	s.Namespace = "default"

	storeIgnition := sfake.NewFakeIgnition()
	storeIgnition.On("LogsCopyEnabled").Return(true)

	storeOrchestrator := sfake.NewFakeOrchestrator()
	storeOrchestrator.On("Namespace").Return(s.Namespace)

	s.Store = sfake.NewFakeStore(nil, storeOrchestrator, storeIgnition)

	var err error
	s.DB, err = gormdb.GetDBFromEnvVars()

	s.Require().NoError(err)

	err = actions.MigrateDB(s.DB)
	s.Require().NoError(err)

	s.GroupID = "aaaa-bbbb-cccc-dddd"

	robot := simfake.NewRobot("test", "x1")

	s.Client = kfake.NewSimpleClientset(
		// Gazebo server
		&corev1.Pod{
			TypeMeta: metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{
				Name:      subtapp.GetPodNameGazeboServer(s.GroupID),
				Namespace: s.Namespace,
				Labels:    subtapp.GetPodLabelsGazeboServer(s.GroupID, nil).Map(),
			},
		},
		// Gazebo copy pod
		&corev1.Pod{
			TypeMeta: metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{
				Name:      subtapp.GetPodNameGazeboServerCopy(s.GroupID),
				Namespace: s.Namespace,
				Labels:    subtapp.GetPodLabelsGazeboServerCopy(s.GroupID, nil).Map(),
			},
		},
		// Field computer robot 1
		&corev1.Pod{
			TypeMeta: metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{
				Name:      subtapp.GetPodNameFieldComputer(s.GroupID, subtapp.GetRobotID(0)),
				Namespace: s.Namespace,
				Labels:    subtapp.GetPodLabelsFieldComputer(s.GroupID, nil).Map(),
			},
		},
		// Comms bridge robot 1
		&corev1.Pod{
			TypeMeta: metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{
				Name:      subtapp.GetPodNameCommsBridge(s.GroupID, subtapp.GetRobotID(0)),
				Namespace: s.Namespace,
				Labels:    subtapp.GetPodLabelsCommsBridge(s.GroupID, nil, robot).Map(),
			},
		},
		// Comms bridge copy pod robot 1
		&corev1.Pod{
			TypeMeta: metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{
				Name:      subtapp.GetPodNameCommsBridgeCopy(s.GroupID, subtapp.GetRobotID(0)),
				Namespace: s.Namespace,
				Labels:    subtapp.GetPodLabelsCommsBridgeCopy(s.GroupID, nil, robot).Map(),
			},
		},
	)

	po := pods.NewPods(s.Client, nil, s.Logger)

	s.Orchestrator = kubernetes.NewCustomKubernetes(kubernetes.Config{
		Pods: po,
	})

	s.Platform = platform.NewPlatform(platform.Components{
		Cluster: s.Orchestrator,
		Store:   s.Store,
	})

	s.SimulationService = simfake.NewService()

	s.Robots = []simulations.Robot{
		robot,
	}

	s.SimulationService.On("GetRobots", s.GroupID).Return(s.Robots, error(nil))

	services := application.NewServices(s.SimulationService, nil)

	s.ApplicationServices = subtapp.NewServices(services, nil, nil)

	s.StopSimulationState = state.NewStopSimulation(s.Platform, s.ApplicationServices, s.GroupID)

	s.ActionStore = actions.NewStore(s.StopSimulationState)
}

func (s *removePodsTestSuite) TestRemovePodsSuccess() {
	_, err := RemovePods.Run(s.ActionStore, s.DB, &actions.Deployment{}, s.StopSimulationState)
	s.Assert().NoError(err)
}
