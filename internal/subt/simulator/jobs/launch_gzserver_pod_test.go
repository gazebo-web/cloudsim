package jobs

import (
	"github.com/stretchr/testify/require"
	subtapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulations/fake"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/tracks"
	tfake "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/tracks/fake"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/pods"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/spdy"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	simfake "gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations/fake"
	sfake "gitlab.com/ignitionrobotics/web/cloudsim/pkg/store/fake"
	"gitlab.com/ignitionrobotics/web/ign-go"
	kfake "k8s.io/client-go/kubernetes/fake"
	"testing"
	"time"
)

func TestLaunchGazeboServerPod(t *testing.T) {
	// Set up logger
	logger := ign.NewLoggerNoRollbar("TestLaunchGazeboServerPod", ign.VerbosityDebug)

	// Set up store
	storeIgnition := sfake.NewFakeIgnition()
	storeOrchestrator := sfake.NewFakeOrchestrator()
	fakeStore := sfake.NewFakeStore(nil, storeOrchestrator, storeIgnition)

	// Mock stores methods used on this test
	storeIgnition.On("GazeboServerLogsPath").Return("/tmp/test")
	storeIgnition.On("IP").Return("127.0.0.1")
	storeIgnition.On("Verbosity").Return("0")

	storeOrchestrator.On("Namespace").Return("default")
	storeOrchestrator.On("TerminationGracePeriod").Return(time.Second)
	storeOrchestrator.On("Nameservers").Return([]string{"8.8.8.8", "8.8.4.4"})

	// Set up SPDY initializer with fake implementation
	spdyInit := spdy.NewSPDYFakeInitializer()

	// Set up kubernetes component
	client := kfake.NewSimpleClientset()
	po := pods.NewPods(client, spdyInit, logger)
	ks := kubernetes.NewCustomKubernetes(kubernetes.Config{
		Nodes:           nil,
		Pods:            po,
		Ingresses:       nil,
		IngressRules:    nil,
		Services:        nil,
		NetworkPolicies: nil,
	})

	// Set up platform using fake store and kubernetes component
	p := platform.NewPlatform(nil, nil, ks, fakeStore)

	// Initialize generic simulation service
	simservice := simfake.NewService()

	// Create a GroupID for testing.
	gid := simulations.GroupID("aaaa-bbbb-cccc-dddd")

	trackName := "Cave Circuit World 1"

	// Create a simulation
	sim := fake.NewSimulation(fake.SimulationConfig{
		GroupID: gid,
		Status:  simulations.StatusRunning,
		Kind:    simulations.SimSingle,
		Error:   nil,
		Image:   "test.org",
		Track:   trackName,
	})

	// Make the get method return the fake simulation when using
	simservice.On("Get", gid).Return(sim, error(nil))

	// Initialize tracks service
	trackService := tfake.NewService()

	trackService.On("Get", trackName).Return(&tracks.Track{
		Name:          trackName,
		Image:         "world-image.org/image",
		BridgeImage:   "bridge-image.org/image",
		StatsTopic:    "test",
		WarmupTopic:   "test",
		MaxSimSeconds: 500,
		Public:        true,
	}, error(nil))

	// Create SubT application service
	app := subtapp.NewServices(application.NewServices(simservice), trackService)

	// Create new state: Start simulation state.
	s := state.NewStartSimulation(p, app, gid)

	// Set up action store
	store := actions.NewStore(s)

	// Run job
	_, err := LaunchGazeboServerPod.Run(store, nil, &actions.Deployment{CurrentJob: "test"}, s)

	// Check if there are any errors.
	require.NoError(t, err)
}
