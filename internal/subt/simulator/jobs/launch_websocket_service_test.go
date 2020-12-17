package jobs

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/services"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	sfake "gitlab.com/ignitionrobotics/web/cloudsim/pkg/store/fake"
	"gitlab.com/ignitionrobotics/web/ign-go"
	kfake "k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestLaunchWebsocketService(t *testing.T) {
	logger := ign.NewLoggerNoRollbar("TestLaunchWebsocketService", ign.VerbosityDebug)

	storeOrchestrator := sfake.NewFakeOrchestrator()
	storeOrchestrator.On("Namespace").Return("default")

	fakeStore := sfake.NewFakeStore(nil, storeOrchestrator, nil)

	gid := simulations.GroupID("aaaa-bbbb-cccc-dddd")

	client := kfake.NewSimpleClientset()

	kss := services.NewServices(client, logger)

	ks := kubernetes.NewCustomKubernetes(kubernetes.Config{Services: kss})

	p := platform.NewPlatform(platform.Components{
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
