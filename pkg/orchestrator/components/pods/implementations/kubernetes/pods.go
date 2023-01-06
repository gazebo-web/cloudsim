package kubernetes

import (
	"context"
	"fmt"
	"github.com/gazebo-web/cloudsim/v4/pkg/orchestrator/components/pods"
	"github.com/gazebo-web/cloudsim/v4/pkg/orchestrator/components/spdy"
	orchestratorResource "github.com/gazebo-web/cloudsim/v4/pkg/orchestrator/resource"
	"github.com/gazebo-web/cloudsim/v4/pkg/waiter"
	"github.com/gazebo-web/gz-go/v7"
	"github.com/gazebo-web/gz-go/v7/kubernetes"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	client "k8s.io/client-go/kubernetes"
)

// kubernetesPods is a pods.Pods implementation.
type kubernetesPods struct {
	API    client.Interface
	SPDY   spdy.Initializer
	Logger gz.Logger
}

// List returns a list of pod resources matching the giving selector in the given namespace.
// If selector is nil or empty (doesn't have any labels specified) it will return all the resources in the given namespace.
func (p *kubernetesPods) List(ctx context.Context, namespace string, selector orchestratorResource.Selector) ([]pods.PodResource, error) {

	if selector == nil {
		selector = orchestratorResource.NewSelector(map[string]string{})
	}
	p.Logger.Debug(fmt.Sprintf("Getting list of pods in namespace [%s] matching the following labels: [%s]", namespace, selector.String()))
	res, err := p.API.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{LabelSelector: selector.String()})
	if err != nil {
		p.Logger.Debug(fmt.Sprintf("Failed to list pods in namespace [%s] matching the following labels: [%s]", namespace, selector.String()))
		return nil, err
	}

	if len(res.Items) == 0 {
		p.Logger.Debug(fmt.Sprintf("No pods available in namespace [%s] matching the following labels: [%s]", namespace, selector.String()))
		return nil, nil
	}

	list := make([]pods.PodResource, len(res.Items))

	for i, po := range res.Items {
		list[i] = kubernetesPodToPodResource(po)
	}

	p.Logger.Debug(fmt.Sprintf("Getting list of pods in namespace [%s] matching the following labels: [%s] succeeded.", namespace, selector.String()))
	return list, nil
}

// Get gets a pod with the certain name and in the given namespace and returns a resource that identifies that pod.
func (p *kubernetesPods) Get(ctx context.Context, name, namespace string) (*pods.PodResource, error) {
	p.Logger.Debug(fmt.Sprintf("Getting pod with name [%s] in namespace [%s]", name, namespace))

	pod, err := p.API.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		p.Logger.Debug(fmt.Sprintf("Getting pod with name [%s] in namespace [%s] failed. Error: %+v.", name, namespace, err))
		return nil, err
	}

	p.Logger.Debug(fmt.Sprintf("Getting pod with name [%s] in namespace [%s] succeeded.", name, namespace))

	res := kubernetesPodToPodResource(*pod)
	return &res, nil
}

// generateKubernetesContainers takes a generic set of cloudsim containers and generate their counterpart for Kubernetes.
func generateKubernetesContainers(containers []pods.Container) []apiv1.Container {
	var result []apiv1.Container

	for _, c := range containers {
		var volumeMounts []apiv1.VolumeMount
		for _, v := range c.Volumes {
			volumeMounts = append(volumeMounts, ParseVolumeMount(v))
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

		var resourceRequests map[apiv1.ResourceName]resource.Quantity
		if len(c.ResourceRequests) > 0 {
			resourceRequests = make(map[apiv1.ResourceName]resource.Quantity, len(c.ResourceRequests))
			for k, v := range c.ResourceRequests {
				resourceRequests[apiv1.ResourceName(k)] = resource.MustParse(v)
			}
		}

		var resourceLimits map[apiv1.ResourceName]resource.Quantity
		if len(c.ResourceLimits) > 0 {
			resourceLimits = make(map[apiv1.ResourceName]resource.Quantity, len(c.ResourceLimits))
			for k, v := range c.ResourceLimits {
				resourceLimits[apiv1.ResourceName(k)] = resource.MustParse(v)
			}
		}

		// Add new container to list of containers
		result = append(result, apiv1.Container{
			Name:    c.Name,
			Image:   c.Image,
			Command: c.Command,
			Args:    c.Args,
			Ports:   ports,
			Env:     envs,
			Resources: apiv1.ResourceRequirements{
				Requests: resourceRequests,
				Limits:   resourceLimits,
			},
			VolumeMounts: volumeMounts,
			SecurityContext: &apiv1.SecurityContext{
				Privileged:               c.Privileged,
				AllowPrivilegeEscalation: c.AllowPrivilegeEscalation,
			},
		})
	}

	return result
}

// Create creates a new pod with the information given in resource.CreatePodInput.
func (p *kubernetesPods) Create(ctx context.Context, input pods.CreatePodInput) (*pods.PodResource, error) {
	p.Logger.Debug(fmt.Sprintf("Creating new pod. Input: %+v", input))

	// Set up init containers
	initContainers := generateKubernetesContainers(input.InitContainers)

	// Set up containers for pod
	containers := generateKubernetesContainers(input.Containers)

	p.Logger.Debug(fmt.Sprintf("List of containers: %+v", containers))

	// Set up volumes
	var volumes []apiv1.Volume
	for _, v := range input.Volumes {
		volumes = append(volumes, ParseVolume(v))
	}

	p.Logger.Debug(fmt.Sprintf("List of volumes: %+v", volumes))

	// Parse image pull secrets
	imagePullSecrets := make([]apiv1.LocalObjectReference, len(input.ImagePullCredentials))
	for i, secret := range input.ImagePullCredentials {
		imagePullSecrets[i].Name = secret
	}

	// Parse termination grace period config.
	terminationGracePeriod := int64(input.TerminationGracePeriodSeconds.Seconds())

	// Configure pod with previous settings
	pod := &apiv1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:   input.Name,
			Labels: input.Labels,
		},
		Spec: apiv1.PodSpec{
			ImagePullSecrets:              imagePullSecrets,
			RestartPolicy:                 apiv1.RestartPolicy(input.RestartPolicy),
			TerminationGracePeriodSeconds: &terminationGracePeriod,
			InitContainers:                initContainers,
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
	created, err := p.API.CoreV1().Pods(input.Namespace).Create(ctx, pod, metav1.CreateOptions{})
	if err != nil {
		p.Logger.Debug(fmt.Sprintf("Creating new pod failed. Input: %+v. Error: %s", input, err))
		return nil, err
	}

	// Create new resource
	res := kubernetesPodToPodResource(*created)

	p.Logger.Debug(fmt.Sprintf("Creating new pod succeeded. Name: %s. Namespace: %s", res.Name(), res.Namespace()))
	return &res, nil
}

// getEnvVarValueFromSource returns an env var source for the given value identified as from where it needs to get the env var.
func getEnvVarValueFromSource(from string) *apiv1.EnvVarSource {
	switch from {
	case pods.EnvVarSourcePodIP:
		return &apiv1.EnvVarSource{
			FieldRef: &apiv1.ObjectFieldSelector{
				FieldPath: pods.EnvVarSourcePodIP,
			},
		}
	}
	return nil
}

// Delete deletes the pod identified by the given resource.
func (p *kubernetesPods) Delete(ctx context.Context, resource orchestratorResource.Resource) (orchestratorResource.Resource, error) {
	p.Logger.Debug(
		fmt.Sprintf("Deleting pod with name [%s] in namespace [%s]", resource.Name(), resource.Namespace()),
	)

	err := p.API.CoreV1().Pods(resource.Namespace()).Delete(ctx, resource.Name(), metav1.DeleteOptions{})
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
func (p *kubernetesPods) Exec(ctx context.Context, pod orchestratorResource.Resource) pods.Executor {
	p.Logger.Debug(fmt.Sprintf("Creating new executor for pod [%s] in namespace [%s]", pod.Name(), pod.Namespace()))
	return newExecutor(p.API, pod, p.SPDY, p.Logger)
}

// Reader creates a new reader.
func (p *kubernetesPods) Reader(ctx context.Context, pod orchestratorResource.Resource) pods.Reader {
	p.Logger.Debug(fmt.Sprintf("Creating new reader for pod [%s]", pod.Name()))
	return newReader(p.API, pod, p.SPDY, p.Logger)
}

// WaitForCondition creates a new wait request that will be used to wait for a resource to match a certain condition.
// The wait request won't be triggered until the method Wait has been called.
// If more than one condition is provided, the waiter will return once any of the conditions are met.
func (p *kubernetesPods) WaitForCondition(ctx context.Context, resource orchestratorResource.Resource,
	conditions ...orchestratorResource.Condition) waiter.Waiter {

	p.Logger.Debug(fmt.Sprintf(
		"Creating wait for conditions [%+v] request on pods matching the following selector: [%s]",
		conditions,
		resource.Selector(),
	))

	// Prepare options
	opts := metav1.ListOptions{
		LabelSelector: resource.Selector().String(),
	}

	// Create job
	job := func() (bool, error) {
		var podsNotReady []*apiv1.Pod

		// Get list of pods
		po, err := p.API.CoreV1().Pods(resource.Namespace()).List(ctx, opts)
		if err != nil {
			p.Logger.Debug("[WaitForCondition] Failed to get pods from orchestrator: ", err)
			return false, nil
		}

		if len(po.Items) == 0 {
			return false, pods.ErrMissingPods
		}

		// Check that pod matches any of the given conditions
		for _, item := range po.Items {
			var ready bool
			for _, condition := range conditions {
				switch condition {
				case orchestratorResource.ReadyCondition:
					ready, err = kubernetes.IsPodReady(&item)
					if err != nil && err != kubernetes.ErrPodCompleted {
						return false, nil
					}
				case orchestratorResource.SucceededCondition:
					ready = kubernetes.IsPodSucceeded(&item)
				case orchestratorResource.FailedCondition:
					ready = kubernetes.IsPodFailed(&item)
				case orchestratorResource.HasIPStatusCondition:
					ready = p.podHasIP(&item)
				}
				if ready {
					break
				}
			}

			// Add pod to list if pod isn't ready.
			if !ready {
				pod := new(apiv1.Pod)
				*pod = item
				podsNotReady = append(podsNotReady, pod)
			}
		}
		return len(podsNotReady) == 0, nil
	}

	p.Logger.Debug(fmt.Sprintf(
		"Wait for conditions [%+v] request on pods matching the following selector: [%s] was created.",
		conditions, resource.Selector(),
	))

	return waiter.NewWaitRequest(job)
}

func (p *kubernetesPods) podHasIP(pod *apiv1.Pod) bool {
	return pod.Status.PodIP != ""
}

// GetIP gets the IP for the pod identified with the given name in the current namespace.
// It will return an error if no IP has been assigned to the pod when calling this method.
// This job assumes that the pod is ready and can be accessed immediately. A WaitForCondition job must be executed at some point
// before executing this job to ensure that the pod is ready and has an IP assigned (resource.HasIPStatusCondition).
func (p *kubernetesPods) GetIP(ctx context.Context, name string, namespace string) (string, error) {
	p.Logger.Debug(fmt.Sprintf("Getting IP from pod with name [%s] in namespace [%s]", name, namespace))

	pod, err := p.API.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		p.Logger.Debug(fmt.Sprintf(
			"Getting IP from pod with name [%s] in namespace [%s] failed. Error: %+v.",
			name, namespace, err,
		))
		return "", err
	}

	if !p.podHasIP(pod) {
		err = pods.ErrPodHasNoIP

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

// NewPods initializes a new pods.Pods implementation for managing Kubernetes Pods.
func NewPods(api client.Interface, spdy spdy.Initializer, logger gz.Logger) pods.Pods {
	return &kubernetesPods{
		API:    api,
		SPDY:   spdy,
		Logger: logger,
	}
}
