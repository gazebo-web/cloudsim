package jobs

import (
	"github.com/stretchr/testify/assert"
	subtapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	network "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/network/implementations/kubernetes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/implementations/kubernetes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	simfake "gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations/fake"
	sfake "gitlab.com/ignitionrobotics/web/cloudsim/pkg/store/implementations/fake"
	"gitlab.com/ignitionrobotics/web/ign-go"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kfake "k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestRemoveNetPols(t *testing.T) {
	logger := ign.NewLoggerNoRollbar("TestRemoveNetPols", ign.VerbosityDebug)

	// Set up store
	storeIgnition := sfake.NewFakeIgnition()
	storeOrchestrator := sfake.NewFakeOrchestrator()
	fakeStore := sfake.NewFakeStore(nil, storeOrchestrator, storeIgnition)

	// Mock orchestrator store methods for this test
	storeOrchestrator.On("Namespace").Return("default")

	// Create a GetGroupID for testing.
	gid := simulations.GroupID("aaaa-bbbb-cccc-dddd")

	client := kfake.NewSimpleClientset(
		&networkingv1.NetworkPolicy{
			ObjectMeta: metav1.ObjectMeta{
				Name:      subtapp.GetPodNameGazeboServer(gid),
				Namespace: "default",
			},
		},
		&networkingv1.NetworkPolicy{
			ObjectMeta: metav1.ObjectMeta{
				Name:      subtapp.GetPodNameFieldComputer(gid, subtapp.GetRobotID(0)),
				Namespace: "default",
			},
		},
		&networkingv1.NetworkPolicy{
			ObjectMeta: metav1.ObjectMeta{
				Name:      subtapp.GetPodNameCommsBridge(gid, subtapp.GetRobotID(0)),
				Namespace: "default",
			},
		},
	)

	nm := network.NewNetworkPolicies(client, logger)
	ks := kubernetes.NewCustomKubernetes(kubernetes.Config{
		NetworkPolicies: nm,
	})

	// Set up platform using fake store and fake kubernetes component
	p := platform.NewPlatform(platform.Components{
		Cluster: ks,
		Store:   fakeStore,
	})

	// Initialize generic simulation service
	simservice := simfake.NewService()

	// Create a simulation for the given track
	robots := []simulations.Robot{
		simfake.NewRobot("testA", "X1"),
	}

	// Make the get method return the fake simulation
	simservice.On("GetRobots", gid).Return(robots, error(nil))

	// Create SubT application service
	app := subtapp.NewServices(application.NewServices(simservice, nil), nil, nil)

	// Create new state: Start simulation state.
	initialState := state.NewStopSimulation(p, app, gid)
	s := actions.NewStore(&initialState)

	_, err := RemoveNetworkPolicies.Run(s, nil, nil, initialState)

	assert.NoError(t, err)

	_, err = RemoveNetworkPolicies.Run(s, nil, nil, initialState)

	assert.Error(t, err)
}
