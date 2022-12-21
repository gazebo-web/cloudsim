package kubernetes

import (
	"github.com/gazebo-web/cloudsim/pkg/orchestrator/components/pods"
	"github.com/gazebo-web/cloudsim/pkg/orchestrator/resource"
	"github.com/gazebo-web/cloudsim/pkg/orchestrator/resource/phase"
	"github.com/gazebo-web/cloudsim/pkg/orchestrator/resource/timestamp"
	apiv1 "k8s.io/api/core/v1"
	"time"
)

// kubernetesPodToPodResource converts the given Kubernetes apiv1.Pod into an orchestrator.PodResource.
func kubernetesPodToPodResource(pod apiv1.Pod) pods.PodResource {
	var deletion *time.Time
	if pod.DeletionTimestamp != nil {
		deletion = &pod.DeletionTimestamp.Time
	}
	return pods.PodResource{
		Resource:          resource.NewResource(pod.Name, pod.Namespace, resource.NewSelector(pod.Labels)),
		ResourcePhase:     phase.NewResourcePhase(phase.Phase(pod.Status.Phase)),
		ResourceTimestamp: timestamp.NewResourceTimestamp(pod.CreationTimestamp.Time, deletion),
	}
}
