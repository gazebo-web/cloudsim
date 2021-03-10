package jobs

import (
	"github.com/stretchr/testify/assert"
	subtapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	services "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/services/implementations/kubernetes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/implementations/kubernetes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	sfake "gitlab.com/ignitionrobotics/web/cloudsim/pkg/store/implementations/fake"
	"gitlab.com/ignitionrobotics/web/ign-go"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kfake "k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestRemoveWebsocketService(t *testing.T) {
	logger := ign.NewLoggerNoRollbar("TestRemoveWebsocketService", ign.VerbosityDebug)

	// Mock store orchestrator
	storeOrchestrator := sfake.NewFakeOrchestrator()
	storeOrchestrator.On("Namespace").Return("default")

	// Initialize the store with the mocked store
	fakeStore := sfake.NewFakeStore(nil, storeOrchestrator, nil)

	// Define the simulation group ID used  by the job.
	gid := simulations.GroupID("aaaa-bbbb-cccc-dddd")

	// Initialize a new fake kubernetes client.
	client := kfake.NewSimpleClientset(&apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      subtapp.GetServiceNameWebsocket(gid),
			Namespace: "default",
		},
	})

	// Initialize a kubernetes service manager
	kss := services.NewServices(client, logger)

	// Initialize a new cluster component with the kubernetes service manager
	ks := kubernetes.NewCustomKubernetes(kubernetes.Config{Services: kss})

	// Initialize a new platform with the orchestrator component and the mocked store.
	p := platform.NewPlatform(platform.Components{
		Cluster: ks,
		Store:   fakeStore,
	})

	s := state.NewStopSimulation(p, nil, gid)

	store := actions.NewStore(s)

	_, err := RemoveWebsocketService.Run(store, nil, nil, s)

	assert.NoError(t, err)
}
