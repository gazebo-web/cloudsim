package kubernetes

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/configurations"
	kubernetesConfigMaps "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/configurations/implementations/kubernetes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/ingresses"
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
	"gitlab.com/ignitionrobotics/web/ign-go/v6"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

// k8s is a orchestrator.Cluster implementation.
type k8s struct {
	// nodes has a reference to a nodes.Nodes implementation.
	nodes nodes.Nodes

	// pods has a reference to a pods.Pods implementation.
	pods pods.Pods

	// ingressRules has a reference to an ingresses.IngressRules implementation.
	ingressRules ingresses.IngressRules

	// ingresses has a reference to an ingresses.Ingresses implementation.
	ingresses ingresses.Ingresses

	// services has a reference to a services.Services implementation.
	services services.Services

	// networkPolicies has a reference to a network.Policies implementation.
	networkPolicies network.Policies

	// configurations has a reference to a configurations.Configurations implementation.
	configurations configurations.Configurations
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

// Configurations returns the Kubernetes configurations.Configurations implementation.
func (k *k8s) Configurations() configurations.Configurations {
	return k.configurations
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
	Configurations  configurations.Configurations
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
		configurations:  config.Configurations,
	}
}

// NewDefaultKubernetes initializes the set of Kubernetes subcomponents using
// the given kubernetes client api, spdy initializer and logger.
func NewDefaultKubernetes(api kubernetes.Interface, spdyInit spdy.Initializer, logger ign.Logger) orchestrator.Cluster {
	return &k8s{
		nodes:           kubernetesNodes.NewNodes(api, logger),
		pods:            kubernetesPods.NewPods(api, spdyInit, logger),
		ingressRules:    rules.NewIngressRules(api, logger),
		services:        kubernetesServices.NewServices(api, logger),
		ingresses:       kubernetesIngresses.NewIngresses(api, logger),
		networkPolicies: kubernetesNetwork.NewNetworkPolicies(api, logger),
		configurations:  kubernetesConfigMaps.NewConfigMaps(api, logger),
	}
}

// NewFakeKubernetes initializes the set of Kubernetes subcomponents using fake implementations.
func NewFakeKubernetes(logger ign.Logger) (orchestrator.Cluster, *fake.Clientset) {
	api := fake.NewSimpleClientset()
	spdyInit := spdy.NewSPDYFakeInitializer()
	return &k8s{
		nodes:           kubernetesNodes.NewNodes(api, logger),
		pods:            kubernetesPods.NewPods(api, spdyInit, logger),
		ingressRules:    rules.NewIngressRules(api, logger),
		services:        kubernetesServices.NewServices(api, logger),
		ingresses:       kubernetesIngresses.NewIngresses(api, logger),
		networkPolicies: kubernetesNetwork.NewNetworkPolicies(api, logger),
		configurations:  kubernetesConfigMaps.NewConfigMaps(api, logger),
	}, api
}

// InitializeKubernetes initializes a new Kubernetes orchestrator.
// `kubeconfig` is the path to the target cluster's kubeconfig file. If it is empty, the default config is used.
func InitializeKubernetes(kubeconfig string, logger ign.Logger) (orchestrator.Cluster, error) {
	config, err := client.GetConfig(kubeconfig)
	if err != nil {
		return nil, err
	}
	client, err := client.NewAPI(config)
	if err != nil {
		return nil, err
	}
	spdyInit := spdy.NewSPDYInitializer(config)
	return NewDefaultKubernetes(client, spdyInit, logger), nil
}
