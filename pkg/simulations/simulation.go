package simulations

import (
	"errors"
)

var (
	// ErrInvalidGroupID is returned when a group id is invalid.
	ErrInvalidGroupID = errors.New("invalid group id")

	// ErrIncorrectStatus is returned when a simulation status is not correct at the time it's being checked.
	ErrIncorrectStatus = errors.New("incorrect status")

	// ErrIncorrectKind is returned when a simulation kind is not correct at the time it's being checked.
	ErrIncorrectKind = errors.New("incorrect kind")

	// ErrSimulationWithError is returned when a simulation has an error.
	ErrSimulationWithError = errors.New("simulation with error")
)

// GroupID is an universally unique identifier that identifies a Simulation.
type GroupID string

// Status represents a stage a Simulation can be in.
type Status string

var (
	// StatusPending is used when a simulation is pending to be scheduled.
	StatusPending Status = "pending"

	// StatusRunning is used when a simulation is running.
	StatusRunning Status = "running"

	// StatusRejected is used when a simulation has been rejected.
	StatusRejected Status = "rejected"

	// StatusLaunchingInstances is used when a simulation has entered the launching instances phase.
	StatusLaunchingInstances Status = "launching-instances"

	// StatusLaunchingPods is used when a simulation has entered the launching pods phase.
	StatusLaunchingPods Status = "launching-pods"

	// StatusWaitingInstances is used when a simulation is waiting for instances to be launched.
	StatusWaitingInstances Status = "waiting-instances"

	// StatusWaitingNodes is used when a simulation is waiting for nodes to be ready
	StatusWaitingNodes Status = "waiting-nodes"

	// StatusWaitingPods is used when a simulation is waiting for pods to be ready.
	StatusWaitingPods Status = "waiting-pods"
)

// Kind is used to identify if a Simulation is a single simulation or a multisim.
// If a simulation is a multisim, different Kind values are used
// to identify if the simulations is a parent simulation or a child simulation.
type Kind uint

var (
	// SimSingle represents a single simulation.
	SimSingle Kind = 0
	// SimParent represents a parent simulation.
	SimParent Kind = 1
	// SimChild represents a child simulation.
	SimChild Kind = 2
)

// Error is used to assign an error to a simulation. Simulations with errors are forbidden to run.
type Error string

// Simulation groups a set of methods to identify a simulation.
type Simulation interface {
	// GroupID returns the current simulation's group id.
	GroupID() GroupID

	// Status returns the current simulation's status.
	Status() Status

	// Kind returns the current simulation's kind.
	Kind() Kind

	// Error returns the current simulation's error. It returns nil if the simulation doesn't have an error.
	Error() *Error

	// Image returns the simulation's docker image. This image is used as the solution image.
	Image() string

	// SetStatus sets a given status to the simulation.
	SetStatus(status Status)
}
