package orchestrator

// Rule is used to return a list of paths.
type Rule interface {
	Host() string
	Paths() []Path
}

// Path matches a certain Regex to a specific Endpoint.
type Path struct {
	Regex    string
	Endpoint Endpoint
}

// Endpoint describes an entrypoint to a certain service name with the given port.
type Endpoint struct {
	// Name is the name of the service.
	Name string
	// Port is the port of the service.
	Port int32
}

// IngressManager groups a set of methods for managing Ingresses.
type IngressManager interface {
	Rules(ingress Resource) Ruler
}

// Ruler groups a set of methods to interact with an Ingress's rules.
type Ruler interface {
	Get(host string) (Rule, error)
	Upsert(host string, paths ...Path) error
	Remove(host string, paths ...Path) error
}
