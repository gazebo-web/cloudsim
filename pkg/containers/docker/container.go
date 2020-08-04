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
}

func (d dockerContainer) ID() string {
	return d.id
}

func (d dockerContainer) Name() string {
	return d.name
}

func (d dockerContainer) Image() string {
	return d.image
}

func (d dockerContainer) Start() {
	d.cli.ContainerStart(context.TODO(), d.id, types.ContainerStartOptions{})
}

func (d dockerContainer) Stop() {
	timeout := time.Second
	d.cli.ContainerStop(context.TODO(), d.id, &timeout)
}

func (d dockerContainer) Remove() {
	d.cli.ContainerRemove(context.TODO(), d.id, types.ContainerRemoveOptions{})
}

func NewContainer(ID, name, image string, ports nat.PortMap, cli *client.Client) containers.Container {
	return &dockerContainer{
		id:    ID,
		name:  name,
		image: image,
		ports: ports,
		cli:   cli,
	}
}
