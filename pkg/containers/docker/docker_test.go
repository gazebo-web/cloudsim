package docker

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewEngine_NotNil(t *testing.T) {
	cli, err := NewClient()
	assert.NoError(t, err)
	assert.NotNil(t, cli)
}

func TestNewEngine_TypeOfStruct(t *testing.T) {
	cli, err := NewClient()
	assert.NoError(t, err)
	_, ok := cli.(*docker)
	assert.True(t, ok)
}

func TestNewEngine_CLINotNil(t *testing.T) {
	cli, err := NewClient()
	assert.NoError(t, err)
	e, ok := cli.(*docker)
	assert.True(t, ok)
	assert.NotNil(t, e.CLI)
}

func TestDocker_Pull(t *testing.T) {
	e, err := NewClient()
	assert.NoError(t, err)
	ctx := context.Background()
	err = e.Pull(ctx, "hello-world")
	assert.NoError(t, err)
}

func TestDocker_CreateAndRemove(t *testing.T) {
	cli, err := NewClient()
	assert.NoError(t, err)
	ctx := context.Background()
	err = cli.Pull(ctx, "hello-world")
	assert.NoError(t, err)
	c, err := cli.Create(ctx, CreateContainerInput{
		Name:         "test",
		Image:        "hello-world",
		Protocol:     nil,
		Port:         nil,
		PortBindings: nil,
		EnvVars:      nil,
	})
	assert.NoError(t, err)
	assert.Equal(t, c.Name(), "test")
	c.Remove()
}

func TestDocker_CreateStartAndRemove(t *testing.T) {
	cli, err := NewClient()
	assert.NoError(t, err)
	ctx := context.Background()
	err = cli.Pull(ctx, "hello-world")
	assert.NoError(t, err)
	c, err := cli.Create(ctx, CreateContainerInput{
		Name:         "test",
		Image:        "hello-world",
		Protocol:     nil,
		Port:         nil,
		PortBindings: nil,
		EnvVars:      nil,
	})
	assert.NoError(t, err)
	assert.Equal(t, c.Name(), "test")
	err = c.Start()
	assert.NoError(t, err)
	err = c.Stop()
	assert.NoError(t, err)
	err = c.Remove()
	assert.NoError(t, err)
}

func TestDocker_ContainerStartStopAndRemove(t *testing.T) {
	cli, err := NewClient()
	assert.NoError(t, err)
	ctx := context.Background()
	err = cli.Pull(ctx, "hello-world")
	assert.NoError(t, err)
	c, err := cli.Create(ctx, CreateContainerInput{
		Name:         "test",
		Image:        "hello-world",
		Protocol:     nil,
		Port:         nil,
		PortBindings: nil,
		EnvVars:      nil,
	})
	assert.NoError(t, err)
	assert.Equal(t, c.Name(), "test")
	err = c.Start()
	assert.NoError(t, err)
	err = c.Stop()
	assert.NoError(t, err)
	err = c.Remove()
	assert.NoError(t, err)
}
