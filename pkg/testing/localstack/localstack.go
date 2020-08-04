package localstack

import (
	"context"
	"github.com/docker/go-connections/nat"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/containers"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/containers/docker"
	"strconv"
)

var (
	netProtocol = "tcp"
)

const (
	dockerImage = "localstack/localstack"
	serviceList = "ec2,"
)

// LocalStack groups a set of methods to wrap LocalStack tool for testing purposes.
type LocalStack interface {
	Start(ctx context.Context) error
}

// localStack is a LocalStack implementation.
type localStack struct {
	host      string
	edgePort  int
	docker    docker.Docker
	container containers.Container
}

// Start starts LocalStack service.
func (l *localStack) Start(ctx context.Context) error {
	err := l.docker.Pull(ctx, dockerImage)
	if err != nil {
		return err
	}
	port := strconv.Itoa(l.edgePort)
	l.container, err = l.docker.Create(ctx, docker.CreateContainerInput{
		Name:     "localstack",
		Image:    dockerImage,
		Protocol: &netProtocol,
		Port:     &port,
		PortBindings: []nat.PortBinding{
			{
				HostIP:   l.host,
				HostPort: port,
			},
		},
		EnvVars: containers.EnvVars{
			"EDGE_PORT": port,
			"HOSTNAME":  l.host,
			"SERVICES":  serviceList,
		},
	})
	if err != nil {
		return err
	}
	err = l.container.Start()
	if err != nil {
		return err
	}
	return nil
}

// New initializes a new LocalStack implementation.
func New(host string, edgePort int) (LocalStack, error) {
	d, err := docker.NewClient()
	if err != nil {
		return nil, err
	}
	return &localStack{
		host:      host,
		edgePort:  edgePort,
		docker:    d,
		container: nil,
	}, nil
}
