package gloo

import (
	"context"
	"errors"
	gatewayFake "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/client/clientset/versioned/fake"
	gateway "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/client/clientset/versioned/typed/gateway.solo.io/v1"
	glooFake "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/kube/client/clientset/versioned/fake"
	gloo "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/kube/client/clientset/versioned/typed/gloo.solo.io/v1"
	"gitlab.com/ignitionrobotics/web/ign-go"
	restclient "k8s.io/client-go/rest"
)

var (
	// ErrFailedToInitializeClientset is returned when the passed configuration options failed to properly configure a
	// clientset.
	ErrFailedToInitializeClientset = errors.New("did not initialize gloo clientset")
)

// Clientset contains the set of supported Kubernetes CRD Gloo interfaces.
type Clientset interface {
	// Gloo returns a Gloo client. This client can be used to interact with Gloo CRDs in a target Kubernetes cluster.
	// Available Gloo resources include Artifacts, Endpoints, Proxies, Secrets, Settings, Upstreams and Upstream Groups.
	Gloo() gloo.GlooV1Interface
	// Gateway returns a Gloo Gateway client. This client can be used to interact with Gloo CRDs in a target Kubernetes
	// cluster.
	// Available Gateway resources include Gateways, Route Tables and Virtual Services.
	Gateway() gateway.GatewayV1Interface
}

// clientset is used to access Gloo resources in a Kubernetes cluster.
type clientset struct {
	gloo.GlooV1Interface
	gateway.GatewayV1Interface
}

// Gloo returns a Gloo client. This client can be used to interact with Gloo CRDs in a target Kubernetes cluster.
func (gc *clientset) Gloo() gloo.GlooV1Interface {
	return gc.GlooV1Interface
}

// Gateway returns a Gloo Gateway client. This client can be used to interact with Gloo CRDs in a target Kubernetes
// cluster.
func (gc *clientset) Gateway() gateway.GatewayV1Interface {
	return gc.GatewayV1Interface
}

// ClientsetConfig is used to initialize a Gloo clientset,
type ClientsetConfig struct {
	// KubeConfig contains the target Kubernetes cluster's connection configuration.
	KubeConfig *restclient.Config
	// IsGoTest indicates if the clientset is being created for a test. If it is, a fake implementation will be
	// returned.
	IsGoTest bool
	// ConnectToCloudServices indicates that the clientset should be able to connect to cloud services and perform
	// persistent operations.
	ConnectToCloudServices bool
}

// newFakeClientset returns a fake Gloo Clientset for unit testing.
func newFakeClientset() Clientset {
	return &clientset{
		GlooV1Interface:    glooFake.NewSimpleClientset().GlooV1(),
		GatewayV1Interface: gatewayFake.NewSimpleClientset().GatewayV1(),
	}
}

// newClientset creates a new Gloo clientset.
func newClientset(kubeconfig *restclient.Config) (Clientset, error) {
	// Prepare the Gloo clientset
	glooClient, err := gloo.NewForConfig(kubeconfig)
	if err != nil {
		return nil, err
	}

	// Prepare the gateway clientset
	gatewayClient, err := gateway.NewForConfig(kubeconfig)
	if err != nil {
		return nil, err
	}

	// Create the clientset
	clientset := &clientset{
		GlooV1Interface:    glooClient,
		GatewayV1Interface: gatewayClient,
	}

	return clientset, nil
}

// NewClientset creates a new instance of the Gloo client.
func NewClientset(ctx context.Context, config *ClientsetConfig) (Clientset, error) {
	// Return the fake clientset if this is a test
	if config.IsGoTest {
		return newFakeClientset(), nil
	}

	// Validate the configuration
	// The clientset cannot be configured if a Kubernetes configuration was not provided
	if config.KubeConfig == nil {
		return nil, ErrFailedToInitializeClientset
	}
	// If the clientset is not initialized if it cannot connect to the cloud
	if !config.ConnectToCloudServices {
		return nil, ErrFailedToInitializeClientset
	}

	// Prepare the clientset
	gloo, err := newClientset(config.KubeConfig)
	if err != nil {
		ign.LoggerFromContext(ctx).Error("Failed to create a Gloo client", err)
		return nil, err
	}

	return gloo, nil
}
