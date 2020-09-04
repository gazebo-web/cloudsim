package env

import "gitlab.com/ignitionrobotics/web/cloudsim/pkg/store"

type orchestrator struct {
}

func (o orchestrator) Namespace() string {
	panic("implement me")
}

func newOrchestratorStore() store.Orchestrator {
	return &orchestrator{}
}
