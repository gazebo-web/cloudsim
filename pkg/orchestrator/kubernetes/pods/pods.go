package pods

import (
	"fmt"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/orchestrator/kubernetes/spdy"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/waiter"
	"gitlab.com/ignitionrobotics/web/ign-go"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	podutil "k8s.io/kubernetes/pkg/api/v1/pod"
	"k8s.io/kubernetes/pkg/client/conditions"
)

// pods is a orchestrator.Pods implementation.
type pods struct {
	API    kubernetes.Interface
	SPDY   spdy.Initializer
	Logger ign.Logger
}

func (p *pods) Create(input orchestrator.CreatePodInput) (orchestrator.Resource, error) {
	p.Logger.Debug(fmt.Sprintf("Creating new pod. Input: %+v", input))

	terminationGracePeriod := int64(input.TerminationGracePeriodSeconds.Seconds())

	// Set up containers for pod.
	var containers []apiv1.Container

	// Iterate over list of containers to create
	for _, c := range input.Containers {
		// Set up volume mounts.
		var volumeMounts []apiv1.VolumeMount
		for _, v := range c.Volumes {
			volumeMounts = append(volumeMounts, apiv1.VolumeMount{
				Name:      v.Name,
				MountPath: v.MountPath,
			})
		}

		// Setup ports
		var ports []apiv1.ContainerPort
		for _, port := range c.Ports {
			ports = append(ports, apiv1.ContainerPort{ContainerPort: port})
		}

		// Setup env vars
		var envs []apiv1.EnvVar
		for k, v := range c.EnvVars {
			envs = append(envs, apiv1.EnvVar{
				Name:  k,
				Value: v,
			})
		}

		// Add new container to list of containers
		containers = append(containers, apiv1.Container{
			Name:  c.Name,
			Image: c.Image,
			Args:  c.Args,
			SecurityContext: &apiv1.SecurityContext{
				Privileged:               c.Privileged,
				AllowPrivilegeEscalation: c.AllowPrivilegeEscalation,
			},
			Ports:        ports,
			VolumeMounts: volumeMounts,
			Env:          envs,
		})
	}

	// Set up volumes
	var volumes []apiv1.Volume
	for _, v := range input.Volumes {
		volumes = append(volumes, apiv1.Volume{
			Name: v.Name,
			VolumeSource: apiv1.VolumeSource{
				HostPath: &apiv1.HostPathVolumeSource{
					Path: v.HostPath,
				},
			},
		})
	}

	// Configure pod with previous settings
	pod := &apiv1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:   input.Name,
			Labels: input.Labels,
		},
		Spec: apiv1.PodSpec{
			RestartPolicy:                 apiv1.RestartPolicy(input.RestartPolicy),
			TerminationGracePeriodSeconds: &terminationGracePeriod,
			NodeSelector:                  input.NodeSelector.Map(),
			Containers:                    containers,
			Volumes:                       volumes,
			// These DNS servers provide alternative DNS server from the internet
			// in case the cluster DNS service isn't available
			DNSConfig: &apiv1.PodDNSConfig{
				Nameservers: input.Nameservers,
			},
		},
	}

	// Create pod in Kubernetes
	_, err := p.API.CoreV1().Pods(input.Namespace).Create(pod)
	if err != nil {
		p.Logger.Debug(fmt.Sprintf("Creating new failed. Input: %+v. Error: %s", input, err))
		return nil, err
	}

	// Create new resource
	res := orchestrator.NewResource(input.Name, input.Namespace, orchestrator.NewSelector(input.Labels))

	p.Logger.Debug(fmt.Sprintf("Creating new pod succeeded. Name: %s. Namespace: %s", res.Name(), res.Namespace()))
	return res, nil
}

// Exec creates a new executor.
func (p *pods) Exec(pod orchestrator.Resource) orchestrator.Executor {
	p.Logger.Debug(fmt.Sprintf("Creating new executor for pod [%s]", pod.Name()))
	return newExecutor(p.API, pod, p.SPDY, p.Logger)
}

// Reader creates a new reader.
func (p *pods) Reader(pod orchestrator.Resource) orchestrator.Reader {
	p.Logger.Debug(fmt.Sprintf("Creating new reader for pod [%s]", pod.Name()))
	return newReader(p.API, pod, p.SPDY, p.Logger)
}

// WaitForCondition creates a new wait request that will be used to wait for a resource to match a certain condition.
// The wait request won't be triggered until the method Wait has been called.
func (p *pods) WaitForCondition(resource orchestrator.Resource, condition orchestrator.Condition) waiter.Waiter {
	p.Logger.Debug(fmt.Sprintf("Creating wait for condition [%+v] request on pods matching the following selector: [%s]",
		condition, resource.Selector(),
	))

	// Prepare options
	opts := metav1.ListOptions{
		LabelSelector: resource.Selector().String(),
	}

	// Create job
	job := func() (bool, error) {
		var podsNotReady []*apiv1.Pod
		po, err := p.API.CoreV1().Pods(resource.Namespace()).List(opts)
		if err != nil {
			return false, err
		}
		for _, i := range po.Items {
			if condition == orchestrator.ReadyCondition {
				ready, err := p.isPodReady(&i)
				if err != nil {
					return false, err
				}
				if !ready {
					pod := new(apiv1.Pod)
					*pod = i
					podsNotReady = append(podsNotReady, pod)
				}
			}
		}
		return len(podsNotReady) == 0, nil
	}

	p.Logger.Debug(fmt.Sprintf(
		"Wait for condition [%+v] request on pods matching the following selector: [%s] was created.",
		condition, resource.Selector(),
	))

	return waiter.NewWaitRequest(job)
}

// isPodReady checks if the given Kubernetes Pod matches the ready condition.
func (p *pods) isPodReady(pod *apiv1.Pod) (bool, error) {
	if pod.Status.Phase == apiv1.PodFailed || pod.Status.Phase == apiv1.PodSucceeded {
		return false, conditions.ErrPodCompleted
	}
	return podutil.IsPodReady(pod), nil
}

// NewPods initializes a new orchestrator.Pods implementation for managing Kubernetes Pods.
func NewPods(api kubernetes.Interface, spdy spdy.Initializer, logger ign.Logger) orchestrator.Pods {
	return &pods{
		API:    api,
		SPDY:   spdy,
		Logger: logger,
	}
}
