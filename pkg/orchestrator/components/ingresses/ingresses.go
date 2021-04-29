package ingresses

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/resource"
)

// Ingresses groups a set of methods for managing Ingresses.
type Ingresses interface {
	Get(name string, namespace string) (resource.Resource, error)
}
