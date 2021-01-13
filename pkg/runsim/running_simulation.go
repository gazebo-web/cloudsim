package runsim

import (
	"context"
	"gitlab.com/ignitionrobotics/web/cloudsim/ign-transport/proto/ignition/msgs"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	ignws "gitlab.com/ignitionrobotics/web/cloudsim/transport/ign"
	"sync"
	"time"
)

// Callback defines a function that is used to read incoming ignition websocket messages.
type Callback func(ctx context.Context, msg ignws.Message) error

// SimulationCallbacks has different Callback implementations used by RunningSimulation.
type SimulationCallbacks interface {
	readWorldStats(ctx context.Context, msg ignws.Message) error
	readWarmup(ctx context.Context, msg ignws.Message) error
}

// State defines a gazebo state. It's used to represent if a simulation is running or paused.
type State string

const (
	stateUnknown        State = "unknown"
	stateRun            State = "run"
	statePause          State = "pause"
	stdoutSkipStatsMsgs       = 100
)

// RunningSimulation is a representation of a simulations.Simulation that has been launched by workers.
type RunningSimulation struct {
	// GroupID has a reference to the simulations.GroupID value to identify a simulation.
	GroupID simulations.GroupID
	// currentState is the last reported state by gazebo
	currentState State
	// lockCurrentState is a mutex used to guard access to currentState field
	lockCurrentState sync.RWMutex
	// publishing is an internal flag to show if we are sending ign-transport messages
	publishing bool
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

// IsExpired returns true if the RunningSimulation has expired.
func (rs *RunningSimulation) IsExpired() bool {
	var secondsExpired bool
	if rs.SimMaxAllowedSeconds > 0 {
		secondsExpired = rs.hasReachedMaxSeconds()
	}
	return secondsExpired || time.Now().After(rs.MaxValidUntil)
}

// readWorldStats is the callback passed to the websocket client. It will be invoked
// each time a message is received in the topic associated to this node's groupID.
func (rs *RunningSimulation) readWorldStats(ctx context.Context, msg ignws.Message) error {
	var m msgs.WorldStatistics

	err := msg.GetPayload(&m)
	if err != nil {
		return err
	}

	rs.stdoutSkipStatsMsgsCount++
	if rs.stdoutSkipStatsMsgsCount > stdoutSkipStatsMsgs {
		rs.stdoutSkipStatsMsgsCount = 0
	}

	rs.lockCurrentState.Lock()
	defer rs.lockCurrentState.Unlock()

	if m.Paused {
		rs.currentState = statePause
	} else {
		rs.currentState = stateRun
	}

	rs.SimTimeSeconds = m.SimTime.Sec

	return nil
}

// hasReachedMaxSeconds defines if a simulation has reached the max allowed amount of seconds.
func (rs *RunningSimulation) hasReachedMaxSeconds() bool {
	return (rs.SimTimeSeconds - rs.SimWarmupSeconds) > rs.SimMaxAllowedSeconds
}

// readWarmup is the callback passed to the websocket client that will be invoked each time
// a message is received at the /warmup/ready topic.
func (rs *RunningSimulation) readWarmup(ctx context.Context, msg ignws.Message) error {
	var m msgs.StringMsg
	err := msg.GetPayload(&m)
	if err != nil {
		return err
	}

	if m.Data == "started" {
		if rs.SimWarmupSeconds == 0 {
			rs.SimWarmupSeconds = rs.SimTimeSeconds
		}
	}

	if !rs.Finished && m.Data == "finished" {
		rs.Finished = true
	}

	return nil
}
