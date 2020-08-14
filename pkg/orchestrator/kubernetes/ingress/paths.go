package ingress

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"k8s.io/api/extensions/v1beta1"
)

// NewPaths returns a generic group of Paths from the given set of Kubernetes HTTPIngressPaths.
func NewPaths(in []v1beta1.HTTPIngressPath) []orchestrator.Path {
	var out []orchestrator.Path
	for _, p := range in {
		out = append(out, orchestrator.Path{
			Regex: p.Path,
			Endpoint: orchestrator.Endpoint{
				Name: p.Backend.ServiceName,
				Port: p.Backend.ServicePort.IntVal,
			},
		})
	}
	return out
}
