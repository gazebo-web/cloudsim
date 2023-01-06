package orchestrator

import (
	"github.com/gazebo-web/cloudsim/pkg/orchestrator/components/configurations"
	"github.com/gazebo-web/cloudsim/pkg/orchestrator/components/ingresses"
	"github.com/gazebo-web/cloudsim/pkg/orchestrator/components/network"
	"github.com/gazebo-web/cloudsim/pkg/orchestrator/components/nodes"
	"github.com/gazebo-web/cloudsim/pkg/orchestrator/components/pods"
	"github.com/gazebo-web/cloudsim/pkg/orchestrator/components/services"
)

// Cluster groups a set of methods for managing a cluster.
type Cluster interface {
	Nodes() nodes.Nodes
	Pods() pods.Pods
	Services() services.Services
	Ingresses() ingresses.Ingresses
	IngressRules() ingresses.IngressRules
	NetworkPolicies() network.Policies
	Configurations() configurations.Configurations
}
