package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/containers"
	"io"
	"os"
)

// Docker represents a docker client.
type Docker interface {
	Pull(ctx context.Context, image string) error
	Create(ctx context.Context, input CreateContainerInput) (containers.Container, error)
}

// CreateContainerInput is used to group a set of required fields when creating a container.
type CreateContainerInput struct {
	Name         string
	Image        string
	Protocol     *string
	Port         *string
	PortBindings []nat.PortBinding
	EnvVars      containers.EnvVars
}

// docker is a Docker implementation.
type docker struct {
	CLI *client.Client
}

func (d docker) Pull(ctx context.Context, image string) error {
	r, err := d.CLI.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	io.Copy(os.Stdout, r)
	defer r.Close()
	return nil
}

// Create creates a new docker container.
func (d docker) Create(ctx context.Context, input CreateContainerInput) (containers.Container, error) {
	var portBindings nat.PortMap
	if input.Protocol != nil && input.Port != nil {
		containerPort, err := nat.NewPort(*input.Protocol, *input.Port)
		if err != nil {
			return nil, err
		}
		portBindings[containerPort] = input.PortBindings
	}

	body, err := d.CLI.ContainerCreate(
		ctx,
		&container.Config{
			Image: input.Image,
			Env:   input.EnvVars.ToSlice(),
		},
		&container.HostConfig{
			PortBindings: portBindings,
		},
		nil,
		input.Name,
	)
	if err != nil {
		return nil, err
	}

	return NewContainer(body.ID, input.Name, input.Image, portBindings, input.EnvVars, d.CLI), nil
}

// NewClient initializes a new Docker client.
func NewClient() (Docker, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}
	return &docker{
		CLI: cli,
	}, nil
}
