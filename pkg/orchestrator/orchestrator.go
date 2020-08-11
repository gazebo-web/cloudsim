package orchestrator

import (
	"time"
)

// Condition represents a state that should be reached.
type Condition string

var (
	// ReadyCondition is used to indicate that Nodes and Pods are ready.
	ReadyCondition Condition = "Ready"
)

// Orchestrator groups a set of methods for managing a cluster.
type Orchestrator interface {
	Nodes() NodeManager
	Pods() PodManager
	Services() ServiceManager
	Ingresses() IngressManager
}

// Resource groups a set of method to identify a resource in a cluster.
type Resource interface {
	// Name returns the name of the resource
	Name() string
	// Selector returns a set of key-value pairs that identify the resource.
	Selector() map[string]string
	// Namespace returns the namespace where the resource lives in.
	Namespace() string
}

// Waiter groups a set of methods to wait.
type Waiter interface {
	Wait(job func(), timeout time.Duration, pollFrequency time.Duration)
}
