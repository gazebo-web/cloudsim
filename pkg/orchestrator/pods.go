package orchestrator

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/waiter"
	"io"
	"time"
)

// RestartPolicy defines a restart policy used for pods.
type RestartPolicy string

var (
	// RestartPolicyNever is used to indicate that a pod won't be restarted.
	RestartPolicyNever RestartPolicy = "Never"

	// RestartPolicyAlways is used to indicate that a pod always will be restarted.
	RestartPolicyAlways RestartPolicy = "Always"

	// RestartPolicyOnFailure is used to indicate that a pod will be restarted only on failures.
	RestartPolicyOnFailure RestartPolicy = "OnFailure"
)

// Volume represents a storage that will be used to persist data from a certain Container.
type Volume struct {
	// Name is the name of the volume.
	Name string
	// HostPath is the mounting path.
	HostPath string

	MountPath string
}

// Container is a represents of a standard unit of software.
type Container struct {
	// Name is the container's name.
	Name string

	// Image is the image running inside the container.
	Image string

	// Args are the commands passed to the container.
	Args []string

	// Privileged defines if the container should run in privileged mode.
	Privileged *bool

	// AllowPrivilegeEscalation is used to define if the container is allowed to scale its privileges.
	AllowPrivilegeEscalation *bool

	// Ports is the list of ports that should be opened.
	Ports []int32

	// Volumes is the list of volumes that should be mounted in the container.
	Volumes []Volume

	// EnvVars is the list of env vars that should be passed into the container.
	EnvVars map[string]string
}

// CreatePodInput is the input of Pods.Create method.
type CreatePodInput struct {
	// Name is the name of the pod that will be created.
	Name string
	// Namespace is the namespace where the pod will live in.
	Namespace string
	// Labels are the map of labels that will be applied to the pod.
	Labels map[string]string
	// RestartPolicy defines how the pod should react after an error.
	RestartPolicy RestartPolicy
	// TerminationGracePeriodSeconds is the time duration in seconds the pod needs to terminate gracefully.
	TerminationGracePeriodSeconds time.Duration
	// NodeSelector defines the node where the pod should run in.
	NodeSelector Selector
	// Containers is the list of containers that should be created inside the pod.
	Containers []Container
	// Volumes are the list of volumes that will be created to persist the data from the Container.Volumes.
	Volumes []Volume
	// Nameservers are the list of DNS Nameservers that will be used to expose the pod to the internet.
	Nameservers []string
}

// Pods groups a set of methods to perform an operation with a Pod.
type Pods interface {
	Create(input CreatePodInput) (Resource, error)
	Exec(resource Resource) Executor
	Reader(resource Resource) Reader
	WaitForCondition(resource Resource, condition Condition) waiter.Waiter
}

// Executor groups a set of methods to run commands and scripts inside a Pod.
type Executor interface {
	// Cmd runs a command inside a container.
	Cmd(command []string) error
	// Script runs a script inside a container.
	// Could be used to run copy_to_s3.sh
	Script(path string) error
}

// Reader groups a set of methods to read files and logs from a Pod.
type Reader interface {
	File(paths ...string) (io.Reader, error)
	Logs(container string, lines int64) (string, error)
}
