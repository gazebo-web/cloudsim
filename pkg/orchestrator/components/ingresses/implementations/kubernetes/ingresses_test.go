package kubernetes

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/ignitionrobotics/web/ign-go/v5"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestNewIngresses(t *testing.T) {
	client := fake.NewSimpleClientset()
	m := NewIngresses(client, ign.NewLoggerNoRollbar("TestIngress", ign.VerbosityDebug))
	assert.IsType(t, &kubernetesIngresses{}, m)
}
