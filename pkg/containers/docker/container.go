package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/containers"
	"time"
)

type dockerContainer struct {
	id    string
	name  string
	image string
	ports nat.PortMap
	cli   *client.Client
	envs  map[string]string
}

// ID returns the container's ID.
func (d dockerContainer) ID() string {
	return d.id
}

// Name returns the container's name.
func (d dockerContainer) Name() string {
	return d.name
}

// Image returns the image's name that were used to create the container.
func (d dockerContainer) Image() string {
	return d.image
}

// EnvVars returns the container's environment variables.
func (d dockerContainer) EnvVars() containers.EnvVars {
	panic("implement me")
}

// Start starts the container.
func (d dockerContainer) Start() error {
	return d.cli.ContainerStart(context.TODO(), d.id, types.ContainerStartOptions{})
}

// Stop stops the container.
func (d dockerContainer) Stop() error {
	timeout := time.Second
	return d.cli.ContainerStop(context.TODO(), d.id, &timeout)
}

// Remove removes the container.
func (d dockerContainer) Remove() error {
	return d.cli.ContainerRemove(context.TODO(), d.id, types.ContainerRemoveOptions{})
}

// NewContainer initializes a new containers.Container implementation using Docker.
func NewContainer(ID, name, image string, ports nat.PortMap, envVars containers.EnvVars, cli *client.Client) containers.Container {
	return &dockerContainer{
		id:    ID,
		name:  name,
		image: image,
		ports: ports,
		cli:   cli,
		envs:  envVars,
	}
}
