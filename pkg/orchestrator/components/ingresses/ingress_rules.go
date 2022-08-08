package ingresses

import (
	"context"
	"errors"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/resource"
)

var (
	// ErrRuleNotFound is returned when a rule doesn't exist.
	ErrRuleNotFound = errors.New("rule not found")
	// ErrRuleEmpty is returned when a rule does not have any paths.
	ErrRuleEmpty = errors.New("rule has no paths")
)

// Rule is used to return a list of available paths to access a certain service.
type Rule interface {
	// Resource returns the resource that implements this rule.
	Resource() resource.Resource
	// Host returns the host where this rule is being applied to.
	Host() string
	// Paths returns the paths that this rule has.
	Paths() []Path
	// UpsertPaths insert and updates underlying list of paths with the given paths.
	UpsertPaths(paths []Path)
	// RemovePaths removes the given paths from the underlying list of paths.
	RemovePaths(paths []Path)
	// ToOutput converts this rule into a implementation-specific type of rule.
	// In Kubernetes, ToOutput is used to return a v1beta1.IngressRule.
	ToOutput() interface{}
}

// Path matches a certain Address to a specific Endpoint.
type Path struct {
	// UID is an unique identifier used to identify different paths.
	// In kubernetes: It's the URL path.
	// In gloo: It's the route name.
	UID string

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
	// Using Kubernetes: It's the name of the service we're pointing to
	// Using Gloo: It's the name of the upstream.
	Name string
	// Port is the port of the service.
	Port int32
}

// IngressRules groups a set of methods to manage rules from a certain Ingresses.
type IngressRules interface {
	Get(ctx context.Context, resource resource.Resource, host string) (Rule, error)
	Upsert(ctx context.Context, rule Rule, paths ...Path) error
	Remove(ctx context.Context, rule Rule, paths ...Path) error
}

// UpsertPaths updates or inserts the given elements into the given list.
// Returns the list after all elements have been updated/added.
func UpsertPaths(list, elements []Path) []Path {
	for _, p := range elements {
		var updated bool
		for i, rulePath := range list {
			if rulePath.UID == p.UID {
				updated = true
				list[i] = p
				break
			}
		}
		if !updated {
			list = append(list, p)
		}
	}
	return list
}

// RemovePaths removes the given elements from the given list.
// Returns the list after all elements have been removed.
func RemovePaths(list, elements []Path) []Path {
	for _, p := range elements {
		for i, rulePath := range list {
			if rulePath.Endpoint == p.Endpoint {
				pathsLen := len(list)
				if pathsLen > 1 {
					list[i] = list[pathsLen-1]
				}
				list = list[:pathsLen-1]
				break
			}
		}
	}
	return list
}
