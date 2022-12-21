package ingresses

import (
	"context"
	"github.com/gazebo-web/cloudsim/pkg/orchestrator/resource"
)

// Ingresses groups a set of methods for managing Ingresses.
type Ingresses interface {
	Get(ctx context.Context, name string, namespace string) (resource.Resource, error)
}
