package kubernetes

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestNewKubernetes(t *testing.T) {
	var nm orchestrator.NodeManager
	var pm orchestrator.PodManager

	ks := NewKubernetes(nm, pm)

	assert.NotNil(t, ks)
	assert.IsType(t, &k8s{}, ks)
}

func TestNewKubernetesWithNodeManager(t *testing.T) {
	var pm orchestrator.PodManager
	client := fake.NewSimpleClientset()
	nm := NewNodeManager(client)
	ks := NewKubernetes(nm, pm)

	assert.NotNil(t, ks)
	assert.IsType(t, &k8s{}, ks)
	assert.NotNil(t, ks.Nodes())
}

func TestNewKubernetesWithPodManager(t *testing.T) {
	var nm orchestrator.NodeManager

	pm := NewPodManager()
	ks := NewKubernetes(nm, pm)

	assert.NotNil(t, ks)
	assert.IsType(t, &k8s{}, ks)
	assert.NotNil(t, ks.Pods())
}
