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
	var nm orchestrator.NodeManager
	var pm orchestrator.PodManager

	ks := NewKubernetes(nm, pm)

	assert.NotNil(t, ks)
	assert.IsType(t, &k8s{}, ks)
}

func TestNewKubernetesWithNodeManager(t *testing.T) {
	var pm orchestrator.PodManager
	client := fake.NewSimpleClientset()
	nm := nodes.NewManager(client)
	ks := NewKubernetes(nm, pm)

	assert.NotNil(t, ks)
	assert.IsType(t, &k8s{}, ks)
	assert.NotNil(t, ks.Nodes())
}

func TestNewKubernetesWithPodManager(t *testing.T) {
	var nm orchestrator.NodeManager

	client := fake.NewSimpleClientset()
	fakeSpdy := spdy.NewSPDYFakeInitializer()
	pm := pods.NewManager(client, fakeSpdy)

	ks := NewKubernetes(nm, pm)

	assert.NotNil(t, ks)
	assert.IsType(t, &k8s{}, ks)
	assert.NotNil(t, ks.Pods())
}
