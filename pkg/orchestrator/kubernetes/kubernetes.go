package kubernetes

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/ingresses"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/ingresses/rules"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/network"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/nodes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/pods"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/services"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/spdy"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

// k8s is a orchestrator.Cluster implementation.
type k8s struct {
	// nodes has a reference to an orchestrator.Nodes implementation.
	nodes orchestrator.Nodes

	// pods has a reference to an orchestrator.Pods implementation.
	pods orchestrator.Pods

	// ingressRules has a reference to an orchestrator.IngressRules implementation.
	ingressRules orchestrator.IngressRules

	// services has a reference to an orchestrator.Services implementation.
	services orchestrator.Services

	// ingresses has a reference to an orchestrator.Ingresses implementation.
	ingresses orchestrator.Ingresses

	// networkPolicies has a reference to an orchestrator.NetworkPolicies implementation.
	networkPolicies orchestrator.NetworkPolicies

	// extensions has a reference to an orchestrator.Extensions implementation.
	extensions orchestrator.Extensions
}

// IngressRules returns the Kubernetes orchestrator.IngressRules implementation.
func (k *k8s) IngressRules() orchestrator.IngressRules {
	return k.ingressRules
}

// Nodes returns the Kubernetes orchestrator.Nodes implementation.
func (k *k8s) Nodes() orchestrator.Nodes {
	return k.nodes
}

// Pods returns the Kubernetes orchestrator.Pods implementation.
func (k *k8s) Pods() orchestrator.Pods {
	return k.pods
}

// Services returns the Kubernetes orchestrator.Services implementation.
func (k *k8s) Services() orchestrator.Services {
	return k.services
}

// Ingresses returns the Kubernetes orchestrator.Ingresses implementation.
func (k *k8s) Ingresses() orchestrator.Ingresses {
	return k.ingresses
}

// NetworkPolicies returns the Kubernetes orchestrator.NetworkPolicies implementation.
func (k *k8s) NetworkPolicies() orchestrator.NetworkPolicies {
	return k.networkPolicies
}

func (k *k8s) Extensions() orchestrator.Extensions {
	if k.extensions == nil {
		panic("no extensions have been added to the cluster implementation")
	}
	return k.extensions
}

// Config is used to group the inputs for NewCustomKubernetes.
// It includes all the needed subcomponents for Kubernetes.
type Config struct {
	Nodes           orchestrator.Nodes
	Pods            orchestrator.Pods
	Ingresses       orchestrator.Ingresses
	IngressRules    orchestrator.IngressRules
	Services        orchestrator.Services
	NetworkPolicies orchestrator.NetworkPolicies
	Extensions      orchestrator.Extensions
}

// NewCustomKubernetes returns a orchestrator.Cluster implementation using Kubernetes.
// All the subcomponents provided by the Config should be already initialized.
func NewCustomKubernetes(config Config) orchestrator.Cluster {
	return &k8s{
		nodes:           config.Nodes,
		pods:            config.Pods,
		ingresses:       config.Ingresses,
		ingressRules:    config.IngressRules,
		services:        config.Services,
		networkPolicies: config.NetworkPolicies,
		extensions:      config.Extensions,
	}
}

// NewDefaultKubernetes initializes the set of Kubernetes subcomponents using
// the given kubernetes client api, spdy initializer and logger.
func NewDefaultKubernetes(api kubernetes.Interface, spdyInit spdy.Initializer, logger ign.Logger) orchestrator.Cluster {
	return &k8s{
		nodes:           nodes.NewNodes(api, logger),
		pods:            pods.NewPods(api, spdyInit, logger),
		ingressRules:    rules.NewIngressRules(api, logger),
		services:        services.NewServices(api, logger),
		ingresses:       ingresses.NewIngresses(api, logger),
		networkPolicies: network.NewNetworkPolicies(api, logger),
		extensions:      nil,
	}
}

// NewFakeKubernetes initializes the set of Kubernetes subcomponents using fake implementations.
func NewFakeKubernetes(logger ign.Logger) orchestrator.Cluster {
	api := fake.NewSimpleClientset()
	spdyInit := spdy.NewSPDYFakeInitializer()
	return &k8s{
		nodes:           nodes.NewNodes(api, logger),
		pods:            pods.NewPods(api, spdyInit, logger),
		ingressRules:    rules.NewIngressRules(api, logger),
		services:        services.NewServices(api, logger),
		ingresses:       ingresses.NewIngresses(api, logger),
		networkPolicies: network.NewNetworkPolicies(api, logger),
		extensions:      nil,
	}
}

// InitializeKubernetes initializes a new Kubernetes orchestrator.
func InitializeKubernetes(logger ign.Logger) (orchestrator.Cluster, error) {
	config, err := GetConfig()
	if err != nil {
		return nil, err
	}
	client, err := NewAPI(config)
	if err != nil {
		return nil, err
	}
	spdyInit := spdy.NewSPDYInitializer(config)
	return NewDefaultKubernetes(client, spdyInit, logger), nil
}
