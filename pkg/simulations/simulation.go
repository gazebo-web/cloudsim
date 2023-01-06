package simulations

import (
	"errors"
	"github.com/gazebo-web/cloudsim/v4/pkg/calculator"
	"time"
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

	// ErrSimulationProcessed is returned when simulation is trying to be processed twice.
	ErrSimulationProcessed = errors.New("simulation has been processed")

	// ErrSimulationPlatformNotDefined is returned when a Simulation does not have a Platform defined.
	ErrSimulationPlatformNotDefined = errors.New("simulation has no platform defined")
)

// GroupID is an universally unique identifier that identifies a Simulation.
type GroupID string

// String returns the string representation of a GroupID.
func (gid GroupID) String() string {
	return string(gid)
}

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

	// StatusTerminateRequested is used when a simulation has been scheduled to be terminated.
	StatusTerminateRequested Status = "terminate-requested"

	// StatusDeletingPods is used when the pods of a certain simulation are being deleted.
	StatusDeletingPods Status = "deleting-pods"

	// StatusDeletingNodes is used when the nodes of a certain simulation are being deleted.
	StatusDeletingNodes Status = "deleting-nodes"

	// StatusProcessingResults is used when a simulation's score and stats are being extracted from a gazebo server.
	StatusProcessingResults Status = "processing-results"

	// StatusTerminatingInstances is used when simulation instances are being deleted.
	StatusTerminatingInstances Status = "terminating-instances"

	// StatusTerminated is used when a simulation has been terminated.
	StatusTerminated Status = "terminated"

	// StatusSuperseded is used when a simulation has been superseded.
	StatusSuperseded Status = "superseded"

	// StatusRestarted is used when a simulation has been restarted.
	StatusRestarted Status = "restarted"

	// StatusUnknown is used to represent an unknown status.
	StatusUnknown Status = "unknown"
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
	// GetGroupID returns the current Simulation's group id.
	GetGroupID() GroupID

	// GetStatus returns the current Simulation's status.
	GetStatus() Status

	// HasStatus checks if the current Simulation has a given status.
	HasStatus(status Status) bool

	// SetStatus sets a given status to the Simulation.
	SetStatus(status Status)

	// GetKind returns the current simulation's kind.
	GetKind() Kind

	// IsKind checks if the current Simulation is of the given kind.
	IsKind(Kind) bool

	// GetError returns the current simulation's error. It returns nil if the simulation doesn't have an error.
	GetError() *Error

	// GetImage returns the simulation's docker image. This image is used as the solution image.
	GetImage() string

	// GetLaunchedAt returns the time and date the simulation was officially launched. This date can differ from the
	// time the simulation was requested due to the simulation having been held, or because it has been unable to
	// launch because of insufficient cloud resources.
	GetLaunchedAt() *time.Time

	// GetValidFor returns the amount of time that the simulation is considered valid.
	GetValidFor() time.Duration

	// IsProcessed returns true if the Simulation has been already processed.
	IsProcessed() bool

	// GetOwner returns the Simulation's owner.
	GetOwner() *string

	// GetCreator returns the Simulation's creator.
	GetCreator() string

	// GetPlatform returns the Simulation's platform.
	GetPlatform() *string

	// SetRate sets the given rate to this simulation.
	SetRate(rate calculator.Rate)

	// GetRate returns the rate at which this simulation should be charged.
	GetRate() calculator.Rate

	// GetStoppedAt returns the date and time when a simulation stopped from running.
	GetStoppedAt() *time.Time

	// GetCost applies the current rate to this simulation resulting in the amount of money that it should be charged.
	GetCost() (uint, calculator.Rate, error)

	// GetChargedAt returns the time and date this simulation was charged.
	GetChargedAt() *time.Time
}
