package simulator

import (
	"context"
	igntransport "gitlab.com/ignitionrobotics/web/cloudsim/third_party/ign-transport"
	"sync"
	"time"
)

// RunningSimulation represents a running simulation. It is created by the
// simulation service when a simulation is lauched. It holds the current state
// reported by gazebo and also holds an ign-transport node to interact with gazebo (pub/sub).
// It uses the given simulation GroupID as the ign-transport's Partition.
type RunningSimulation struct {
	// The simulation GroupID assigned by the sim_service
	GroupID string
	// The user (or Org) that launched this simulation
	Owner string
	// The last reported state by gazebo
	currentState gazeboState
	// The desired state. Based on it, the RunningSimulation will use ign-transport
	// to request gazebo to switch to that state
	desiredState gazeboState
	// The ign-transport node to interact with Gazebo
	ignTransportNode *igntransport.GoIgnTransportNode
	// An internal flag to show if we are sending ign-transport messages
	publishing bool
	// A mutex used to guard access to currentState field
	lockCurrentState sync.RWMutex
	// A mutex used to guard access to desiredState field
	lockDesiredState sync.RWMutex
	// SimTimeSeconds tracks the current "sim seconds" reported by the simulation /stats topic.
	SimTimeSeconds int64
	// SimWarmupSeconds holds the "sim seconds" value when the simulation notified
	// the warmup was finished.
	SimWarmupSeconds int64
	// SimMaxAllowedSeconds allows to configure an Expiration time based on the simulation time.
	SimMaxAllowedSeconds int64
	// SimCreatedAtTime keeps track of the entire simulation group launch time.
	SimCreatedAtTime time.Time
	// MaxValidUntil keeps track of the Max time this simulation should be automatically
	// terminated. It is used to avoid dangling simulations.
	MaxValidUntil time.Time
	// Finished indicates if the simulation has finished running. A "finished" message in the warmup topic will mark the
	// simulation as finished.
	Finished bool
	// stdoutSkipStatsMsgsCount is an internal variable used to control throttling while printing debug messages to stdout
	stdoutSkipStatsMsgsCount int
}

type gazeboState string

const (
	GazeboStateUnknown        gazeboState = "unknown"
	GazeboStateRun            gazeboState = "run"
	GazeboStatePause          gazeboState = "pause"
	stdoutSkipStatsMsgs             = 100
)

// Free releases the resources of this running simulation.
func (s *RunningSimulation) Free(ctx context.Context) {
	// logger(ctx).Info(fmt.Sprintf("RunningSimulation groupID[%s] Free() invoked", s.GroupID))
	s.publishing = false
	if s.ignTransportNode != nil {
		s.ignTransportNode.Free()
	}
	s.ignTransportNode = nil
}