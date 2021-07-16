package jobs

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	kubernetesServices "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/services/implementations/kubernetes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/implementations/kubernetes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	sfake "gitlab.com/ignitionrobotics/web/cloudsim/pkg/store/implementations/fake"
	"gitlab.com/ignitionrobotics/web/ign-go"
	kfake "k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestLaunchWebsocketService(t *testing.T) {
	// Initialize logger
	logger := ign.NewLoggerNoRollbar("TestLaunchWebsocketService", ign.VerbosityDebug)

	// Mock store orchestrator
	storeOrchestrator := sfake.NewFakeOrchestrator()
	storeOrchestrator.On("Namespace").Return("default")

	// Initialize the store with the mocked store
	fakeStore := sfake.NewFakeStore(nil, storeOrchestrator, nil, nil)

	// Define the simulation group ID used  by the job.
	gid := simulations.GroupID("aaaa-bbbb-cccc-dddd")

	// Initialize a new fake kubernetes client.
	client := kfake.NewSimpleClientset()

	// Initialize a kubernetes service manager
	kss := kubernetesServices.NewServices(client, logger)

	// Initialize a new cluster component with the kubernetes service manager
	ks := kubernetes.NewCustomKubernetes(kubernetes.Config{Services: kss})

	// Initialize a new platform with the orchestrator component and the mocked store.
	p, _ := platform.NewPlatform("test", platform.Components{
		Cluster: ks,
		Store:   fakeStore,
	})

	// Create new state: Start simulation state.
	s := state.NewStartSimulation(p, nil, gid)

	// Set up action store
	store := actions.NewStore(s)

	_, err := LaunchWebsocketService.Run(store, nil, nil, s)

	assert.NoError(t, err)
}
