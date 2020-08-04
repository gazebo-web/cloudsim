package docker

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewEngine_NotNil(t *testing.T) {
	engine, err := NewEngine()
	assert.NoError(t, err)
	assert.NotNil(t, engine)
}

func TestNewEngine_TypeOfStruct(t *testing.T) {
	engine, err := NewEngine()
	assert.NoError(t, err)
	_, ok := engine.(*docker)
	assert.True(t, ok)
}

func TestNewEngine_CLINotNil(t *testing.T) {
	engine, err := NewEngine()
	assert.NoError(t, err)
	e, ok := engine.(*docker)
	assert.True(t, ok)
	assert.NotNil(t, e.CLI)
}

func TestDocker_Pull(t *testing.T) {
	e, err := NewEngine()
	assert.NoError(t, err)
	ctx := context.Background()
	err = e.Pull(ctx, "dockerhub", "hello-world")
	assert.NoError(t, err)
}
