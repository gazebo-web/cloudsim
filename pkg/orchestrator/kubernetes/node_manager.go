package kubernetes

import "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"

type nodeManager struct {
}

func (n nodeManager) Condition(node orchestrator.Resource, condition orchestrator.Condition) orchestrator.Waiter {
	panic("implement me")
}

func NewNodeManager() orchestrator.NodeManager {
	return &nodeManager{}
}
