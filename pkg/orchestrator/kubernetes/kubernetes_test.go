package kubernetes

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/nodes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/pods"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/spdy"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestNewKubernetes(t *testing.T) {
	var no orchestrator.Nodes
	var po orchestrator.Pods
	var ig orchestrator.Ingresses
	var ru orchestrator.IngressRules

	ks := NewKubernetes(no, po, ig, ru)

	assert.NotNil(t, ks)
	assert.IsType(t, &k8s{}, ks)
}

func TestNewKubernetesWithNodeManager(t *testing.T) {
	var po orchestrator.Pods
	var ig orchestrator.Ingresses
	var ru orchestrator.IngressRules
	client := fake.NewSimpleClientset()
	no := nodes.NewNodes(client)
	ks := NewKubernetes(no, po, ig, ru)

	assert.NotNil(t, ks)
	assert.IsType(t, &k8s{}, ks)
	assert.NotNil(t, ks.Nodes())
}

func TestNewKubernetesWithPodManager(t *testing.T) {
	var no orchestrator.Nodes
	var ig orchestrator.Ingresses
	var ru orchestrator.IngressRules

	client := fake.NewSimpleClientset()
	fakeSpdy := spdy.NewSPDYFakeInitializer()
	po := pods.NewPods(client, fakeSpdy)

	ks := NewKubernetes(no, po, ig, ru)

	assert.NotNil(t, ks)
	assert.IsType(t, &k8s{}, ks)
	assert.NotNil(t, ks.Pods())
}
