package orchestrator

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/configurations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/ingresses"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/network"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/nodes"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/pods"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/services"
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
