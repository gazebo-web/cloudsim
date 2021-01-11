package jobs

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/pods"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/spdy"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/ign-go"
	v1 "k8s.io/api/core/v1"
	kfake "k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestGetGazeboServerIPFailsWhenNotSet(t *testing.T) {
	logger := ign.NewLoggerNoRollbar("TestGetGazeboServerIPFailsWhenNotSet", ign.VerbosityDebug)

	spdyInit := spdy.NewSPDYFakeInitializer()

	pod := &v1.Pod{
		Status: v1.PodStatus{
			PodIP: "",
		},
	}

	client := kfake.NewSimpleClientset(pod)
	po := pods.NewPods(client, spdyInit, logger)

	ks := kubernetes.NewCustomKubernetes(kubernetes.Config{
		Pods: po,
	})

	p := platform.NewPlatform(platform.Components{
		Cluster: ks,
	})

	gid := simulations.GroupID("test-12345")

	s := state.NewStartSimulation(p, nil, gid)

	store := actions.NewStore(s)

	_, err := GetGazeboIP.Run(store, nil, nil, s)

	assert.Error(t, err)
}
