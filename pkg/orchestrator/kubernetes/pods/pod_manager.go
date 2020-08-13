package pods

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/spdy"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/waiter"
	"k8s.io/client-go/kubernetes"
)

// manager is a orchestrator.PodManager implementation.
type manager struct {
	API  kubernetes.Interface
	SPDY spdy.Initializer
}

// Exec creates a new executor.
func (p manager) Exec(pod orchestrator.Resource) orchestrator.Executor {
	return newExecutor(p.API, pod, p.SPDY)
}

// Reader creates a new reader.
func (p manager) Reader(pod orchestrator.Resource) orchestrator.Reader {
	return newReader(p.API, pod, p.SPDY)
}

// Condition creates a new wait request.
func (p manager) Condition(pod orchestrator.Resource, condition orchestrator.Condition) waiter.Waiter {
	panic("implement me")
}

// NewManager initializes a new manager.
func NewManager(api kubernetes.Interface, spdy spdy.Initializer) orchestrator.PodManager {
	return &manager{
		API:  api,
		SPDY: spdy,
	}
}
