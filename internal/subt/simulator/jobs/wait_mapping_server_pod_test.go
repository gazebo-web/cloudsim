package jobs

import (
	"github.com/stretchr/testify/suite"
	subtapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	subtsimfake "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulations/fake"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/tracks"
	tfake "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/tracks/fake"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	kubernetesPods "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/pods/implementations/kubernetes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/spdy"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/implementations/kubernetes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	simfake "gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations/fake"
	sfake "gitlab.com/ignitionrobotics/web/cloudsim/pkg/store/implementations/fake"
	"gitlab.com/ignitionrobotics/web/ign-go"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"knative.dev/pkg/ptr"
	"testing"
	"time"
)

func TestWaitMappingServerPod(t *testing.T) {
	suite.Run(t, new(waitMappingServerPodTestSuite))
}

type waitMappingServerPodTestSuite struct {
	suite.Suite
}

func (suite *waitMappingServerPodTestSuite) waitForPod(sim simulations.Simulation, client *fake.Clientset,
	mappingImage bool) (interface{}, error) {

	logger := ign.NewLoggerNoRollbar("TestWaitForMappingServerPod", ign.VerbosityDebug)
	orchestratorStore := sfake.NewFakeOrchestrator()
	fakeStore := sfake.NewFakeStore(nil, orchestratorStore, nil, nil)
	spdyInit := spdy.NewSPDYFakeInitializer()

	po := kubernetesPods.NewPods(client, spdyInit, logger)
	ks := kubernetes.NewCustomKubernetes(kubernetes.Config{
		Nodes:           nil,
		Pods:            po,
		Ingresses:       nil,
		IngressRules:    nil,
		Services:        nil,
		NetworkPolicies: nil,
	})

	p, _ := platform.NewPlatform("test", platform.Components{
		Cluster: ks,
		Store:   fakeStore,
	})

	// Initialize tracks service
	track := &tracks.Track{
		Name:          "test",
		Image:         "world-image.org/image",
		BridgeImage:   "bridge-image.org/image",
		StatsTopic:    "test",
		WarmupTopic:   "test",
		MaxSimSeconds: 500,
		Public:        true,
	}
	if mappingImage {
		track.MappingImage = ptr.String("mapping-image")
	}
	trackService := tfake.NewService()
	trackService.On("Get", "test", 0, 0).Return(track, error(nil))

	// Initialize fake simulation service
	svc := simfake.NewService()
	svc.On("Get", sim.GetGroupID()).Return(sim, nil)

	// Create SubT application service
	app := subtapp.NewServices(application.NewServices(svc, nil), trackService, nil)

	s := state.NewStartSimulation(p, app, sim.GetGroupID())
	store := actions.NewStore(s)

	orchestratorStore.On("Timeout").Return(1 * time.Second)
	orchestratorStore.On("PollFrequency").Return(1 * time.Second)
	orchestratorStore.On("Namespace").Return("default")

	return WaitForMappingServerPod.Run(store, nil, &actions.Deployment{CurrentJob: "test"}, s)
}

func (suite *waitMappingServerPodTestSuite) TestWaitForMappingServerPod() {
	gid := simulations.GroupID("aaaa-bbbb-cccc-dddd")
	sim := subtsimfake.NewSimulation(subtsimfake.SimulationConfig{
		GroupID: gid,
		Status:  simulations.StatusLaunchingPods,
		Kind:    simulations.SimSingle,
		Error:   nil,
		Image:   "test",
		Track:   "test",
		Robots: []simulations.Robot{
			simfake.NewRobot("testA", "X1"),
			simfake.NewRobot("testB", "X2"),
			simfake.NewRobot("testC", "X3"),
		},
	})

	initialPod := &apiv1.Pod{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      subtapp.GetPodNameMappingServer(gid),
			Namespace: "default",
			Labels:    subtapp.GetPodLabelsMappingServer(gid, nil).Map(),
		},
		Spec: apiv1.PodSpec{},
		Status: apiv1.PodStatus{
			PodIP: "127.0.0.1",
		},
	}
	client := fake.NewSimpleClientset(initialPod)

	result, err := suite.waitForPod(sim, client, true)

	suite.Require().NoError(err)

	output, ok := result.(*state.StartSimulation)
	suite.Require().True(ok)

	suite.Require().Equal(gid, output.GroupID)
}

func (suite *waitMappingServerPodTestSuite) TestWaitForMappingServerNoPod() {
	gid := simulations.GroupID("aaaa-bbbb-cccc-dddd")
	sim := subtsimfake.NewSimulation(subtsimfake.SimulationConfig{
		GroupID: gid,
		Status:  simulations.StatusLaunchingPods,
		Kind:    simulations.SimSingle,
		Error:   nil,
		Image:   "test",
		Track:   "test",
		Robots: []simulations.Robot{
			simfake.NewRobot("testA", "X1"),
			simfake.NewRobot("testB", "X2"),
			simfake.NewRobot("testC", "X3"),
		},
	})

	client := fake.NewSimpleClientset()

	result, err := suite.waitForPod(sim, client, false)

	suite.Require().NoError(err)

	output, ok := result.(*state.StartSimulation)
	suite.Require().True(ok)

	suite.Require().Equal(gid, output.GroupID)
}
