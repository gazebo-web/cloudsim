package localstack

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew_NotNil(t *testing.T) {
	l, err := New("localhost", 4566)
	assert.NotNil(t, l)
	assert.NoError(t, err)
}

func TestNew_TypeOfStruct(t *testing.T) {
	l, err := New("localhost", 4566)
	_, ok := l.(*localStack)
	assert.True(t, ok)
	assert.NoError(t, err)
}

func TestStart_NoError(t *testing.T) {
	l, err := New("localhost", 4566)
	assert.NoError(t, err)

	err = l.Start(context.Background())
	assert.NoError(t, err)
}
