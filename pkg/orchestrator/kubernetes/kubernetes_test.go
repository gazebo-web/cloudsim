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
	var nm orchestrator.Nodes
	var pm orchestrator.Pods
	var rm orchestrator.IngressRules

	ks := NewKubernetes(nm, pm, rm)

	assert.NotNil(t, ks)
	assert.IsType(t, &k8s{}, ks)
}

func TestNewKubernetesWithNodeManager(t *testing.T) {
	var pm orchestrator.Pods
	var rm orchestrator.IngressRules
	client := fake.NewSimpleClientset()
	nm := nodes.NewNodes(client)
	ks := NewKubernetes(nm, pm, rm)

	assert.NotNil(t, ks)
	assert.IsType(t, &k8s{}, ks)
	assert.NotNil(t, ks.Nodes())
}

func TestNewKubernetesWithPodManager(t *testing.T) {
	var nm orchestrator.Nodes
	var rm orchestrator.IngressRules

	client := fake.NewSimpleClientset()
	fakeSpdy := spdy.NewSPDYFakeInitializer()
	pm := pods.NewPods(client, fakeSpdy)

	ks := NewKubernetes(nm, pm, rm)

	assert.NotNil(t, ks)
	assert.IsType(t, &k8s{}, ks)
	assert.NotNil(t, ks.Pods())
}
