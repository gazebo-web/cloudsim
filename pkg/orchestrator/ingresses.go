package orchestrator

// Ingresses groups a set of methods for managing Ingresses.
type Ingresses interface {
	Get(name string, namespace string) (Resource, error)
}
