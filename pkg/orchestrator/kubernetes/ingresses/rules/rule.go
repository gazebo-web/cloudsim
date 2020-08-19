package rules

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// rule is an orchestrator.Rule implementation.
type rule struct {
	host     string
	paths    []orchestrator.Path
	resource orchestrator.Resource
}

func (r *rule) Resource() orchestrator.Resource {
	return r.resource
}

func (r *rule) UpsertPaths(paths []orchestrator.Path) {
	for _, p := range paths {
		var updated bool
		for i, rulePath := range r.paths {
			if rulePath.Endpoint == p.Endpoint {
				updated = true
				r.paths[i] = p
				break
			}
		}
		if !updated {
			r.paths = append(r.paths, p)
		}
	}
}

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
func (r *rule) Paths() []orchestrator.Path {
	return r.paths
}

// NewRule initializes a new orchestrator.Rule.
func NewRule(resource orchestrator.Resource, host string, paths []orchestrator.Path) orchestrator.Rule {
	return &rule{
		resource: resource,
		host:     host,
		paths:    paths,
	}
}

// NewPaths returns a generic group of Paths from the given set of Kubernetes HTTPIngressPaths.
func NewPaths(in []v1beta1.HTTPIngressPath) []orchestrator.Path {
	var out []orchestrator.Path
	for _, p := range in {
		out = append(out, orchestrator.Path{
			Address: p.Path,
			Endpoint: orchestrator.Endpoint{
				Name: p.Backend.ServiceName,
				Port: p.Backend.ServicePort.IntVal,
			},
		})
	}
	return out
}
