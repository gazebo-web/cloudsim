package kubernetes

import (
	gatewayV1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/client/clientset/versioned/typed/gateway.solo.io/v1"
	glooV1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/kube/client/clientset/versioned/typed/gloo.solo.io/v1"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/ingresses"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/ingresses/implementations/gloo"
	kubernetesIngresses "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/ingresses/implementations/kubernetes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/ingresses/implementations/kubernetes/rules"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/network"
	kubernetesNetwork "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/network/implementations/kubernetes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/nodes"
	kubernetesNodes "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/nodes/implementations/kubernetes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/pods"
	kubernetesPods "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/pods/implementations/kubernetes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/services"
	kubernetesServices "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/services/implementations/kubernetes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/spdy"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/implementations/kubernetes/client"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

// k8s is a orchestrator.Cluster implementation.
type k8s struct {
	// nodes has a reference to an nodes.Nodes implementation.
	nodes nodes.Nodes

	// pods has a reference to an pods.Pods implementation.
	pods pods.Pods

	// ingressRules has a reference to an ingresses.IngressRules implementation.
	ingressRules ingresses.IngressRules

	// services has a reference to an services.Services implementation.
	services services.Services

	// ingresses has a reference to an ingresses.Ingresses implementation.
	ingresses ingresses.Ingresses

	// networkPolicies has a reference to an network.Policies implementation.
	networkPolicies network.Policies
}

// IngressRules returns the Kubernetes ingresses.IngressRules implementation.
func (k *k8s) IngressRules() ingresses.IngressRules {
	return k.ingressRules
}

// Nodes returns the Kubernetes nodes.Nodes implementation.
func (k *k8s) Nodes() nodes.Nodes {
	return k.nodes
}

// Pods returns the Kubernetes pods.Pods implementation.
func (k *k8s) Pods() pods.Pods {
	return k.pods
}

// Services returns the Kubernetes services.Services implementation.
func (k *k8s) Services() services.Services {
	return k.services
}

// Ingresses returns the Kubernetes ingresses.Ingresses implementation.
func (k *k8s) Ingresses() ingresses.Ingresses {
	return k.ingresses
}

// NetworkPolicies returns the Kubernetes network.Policies implementation.
func (k *k8s) NetworkPolicies() network.Policies {
	return k.networkPolicies
}

// Config is used to group the inputs for NewCustomKubernetes.
// It includes all the needed subcomponents required by Kubernetes.
type Config struct {
	Nodes           nodes.Nodes
	Pods            pods.Pods
	Ingresses       ingresses.Ingresses
	IngressRules    ingresses.IngressRules
	Services        services.Services
	NetworkPolicies network.Policies
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
	}
}

// newDefaultKubernetes initializes the set of Kubernetes subcomponents using
// the given kubernetes client api, spdy initializer and logger.
func newDefaultKubernetes(api kubernetes.Interface, spdyInit spdy.Initializer, logger ign.Logger) orchestrator.Cluster {
	return &k8s{
		nodes:           kubernetesNodes.NewNodes(api, logger),
		pods:            kubernetesPods.NewPods(api, spdyInit, logger),
		ingressRules:    rules.NewIngressRules(api, logger),
		services:        kubernetesServices.NewServices(api, logger),
		ingresses:       kubernetesIngresses.NewIngresses(api, logger),
		networkPolicies: kubernetesNetwork.NewNetworkPolicies(api, logger),
	}
}

// newKubernetesWithGloo initializes the set of Kubernetes subcomponents using the given kubernetes client api,
// spdy initializer and logger. It also uses the given gloo and gateway clients to initialize Gloo to manage
// ingresses and ingress rules.
func newKubernetesWithGloo(api kubernetes.Interface, glooClient glooV1.GlooV1Interface, gatewayClient gatewayV1.GatewayV1Interface, spdyInit spdy.Initializer, logger ign.Logger) orchestrator.Cluster {
	return &k8s{
		nodes:           kubernetesNodes.NewNodes(api, logger),
		pods:            kubernetesPods.NewPods(api, spdyInit, logger),
		ingressRules:    gloo.NewVirtualHosts(gatewayClient, logger),
		services:        kubernetesServices.NewServices(api, logger),
		ingresses:       gloo.NewVirtualServices(gatewayClient, logger, glooClient),
		networkPolicies: kubernetesNetwork.NewNetworkPolicies(api, logger),
	}
}

// NewFakeKubernetes initializes the set of Kubernetes subcomponents using fake implementations.
func NewFakeKubernetes(logger ign.Logger) orchestrator.Cluster {
	api := fake.NewSimpleClientset()
	spdyInit := spdy.NewSPDYFakeInitializer()
	return &k8s{
		nodes:           kubernetesNodes.NewNodes(api, logger),
		pods:            kubernetesPods.NewPods(api, spdyInit, logger),
		ingressRules:    rules.NewIngressRules(api, logger),
		services:        kubernetesServices.NewServices(api, logger),
		ingresses:       kubernetesIngresses.NewIngresses(api, logger),
		networkPolicies: kubernetesNetwork.NewNetworkPolicies(api, logger),
	}
}

// InitializeKubernetes initializes a new Kubernetes orchestrator.
// `kubeconfig` is the path to the target cluster's kubeconfig file. If it is empty, the default config is used.
func InitializeKubernetes(kubeconfig string, logger ign.Logger) (orchestrator.Cluster, error) {
	config, err := client.GetConfig(kubeconfig)
	if err != nil {
		return nil, err
	}
	c, err := client.NewAPI(config)
	if err != nil {
		return nil, err
	}
	spdyInit := spdy.NewSPDYInitializer(config)
	return newDefaultKubernetes(c, spdyInit, logger), nil
}

// InitializeKubernetesWithGloo initializes a new Kubernetes orchestrator with Gloo to manage ingress and ingress rules.
// `kubeconfig` is the path to the target cluster's kubeconfig file. If it is empty, the default config is used.
func InitializeKubernetesWithGloo(kubeconfig string, logger ign.Logger) (orchestrator.Cluster, error) {
	config, err := client.GetConfig(kubeconfig)
	if err != nil {
		return nil, err
	}
	c, err := client.NewAPI(config)
	if err != nil {
		return nil, err
	}
	g, err := glooV1.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	gw, err := gatewayV1.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	spdyInit := spdy.NewSPDYInitializer(config)
	return newKubernetesWithGloo(c, g, gw, spdyInit, logger), nil
}
