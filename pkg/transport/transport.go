package transport

import (
	"github.com/caarlos0/env"
	igntransport "gitlab.com/ignitionrobotics/web/cloudsim/third_party/ign-transport"
)

// config
type config struct {
	Topic string `env:"IGN_TRANSPORT_TEST_TOPIC" envDefault:"/foo"`
}

// Transport
type Transport struct {
	Node  *igntransport.GoIgnTransportNode
	Topic string
}

// New
func New() (*Transport, error) {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	node, err := igntransport.NewIgnTransportNode(nil)
	if err != nil {
		return nil, err
	}
	return &Transport{
		Node:  node,
		Topic: cfg.Topic,
	}, nil
}

// Stop
func (t *Transport) Stop() {
	t.Node.Free()
}
