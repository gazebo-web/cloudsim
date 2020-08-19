package ingresses

import (
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/fake"
	"testing"
)

func TestNewIngresses(t *testing.T) {
	client := fake.NewSimpleClientset()
	m := NewIngresses(client)
	assert.IsType(t, &ingresses{}, m)
}
