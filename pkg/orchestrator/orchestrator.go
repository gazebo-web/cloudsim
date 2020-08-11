package orchestrator

import (
	"time"
)

// Condition represents a state that should be reached.
type Condition struct {
	Type   string
	Status string
}

var (
	// ReadyCondition is used to indicate that Nodes and Pods are ready.
	ReadyCondition = Condition{
		Type:   "Ready",
		Status: "True",
	}
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
	// Selector returns the resource's selector.
	Selector() string
	// Namespace returns the namespace where the resource lives in.
	Namespace() string
}

// Waiter groups a set of methods to wait for a certain job to be executed in
// regular periods of time until a given timeout.
type Waiter interface {
	Wait(timeout time.Duration, pollFrequency time.Duration) error
}
