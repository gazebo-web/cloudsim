package ingress

import (
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestNewManager(t *testing.T) {
	client := fake.NewSimpleClientset()
	m := NewIngresses(client)
	assert.IsType(t, &ingresses{}, m)
}
