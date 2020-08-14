package orchestrator

// IngressManager groups a set of methods for managing Ingresses.
type IngressManager interface {
	Get(name string, namespace string) (Resource, error)
}
