package jobs

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	kubernetesPods "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/pods/implementations/kubernetes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/spdy"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/implementations/kubernetes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	sfake "gitlab.com/ignitionrobotics/web/cloudsim/pkg/store/implementations/fake"
	"gitlab.com/ignitionrobotics/web/ign-go"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kfake "k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestGetGazeboServerIPFailsWhenNotSet(t *testing.T) {
	logger := ign.NewLoggerNoRollbar("TestGetGazeboServerIPFailsWhenNotSet", ign.VerbosityDebug)

	storeOrchestrator := sfake.NewFakeOrchestrator()
	storeOrchestrator.On("Namespace").Return("default")
	fakeStore := sfake.NewFakeStore(nil, storeOrchestrator, nil, nil)

	spdyInit := spdy.NewSPDYFakeInitializer()

	pod := &v1.Pod{
		Status: v1.PodStatus{
			PodIP: "",
		},
	}

	client := kfake.NewSimpleClientset(pod)
	po := kubernetesPods.NewPods(client, spdyInit, logger)

	ks := kubernetes.NewCustomKubernetes(kubernetes.Config{
		Pods: po,
	})

	p, _ := platform.NewPlatform("test", platform.Components{
		Cluster: ks,
		Store:   fakeStore,
	})

	gid := simulations.GroupID("test-12345")

	s := state.NewStartSimulation(p, nil, gid)

	store := actions.NewStore(s)

	_, err := GetGazeboIP.Run(store, nil, nil, s)

	assert.Error(t, err)
}

func TestGetGazeboIP(t *testing.T) {
	logger := ign.NewLoggerNoRollbar("TestGetGazeboIP", ign.VerbosityDebug)

	storeOrchestrator := sfake.NewFakeOrchestrator()
	storeOrchestrator.On("Namespace").Return("default")
	fakeStore := sfake.NewFakeStore(nil, storeOrchestrator, nil, nil)

	spdyInit := spdy.NewSPDYFakeInitializer()

	ip := "127.0.0.1"
	gid := simulations.GroupID("test-12345")

	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      application.GetPodNameGazeboServer(gid),
			Namespace: "default",
		},
		Status: v1.PodStatus{
			PodIP: ip,
		},
	}

	client := kfake.NewSimpleClientset(pod)
	po := kubernetesPods.NewPods(client, spdyInit, logger)

	ks := kubernetes.NewCustomKubernetes(kubernetes.Config{
		Pods: po,
	})

	p, _ := platform.NewPlatform("test", platform.Components{
		Cluster: ks,
		Store:   fakeStore,
	})

	s := state.NewStartSimulation(p, nil, gid)

	store := actions.NewStore(s)

	out, err := GetGazeboIP.Run(store, nil, nil, s)

	require.NoError(t, err)

	result, ok := out.(*state.StartSimulation)

	assert.True(t, ok)
	assert.Equal(t, ip, result.GazeboServerIP)
}
