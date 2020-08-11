package orchestrator

// ServiceManager groups a set of methods for managing Services.
// Services abstract a group of pods behind a single endpoint.
// TODO: Add CRUD operations.
type ServiceManager interface {
	Get(selector, name string)
}
