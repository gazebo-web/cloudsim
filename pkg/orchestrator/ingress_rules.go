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

// Path matches a certain Address to a specific Endpoint.
type Path struct {
	// Address is an extended POSIX regex as defined by IEEE Std 1003.1,
	// (i.e this follows the egrep/unix syntax, not the perl syntax)
	// matched against the path of an incoming request. Currently it can
	// contain characters disallowed from the conventional "path"
	// part of a URL as defined by RFC 3986. Paths must begin with
	// a '/'. If unspecified, the path defaults to a catch all sending
	// traffic to the backend.
	Address string

	// Endpoint has the information needed to route the incoming traffic in this Path
	// into a specific service running inside the cluster.
	Endpoint Endpoint
}

// Endpoint describes an entrypoint to a certain service name with the given port.
type Endpoint struct {
	// Name is the name of the service.
	Name string
	// Port is the port of the service.
	Port int32
}

// IngressRules groups a set of methods to manage rules from a certain Ingresses.
type IngressRules interface {
	Get(resource Resource, host string) (Rule, error)
	Upsert(rule Rule, paths ...Path) error
	Remove(host string, paths ...Path) error
}
