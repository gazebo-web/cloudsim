package orchestrator

// IngressManager groups a set of methods for managing Ingresses.
type IngressManager interface {
	GetByName(name string)
	Update(ingress Resource, data interface{})
	Rules(ingress Resource) Ruler
}

// Ruler groups a set of methods to interact with an Ingress's rules.
type Ruler interface {
	Get(host string)
	Upsert(host string, paths ...string)
	Remove(host string, paths ...string)
}
