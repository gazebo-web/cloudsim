package spdy

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"net/url"
)

// Initializer is used to initialize different SPDY executors.
// We created this interface to inject Kubernetes a way to mock an SPDY executor.
type Initializer interface {
	NewSPDYExecutor(method string, url *url.URL) (remotecommand.Executor, error)
}

// initializer is a Initializer implementation.
// It's a wrapper for the default Kubernetes remotecommand.NewSPDYExecutor implementation.
type initializer struct {
	config *rest.Config
}

// NewSPDYExecutor creates a new remotecommand.Executor.
func (i initializer) NewSPDYExecutor(method string, url *url.URL) (remotecommand.Executor, error) {
	exec, err := remotecommand.NewSPDYExecutor(i.config, method, url)
	if err != nil {
		return nil, err
	}
	return exec, nil
}

// NewSPDYInitializer initializes a new Initializer using the default implementation.
func NewSPDYInitializer(config *rest.Config) Initializer {
	return &initializer{
		config: config,
	}
}
