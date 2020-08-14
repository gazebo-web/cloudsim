package ingress

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// rule is an orchestrator.Rule implementation.
type rule struct {
	host  string
	paths []orchestrator.Path
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
			Path: p.Regex,
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
func NewRule(host string, paths []orchestrator.Path) orchestrator.Rule {
	return &rule{
		host:  host,
		paths: paths,
	}
}
