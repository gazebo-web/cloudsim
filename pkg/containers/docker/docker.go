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
	Pull(ctx context.Context, repository, image string) error
	Create(ctx context.Context, name, image string, hostPortBinding nat.PortBinding, protocol string, port string) (containers.Container, error)
}

// docker is a Docker implementation.
type docker struct {
	CLI          *client.Client
	repositories map[string]string
}

func (d docker) Pull(ctx context.Context, repository, image string) error {
	r, err := d.CLI.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	io.Copy(os.Stdout, r)
	defer r.Close()
	return nil
}

// Create creates a docker container.
func (d docker) Create(ctx context.Context, name, image string, hostPortBinding nat.PortBinding, protocol string, port string) (containers.Container, error) {
	containerPort, err := nat.NewPort(protocol, port)

	portBindings := nat.PortMap{
		containerPort: []nat.PortBinding{hostPortBinding},
	}

	body, err := d.CLI.ContainerCreate(
		ctx,
		&container.Config{
			Image: image,
		},
		&container.HostConfig{
			PortBindings: portBindings,
		},
		nil,
		name,
	)
	if err != nil {
		return nil, err
	}

	return NewContainer(body.ID, name, image, portBindings, d.CLI), nil
}

func NewEngine() (Docker, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}
	return &docker{
		repositories: map[string]string{
			"dockerhub": "docker.io/library/",
		},
		CLI: cli,
	}, nil
}
