package kubernetes

import "gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"

type podManager struct {
}

func (p podManager) Exec(pod orchestrator.Resource) orchestrator.Executor {
	panic("implement me")
}

func (p podManager) Reader(pod orchestrator.Resource) orchestrator.Reader {
	panic("implement me")
}

func (p podManager) Condition(pod orchestrator.Resource, condition orchestrator.Condition) orchestrator.Waiter {
	panic("implement me")
}

func NewPodManager() orchestrator.PodManager {
	return &podManager{}
}
