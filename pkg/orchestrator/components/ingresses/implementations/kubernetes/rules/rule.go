package rules

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/components/ingresses"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/resource"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// rule is an ingresses.Rule implementation.
type rule struct {
	host     string
	paths    []ingresses.Path
	resource resource.Resource
}

// Resource returns the resource associated with the current rule.
func (r *rule) Resource() resource.Resource {
	return r.resource
}

// UpsertPaths inserts and update the given paths into the current rule.
func (r *rule) UpsertPaths(paths []ingresses.Path) {
	r.paths = ingresses.UpsertPaths(r.paths, paths)
}

// RemovePaths removes the paths from the current rule.
func (r *rule) RemovePaths(paths []ingresses.Path) {
	r.paths = ingresses.RemovePaths(r.paths, paths)
}

// toIngressPaths converts the current rule paths into an slice of v1beta1.HTTPIngressPath.
func (r *rule) toIngressPaths() []v1beta1.HTTPIngressPath {
	var result []v1beta1.HTTPIngressPath
	for _, p := range r.paths {
		result = append(result, v1beta1.HTTPIngressPath{
			Path: p.Address,
			Backend: v1beta1.IngressBackend{
				ServiceName: p.Endpoint.Name,
				ServicePort: intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: p.Endpoint.Port,
				},
			},
		})
	}
	return result
}

// ToOutput returns the current rule as a v1beta1.IngressRule.
func (r *rule) ToOutput() interface{} {
	return v1beta1.IngressRule{
		Host: r.host,
		IngressRuleValue: v1beta1.IngressRuleValue{
			HTTP: &v1beta1.HTTPIngressRuleValue{
				Paths: r.toIngressPaths(),
			},
		},
	}
}

// Host returns the rule's host.
func (r *rule) Host() string {
	return r.host
}

// Paths returns an array of paths.
func (r *rule) Paths() []ingresses.Path {
	return r.paths
}

// NewRule initializes a new ingresses.Rule.
func NewRule(resource resource.Resource, host string, paths []ingresses.Path) ingresses.Rule {
	return &rule{
		resource: resource,
		host:     host,
		paths:    paths,
	}
}

// NewPaths returns a generic group of Paths from the given set of Kubernetes HTTPIngressPaths.
func NewPaths(in []v1beta1.HTTPIngressPath) []ingresses.Path {
	var out []ingresses.Path
	for _, p := range in {
		out = append(out, ingresses.Path{
			UID:     p.Backend.ServiceName,
			Address: p.Path,
			Endpoint: ingresses.Endpoint{
				Name: p.Backend.ServiceName,
				Port: p.Backend.ServicePort.IntVal,
			},
		})
	}
	return out
}
