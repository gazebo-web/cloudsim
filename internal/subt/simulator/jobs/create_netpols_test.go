package jobs

import (
	"github.com/stretchr/testify/assert"
	subtapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	kubernetesNetwork "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/network/implementations/kubernetes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/implementations/kubernetes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	simfake "gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations/fake"
	sfake "gitlab.com/ignitionrobotics/web/cloudsim/pkg/store/implementations/fake"
	"gitlab.com/ignitionrobotics/web/ign-go"
	kfake "k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestCreateNetPolsGazeboServer(t *testing.T) {
	logger := ign.NewLoggerNoRollbar("TestCreateNetPolsGazeboServer", ign.VerbosityDebug)

	// Set up store
	storeIgnition := sfake.NewFakeIgnition()
	storeOrchestrator := sfake.NewFakeOrchestrator()
	fakeStore := sfake.NewFakeStore(nil, storeOrchestrator, storeIgnition)

	// Mock ignition store methods for this test
	storeIgnition.On("IP").Return("127.0.0.1")

	// Mock orchestrator store methods for this test
	storeOrchestrator.On("Namespace").Return("default")

	client := kfake.NewSimpleClientset()

	nm := kubernetesNetwork.NewNetworkPolicies(client, logger)
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

	// Create a GetGroupID for testing.
	gid := simulations.GroupID("aaaa-bbbb-cccc-dddd")

	// Create a simulation for the given track
	robots := []simulations.Robot{
		simfake.NewRobot("testA", "X1"),
		simfake.NewRobot("testB", "X2"),
		simfake.NewRobot("testC", "X3"),
	}

	// Make the get method return the fake simulation
	simservice.On("GetRobots", gid).Return(robots, error(nil))

	// Create SubT application service
	app := subtapp.NewServices(application.NewServices(simservice, nil), nil, nil)

	// Create new state: Start simulation state.
	initialState := state.NewStartSimulation(p, app, gid)
	s := actions.NewStore(&initialState)

	_, err := CreateNetworkPolicyGazeboServer.Run(s, nil, nil, initialState)

	assert.NoError(t, err)

}

func TestCreateNetPolsCommsBridge(t *testing.T) {
	logger := ign.NewLoggerNoRollbar("CreateNetworkPolicyCommsBridges", ign.VerbosityDebug)

	// Set up store
	storeIgnition := sfake.NewFakeIgnition()
	storeOrchestrator := sfake.NewFakeOrchestrator()
	fakeStore := sfake.NewFakeStore(nil, storeOrchestrator, storeIgnition)

	// Mock ignition store methods for this test
	storeIgnition.On("IP").Return("127.0.0.1")

	// Mock orchestrator store methods for this test
	storeOrchestrator.On("Namespace").Return("default")

	client := kfake.NewSimpleClientset()

	nm := kubernetesNetwork.NewNetworkPolicies(client, logger)
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

	// Create a GetGroupID for testing.
	gid := simulations.GroupID("aaaa-bbbb-cccc-dddd")

	// Create a simulation for the given track
	robots := []simulations.Robot{
		simfake.NewRobot("testA", "X1"),
		simfake.NewRobot("testB", "X2"),
		simfake.NewRobot("testC", "X3"),
	}

	// Make the get method return the fake simulation
	simservice.On("GetRobots", gid).Return(robots, error(nil))

	// Create SubT application service
	app := subtapp.NewServices(application.NewServices(simservice, nil), nil, nil)

	// Create new state: Start simulation state.
	initialState := state.NewStartSimulation(p, app, gid)
	s := actions.NewStore(&initialState)

	_, err := CreateNetworkPolicyCommsBridges.Run(s, nil, nil, initialState)

	assert.NoError(t, err)

}

func TestCreateNetPolsFieldComputer(t *testing.T) {
	logger := ign.NewLoggerNoRollbar("CreateNetworkPolicyFieldComputers", ign.VerbosityDebug)

	// Set up store
	storeIgnition := sfake.NewFakeIgnition()
	storeOrchestrator := sfake.NewFakeOrchestrator()
	fakeStore := sfake.NewFakeStore(nil, storeOrchestrator, storeIgnition)

	// Mock ignition store methods for this test
	storeIgnition.On("IP").Return("127.0.0.1")

	// Mock orchestrator store methods for this test
	storeOrchestrator.On("Namespace").Return("default")

	client := kfake.NewSimpleClientset()

	nm := kubernetesNetwork.NewNetworkPolicies(client, logger)
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

	// Create a GetGroupID for testing.
	gid := simulations.GroupID("aaaa-bbbb-cccc-dddd")

	// Create a simulation for the given track
	robots := []simulations.Robot{
		simfake.NewRobot("testA", "X1"),
		simfake.NewRobot("testB", "X2"),
		simfake.NewRobot("testC", "X3"),
	}

	// Make the get method return the fake simulation
	simservice.On("GetRobots", gid).Return(robots, error(nil))

	// Create SubT application service
	app := subtapp.NewServices(application.NewServices(simservice, nil), nil, nil)

	// Create new state: Start simulation state.
	initialState := state.NewStartSimulation(p, app, gid)
	s := actions.NewStore(&initialState)

	_, err := CreateNetworkPolicyFieldComputers.Run(s, nil, nil, initialState)

	assert.NoError(t, err)

}
