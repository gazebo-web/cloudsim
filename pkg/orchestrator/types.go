package orchestrator

import (
	"fmt"
	"strings"
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
	// HasIPStatusCondition is used to indicate that pods have ips available.
	HasIPStatusCondition = Condition{
		Type:   "HasIPStatus",
		Status: "True",
	}
)

// Phase represents a certain point in the lifecycle of a Resource.
type Phase string

const (
	// PhasePending is used to represent when a Resource is on a Pending Phase.
	// Used by: Pods, Nodes.
	PhasePending Phase = "Pending"
	// PhaseRunning is used to represent when a Resource is on a Running Phase.
	// Used by Pods, Nodes.
	PhaseRunning Phase = "Running"
	// PhaseSucceeded is used to represent when a Resource is on a Succeeded Phase.
	// Used by: Pods.
	PhaseSucceeded Phase = "Succeeded"
	// PhaseFailed is used to represent when a Resource is on a Failed Phase.
	// Used by: Pods.
	PhaseFailed Phase = "Failed"
	// PhaseUnknown is used to represent when a Resource is on a Unknown Phase.
	// Used by: Pods.
	PhaseUnknown Phase = "Unknown"
	// PhaseTerminated is used to represent when a Resource is on a Terminated Phase.
	// Used by: Nodes.
	PhaseTerminated Phase = "Terminated"
)

// Selector is used to represent the state a certain resource.
type Selector interface {
	// String returns the selector represented in string format.
	String() string
	// Map returns the underlying selector's map.
	Map() map[string]string
	// Extend extends the underlying base map with the extension selector.
	// NOTE: If a certain key already exists in the base map, it will be overwritten by the extension value.
	Extend(extension Selector) Selector
	// Set sets the given value to the given key. If the key already exists, it will be overwritten.
	Set(key string, value string)
}

// selector is a group of key-pair values that identify a resource.
type selector map[string]string

// Set sets the given value to the given key. If the key already exists, it will be overwritten.
func (s selector) Set(key string, value string) {
	s[key] = value
}

// Extend extends the underlying base map with the extension selector.
// NOTE: If a certain key already exists in the base map, it will be overwritten by the extension value.
func (s selector) Extend(extension Selector) Selector {
	for k, v := range extension.Map() {
		s[k] = v
	}
	return s
}

// Map returns the selector in map format.
func (s selector) Map() map[string]string {
	return s
}

// String returns the selector in string format.
func (s selector) String() string {
	var out string
	var labels []string
	for key, value := range s {
		out = fmt.Sprintf("%s=%s", key, value)
		labels = append(labels, out)
	}
	return strings.Join(labels, ",")
}

// NewSelector initializes a new orchestrator.Selector from the given map.
// If `nil` is passed as input, an empty selector will be returned.
func NewSelector(input map[string]string) Selector {
	if input == nil {
		input = map[string]string{}
	}
	return selector(input)
}
