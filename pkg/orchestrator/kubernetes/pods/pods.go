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

// Get gets a pod with the certain name and in the given namespace and returns a resource that identifies that pod.
func (p *pods) Get(name, namespace string) (orchestrator.Resource, error) {
	p.Logger.Debug(fmt.Sprintf("Getting pod with name [%s] in namespace [%s]", name, namespace))

	pod, err := p.API.CoreV1().Pods(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		p.Logger.Debug(fmt.Sprintf(
			"Getting pod with name [%s] in namespace [%s] failed. Error: %+v.",
			name, namespace, err,
		))
		return nil, err
	}

	selector := orchestrator.NewSelector(pod.Labels)

	p.Logger.Debug(fmt.Sprintf(
		"Getting pod with name [%s] in namespace [%s] succeeded.",
		name, namespace,
	))

	return orchestrator.NewResource(name, namespace, selector), nil
}

// Create creates a new pod with the information given in orchestrator.CreatePodInput.
func (p *pods) Create(input orchestrator.CreatePodInput) (orchestrator.Resource, error) {
	p.Logger.Debug(fmt.Sprintf("Creating new pod. Input: %+v", input))

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
				SubPath:   v.SubPath,
			})
		}

		// Setup ports
		var ports []apiv1.ContainerPort
		for _, port := range c.Ports {
			ports = append(ports, apiv1.ContainerPort{ContainerPort: port})
		}

		// Setup env vars
		var envs []apiv1.EnvVar

		for key, from := range c.EnvVarsFrom {
			envs = append(envs, apiv1.EnvVar{
				Name:      key,
				ValueFrom: getEnvVarValueFromSource(from),
			})
		}

		for k, v := range c.EnvVars {
			envs = append(envs, apiv1.EnvVar{
				Name:  k,
				Value: v,
			})
		}

		// Add new container to list of containers
		containers = append(containers, apiv1.Container{
			Name:    c.Name,
			Image:   c.Image,
			Command: c.Command,
			Args:    c.Args,
			SecurityContext: &apiv1.SecurityContext{
				Privileged:               c.Privileged,
				AllowPrivilegeEscalation: c.AllowPrivilegeEscalation,
			},
			Ports:        ports,
			VolumeMounts: volumeMounts,
			Env:          envs,
		})
	}

	p.Logger.Debug(fmt.Sprintf("List of containers: %+v", containers))

	// Set up volumes
	var volumes []apiv1.Volume

	for _, v := range input.Volumes {
		hostPathType := apiv1.HostPathType(v.HostPathType)
		volumes = append(volumes, apiv1.Volume{
			Name: v.Name,
			VolumeSource: apiv1.VolumeSource{
				HostPath: &apiv1.HostPathVolumeSource{
					Path: v.HostPath,
					Type: &hostPathType,
				},
			},
		})
	}

	p.Logger.Debug(fmt.Sprintf("List of volumes: %+v", volumes))

	// Parse termination grace period config.
	terminationGracePeriod := int64(input.TerminationGracePeriodSeconds.Seconds())

	// Configure pod with previous settings
	pod := &apiv1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:   input.Name,
			Labels: input.Labels,
		},
		Spec: apiv1.PodSpec{
			RestartPolicy:                 apiv1.RestartPolicy(input.RestartPolicy),
			TerminationGracePeriodSeconds: &terminationGracePeriod,
			Containers:                    containers,
			Volumes:                       volumes,
			// These DNS servers provide alternative DNS server from the internet
			// in case the cluster DNS service isn't available
			DNSConfig: &apiv1.PodDNSConfig{
				Nameservers: input.Nameservers,
			},
		},
	}

	if input.NodeSelector != nil {
		pod.Spec.NodeSelector = input.NodeSelector.Map()
	}

	// Create pod in Kubernetes
	_, err := p.API.CoreV1().Pods(input.Namespace).Create(pod)
	if err != nil {
		p.Logger.Debug(fmt.Sprintf("Creating new pod failed. Input: %+v. Error: %s", input, err))
		return nil, err
	}

	// Create new resource
	res := orchestrator.NewResource(input.Name, input.Namespace, orchestrator.NewSelector(input.Labels))

	p.Logger.Debug(fmt.Sprintf("Creating new pod succeeded. Name: %s. Namespace: %s", res.Name(), res.Namespace()))
	return res, nil
}

// getEnvVarValueFromSource returns an env var source for the given value identified as from where it needs to get the env var.
func getEnvVarValueFromSource(from string) *apiv1.EnvVarSource {
	switch from {
	case orchestrator.EnvVarSourcePodIP:
		return &apiv1.EnvVarSource{
			FieldRef: &apiv1.ObjectFieldSelector{
				FieldPath: orchestrator.EnvVarSourcePodIP,
			},
		}
	}
	return nil
}

// Delete deletes the pod identified by the given resource.
func (p *pods) Delete(resource orchestrator.Resource) (orchestrator.Resource, error) {
	p.Logger.Debug(
		fmt.Sprintf("Deleting pod with name [%s] in namespace [%s]", resource.Name(), resource.Namespace()),
	)

	err := p.API.CoreV1().Pods(resource.Namespace()).Delete(resource.Name(), &metav1.DeleteOptions{})
	if err != nil {
		p.Logger.Debug(fmt.Sprintf(
			"Deleting pod with name [%s] in namespace [%s] failed. Error: %+v.",
			resource.Name(), resource.Namespace(), err,
		))
		return nil, err
	}

	p.Logger.Debug(fmt.Sprintf(
		"Deleting pod with name [%s] in namespace [%s] succeeded.",
		resource.Name(), resource.Namespace(),
	))

	return resource, nil
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

		// Get list of pods
		po, err := p.API.CoreV1().Pods(resource.Namespace()).List(opts)
		if err != nil {
			return false, err
		}

		if len(po.Items) == 0 {
			return false, orchestrator.ErrMissingPods
		}

		// Iterate over list of pods
		for _, i := range po.Items {
			var ready bool

			// Check that pod doesn't match the given condition.
			switch condition {
			case orchestrator.ReadyCondition:
				ready, err = p.isPodReady(&i)
				if err != nil {
					return false, err
				}
				break
			case orchestrator.HasIPStatusCondition:
				ready = p.podHasIP(&i)
				break
			}

			// Add pod to list if pod isn't ready.
			if !ready {
				pod := new(apiv1.Pod)
				*pod = i
				podsNotReady = append(podsNotReady, pod)
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

func (p *pods) podHasIP(pod *apiv1.Pod) bool {
	return pod.Status.PodIP != ""
}

// GetIP gets the IP for the pod identified with the given name in the current namespace.
// It will return an error if no IP has been assigned to the pod when calling this method.
// This job assumes that the pod is ready and can be accessed immediately. A WaitForCondition job must be executed at some point
// before executing this job to ensure that the pod is ready and has an IP assigned (orchestrator.HasIPStatusCondition).
func (p *pods) GetIP(name, namespace string) (string, error) {
	p.Logger.Debug(fmt.Sprintf("Getting IP from pod with name [%s] in namespace [%s]", name, namespace))

	pod, err := p.API.CoreV1().Pods(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		p.Logger.Debug(fmt.Sprintf(
			"Getting IP from pod with name [%s] in namespace [%s] failed. Error: %+v.",
			name, namespace, err,
		))
		return "", err
	}

	if !p.podHasIP(pod) {
		err = orchestrator.ErrPodHasNoIP

		p.Logger.Debug(fmt.Sprintf(
			"Getting IP from pod with name [%s] in namespace [%s] failed. Error: %+v.",
			name, namespace, err,
		))
		return "", err
	}

	p.Logger.Debug(fmt.Sprintf(
		"Getting IP from pod with name [%s] in namespace [%s] succeeded.",
		name, namespace,
	))

	return pod.Status.PodIP, nil
}

// NewPods initializes a new orchestrator.Pods implementation for managing Kubernetes Pods.
func NewPods(api kubernetes.Interface, spdy spdy.Initializer, logger ign.Logger) orchestrator.Pods {
	return &pods{
		API:    api,
		SPDY:   spdy,
		Logger: logger,
	}
}
