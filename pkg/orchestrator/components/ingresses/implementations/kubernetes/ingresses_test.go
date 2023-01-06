package kubernetes

import (
	"github.com/gazebo-web/gz-go/v7"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestNewIngresses(t *testing.T) {
	client := fake.NewSimpleClientset()
	m := NewIngresses(client, gz.NewLoggerNoRollbar("TestIngress", gz.VerbosityDebug))
	assert.IsType(t, &kubernetesIngresses{}, m)
}
