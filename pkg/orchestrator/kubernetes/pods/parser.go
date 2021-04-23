package pods

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	apiv1 "k8s.io/api/core/v1"
	"time"
)

// parseKubernetesPod parses the given Kubernetes apiv1.Pod into a orchestrator.PodResource.
func parseKubernetesPod(pod apiv1.Pod) orchestrator.PodResource {
	var deletion *time.Time
	if pod.DeletionTimestamp != nil {
		deletion = new(time.Time)
		*deletion = pod.DeletionTimestamp.Time
	}
	return orchestrator.PodResource{
		Resource:          orchestrator.NewResource(pod.Name, pod.Namespace, orchestrator.NewSelector(pod.Labels)),
		ResourcePhase:     orchestrator.NewResourcePhase(orchestrator.Phase(pod.Status.Phase)),
		ResourceTimestamp: orchestrator.NewResourceTimestamp(pod.CreationTimestamp.Time, deletion),
	}
}
