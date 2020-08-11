package kubernetes

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
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
	var nm orchestrator.NodeManager
	var pm orchestrator.PodManager

	ks := NewKubernetes(nm, pm)

	assert.NotNil(t, ks)
	assert.IsType(t, &k8s{}, ks)
	assert.NotNil(t, ks.Nodes())
}
