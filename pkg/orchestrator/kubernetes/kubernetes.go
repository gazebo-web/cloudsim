package kubernetes

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/ingresses"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/ingresses/rules"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/nodes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/pods"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/services"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/spdy"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"k8s.io/client-go/kubernetes"
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
}

// IngressRules returns the Kubernetes orchestrator.IngressRules implementation.
func (k k8s) IngressRules() orchestrator.IngressRules {
	return k.ingressRules
}

// Nodes returns the Kubernetes orchestrator.Nodes implementation.
func (k k8s) Nodes() orchestrator.Nodes {
	return k.nodes
}

// Pods returns the Kubernetes orchestrator.Pods implementation.
func (k k8s) Pods() orchestrator.Pods {
	return k.pods
}

// Services returns the Kubernetes orchestrator.Services implementation.
func (k k8s) Services() orchestrator.Services {
	return k.services
}

// Ingresses returns the Kubernetes orchestrator.Ingresses implementation.
func (k k8s) Ingresses() orchestrator.Ingresses {
	return k.ingresses
}

// Config is used to group the inputs for NewCustomKubernetes.
// It includes all the needed subcomponents for Kubernetes.
type Config struct {
	Nodes        orchestrator.Nodes
	Pods         orchestrator.Pods
	Ingresses    orchestrator.Ingresses
	IngressRules orchestrator.IngressRules
	Services     orchestrator.Services
}

// NewCustomKubernetes returns a orchestrator.Cluster implementation using Kubernetes.
// All the subcomponents provided by the Config should be already initialized.
func NewCustomKubernetes(config Config) orchestrator.Cluster {
	return &k8s{
		nodes:        config.Nodes,
		pods:         config.Pods,
		ingresses:    config.Ingresses,
		ingressRules: config.IngressRules,
		services:     config.Services,
	}
}

// NewDefaultKubernetes initializes the set of Kubernetes subcomponents using
// the given kubernetes client api, spdy initializer and logger.
func NewDefaultKubernetes(api kubernetes.Interface, spdyInit spdy.Initializer, logger ign.Logger) orchestrator.Cluster {
	return &k8s{
		nodes:        nodes.NewNodes(api, logger),
		pods:         pods.NewPods(api, spdyInit, logger),
		ingressRules: rules.NewIngressRules(api, logger),
		services:     services.NewServices(api, logger),
		ingresses:    ingresses.NewIngresses(api, logger),
	}
}
