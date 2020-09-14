package jobs

import (
	"context"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/pods"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/spdy"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	sfake "gitlab.com/ignitionrobotics/web/cloudsim/pkg/store/fake"
	"gitlab.com/ignitionrobotics/web/ign-go"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestWaitForGazeboServerPod(t *testing.T) {
	ctx := context.Background()
	logger := ign.NewLoggerNoRollbar("TestWaitForGazeboServerPod", ign.VerbosityDebug)
	storeOrchestrator := sfake.NewFakeOrchestrator()
	sfake.NewFakeStore(nil, storeOrchestrator, nil)
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
		Spec:   apiv1.PodSpec{},
		Status: apiv1.PodStatus{},
	}
	client := fake.NewSimpleClientset(initialPod)
	po := pods.NewPods(client, spdyInit, logger)
	ks := kubernetes.NewCustomKubernetes(kubernetes.Config{
		Nodes:           nil,
		Pods:            po,
		Ingresses:       nil,
		IngressRules:    nil,
		Services:        nil,
		NetworkPolicies: nil,
	})

	p := platform.NewPlatform(nil, nil, ks, nil)
	ctx = context.WithValue(ctx, "cloudsim_platform", p)

	input := &StartSimulationData{
		GroupID: "aaaa-bbbb-cccc-dddd",
		GazeboServerPod: orchestrator.NewResource(
			"test",
			"default",
			orchestrator.NewSelector(map[string]string{
				"app": "test",
			}),
		),
	}

	result, err := WaitForGazeboServerPod.Run(ctx, nil, &actions.Deployment{
		CurrentJob: "test",
	}, input)
	assert.NoError(t, err)

	output, ok := result.(*StartSimulationData)
	assert.True(t, ok)

	assert.Equal(t, input.GroupID, output.GroupID)
	assert.NotNil(t, output.GazeboServerPod)
}
