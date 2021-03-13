package jobs

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	kubernetesPods "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/pods/implementations/kubernetes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/spdy"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/implementations/kubernetes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/resource"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	sfake "gitlab.com/ignitionrobotics/web/cloudsim/pkg/store/implementations/fake"
	"gitlab.com/ignitionrobotics/web/ign-go"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
	"time"
)

func TestWaitForGazeboServerPod(t *testing.T) {
	logger := ign.NewLoggerNoRollbar("TestWaitForGazeboServerPod", ign.VerbosityDebug)
	storeMachines := sfake.NewFakeMachines()
	orchestratorStore := sfake.NewFakeOrchestrator()
	fakeStore := sfake.NewFakeStore(storeMachines, orchestratorStore, nil)
	spdyInit := spdy.NewSPDYFakeInitializer()

	initialPod := &apiv1.Pod{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "default",
			Labels: map[string]string{
				"app": "test",
			},
		},
		Spec: apiv1.PodSpec{},
		Status: apiv1.PodStatus{
			PodIP: "127.0.0.1",
		},
	}
	client := fake.NewSimpleClientset(initialPod)
	po := kubernetesPods.NewPods(client, spdyInit, logger)
	ks := kubernetes.NewCustomKubernetes(kubernetes.Config{
		Nodes:           nil,
		Pods:            po,
		Ingresses:       nil,
		IngressRules:    nil,
		Services:        nil,
		NetworkPolicies: nil,
	})

	p := platform.NewPlatform(platform.Components{
		Cluster: ks,
		Store:   fakeStore,
	})

	gid := simulations.GroupID("aaaa-bbbb-cccc-dddd")
	s := state.NewStartSimulation(p, nil, gid)
	store := actions.NewStore(s)

	// Mock gazebo server pod
	s.GazeboServerPod = resource.NewResource(
		"test",
		"default",
		resource.NewSelector(map[string]string{
			"app": "test",
		}),
	)
	store.SetState(s)

	storeMachines.On("Timeout").Return(1 * time.Second)
	storeMachines.On("PollFrequency").Return(1 * time.Second)
	orchestratorStore.On("Namespace").Return("default")

	result, err := WaitForGazeboServerPod.Run(store, nil, &actions.Deployment{CurrentJob: "test"}, s)
	assert.NoError(t, err)

	output, ok := result.(*state.StartSimulation)
	assert.True(t, ok)

	assert.Equal(t, gid, output.GroupID)
	assert.NotNil(t, output.GazeboServerPod)
}
