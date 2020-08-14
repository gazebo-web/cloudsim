package orchestrator

import "errors"

var (
	// ErrRuleNotFound is returned when a rule doesn't exist.
	ErrRuleNotFound = errors.New("rule not found")
)

// Rule is used to return a list of available paths to access a certain service.
type Rule interface {
	Resource() Resource
	Host() string
	Paths() []Path
	UpsertPaths(paths []Path)
	ToOutput() interface{}
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

type IngressRulesManager interface {
	Get(resource Resource, host string) (Rule, error)
	Upsert(rule Rule, paths ...Path) error
	Remove(host string, paths ...Path) error
}
