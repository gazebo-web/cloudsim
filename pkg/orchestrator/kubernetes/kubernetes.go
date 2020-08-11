package kubernetes

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"k8s.io/client-go/kubernetes"
)

type k8s struct {
	API         kubernetes.Interface
	nodeManager orchestrator.NodeManager
	podManager  orchestrator.PodManager
}

func (k k8s) Nodes() orchestrator.NodeManager {
	return k.nodeManager
}

func (k k8s) Pods() orchestrator.PodManager {
	return k.podManager
}

func (k k8s) Services() orchestrator.ServiceManager {
	panic("implement me")
}

func (k k8s) Ingresses() orchestrator.IngressManager {
	panic("implement me")
}

func NewKubernetes(nodeManager orchestrator.NodeManager, podManager orchestrator.PodManager) orchestrator.Orchestrator {
	return &k8s{
		nodeManager: nodeManager,
		podManager:  podManager,
	}
}
