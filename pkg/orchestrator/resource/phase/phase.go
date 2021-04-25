package phase

// Phase represents a certain point in the lifecycle of a Resource.
type Phase string

const (
	// Pending is used to represent when a Resource is in the Pending Phase.
	// Used by: Pods, Nodes.
	Pending Phase = "Pending"
	// Running is used to represent when a Resource is in the Running Phase.
	// Used by Pods, Nodes.
	Running Phase = "Running"
	// Succeeded is used to represent when a Resource is in the Succeeded Phase.
	// Used by: Pods.
	Succeeded Phase = "Succeeded"
	// Failed is used to represent when a Resource is in the Failed Phase.
	// Used by: Pods.
	Failed Phase = "Failed"
	// Unknown is used to represent when a Resource is in an Unknown Phase.
	// Used by: Pods.
	Unknown Phase = "Unknown"
	// Terminated is used to represent when a Resource is in the Terminated Phase.
	// Used by: Nodes.
	Terminated Phase = "Terminated"
)

// ResourcePhase has a method to return the phase of a certain Resource.
type ResourcePhase interface {
	// Phase is a simple, high-level summary of where the Resource is in its lifecycle.
	Phase() Phase
}

type resourcePhase struct {
	// phase is a simple, high-level summary of where the Resource is in its lifecycle.
	phase Phase
}

// Phase is a simple, high-level summary of where the resource is in its lifecycle.
func (r *resourcePhase) Phase() Phase {
	return r.phase
}

// NewResourcePhase initializes a new ResourcePhase implementation.
func NewResourcePhase(phase Phase) ResourcePhase {
	return &resourcePhase{phase: phase}
}