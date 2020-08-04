package localstack

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew_NotNil(t *testing.T) {
	l := New()
	assert.NotNil(t, l)
}

func TestNew_TypeOfStruct(t *testing.T) {
	l := New()
	_, ok := l.(*localStack)
	assert.True(t, ok)
}
