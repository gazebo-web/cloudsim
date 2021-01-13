package runsim

import (
	"context"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	ignws "gitlab.com/ignitionrobotics/web/cloudsim/transport/ign"
	"sync"
	"time"
)

type Callback func(ctx context.Context, msg ignws.Message)

type Callbacks interface {
	readWorldStats(ctx context.Context, msg ignws.Message)
	readWarmup(ctx context.Context, msg ignws.Message)
}

type State string

type RunningSimulation struct {
	// GroupID has a reference to the simulations.GroupID value to identify a simulation.
	GroupID simulations.GroupID
	// currentState is the last reported state by gazebo
	currentState State
	// desiredState is used to send through a transport node a request to gazebo to apply this state.
	desiredState State
	// publishing is an internal flag to show if we are sending ign-transport messages
	publishing bool
	// lockCurrentState is a mutex used to guard access to currentState field
	lockCurrentState sync.RWMutex
	// lockDesiredState is a mutex used to guard access to desiredState field
	lockDesiredState sync.RWMutex
	// SimTimeSeconds tracks the current "sim seconds" reported by the simulation /stats topic.
	SimTimeSeconds int64
	// SimWarmupSeconds holds the "sim seconds" value when the simulation notified
	// the warmup was finished.
	SimWarmupSeconds int64
	// SimMaxAllowedSeconds allows to configure an Expiration time based on the simulation time.
	SimMaxAllowedSeconds int64
	// CreatedAt keeps track of the entire simulation group launch time.
	CreatedAt time.Time
	// MaxValidUntil keeps track of the Max time this simulation should be automatically
	// terminated. It is used to avoid dangling simulations.
	MaxValidUntil time.Time
	// Finished indicates if the simulation has finished running. A "finished" message in the warmup topic will mark the
	// simulation as finished.
	Finished bool
	// stdoutSkipStatsMsgsCount is an internal variable used to control throttling while printing debug messages to stdout
	stdoutSkipStatsMsgsCount int
}

func (rs *RunningSimulation) Free(ctx context.Context) error {
	panic("implement me")
}

func (rs *RunningSimulation) IsExpired() bool {
	panic("implement me")
}

func (rs *RunningSimulation) readWorldStats(ctx context.Context, msg ignws.Message) {
	panic("implement me")
}

func (rs *RunningSimulation) readWarmup(ctx context.Context, msg ignws.Message) {
	panic("implement me")
}
