package runsim

import (
	"context"
	"gitlab.com/ignitionrobotics/web/cloudsim/ign-transport/proto/ignition/msgs"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/transport"
	ignws "gitlab.com/ignitionrobotics/web/cloudsim/pkg/transport/ign"
	"sync"
	"time"
)

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
	// Transport has a reference to a publisher/subscriber transporting mechanism using websockets.
	Transport ignws.PubSubWebsocketTransporter
}

// IsExpired returns true if the RunningSimulation has expired.
func (rs *RunningSimulation) IsExpired() bool {
	return time.Now().After(rs.MaxValidUntil)
}

// ReadWarmup is the callback passed to the websocket client that will be invoked each time
// a message is received at the /subt/start topic.
func (rs *RunningSimulation) ReadWarmup(ctx context.Context, msg transport.Message) error {
	var m msgs.StringMsg
	err := msg.GetPayload(&m)
	if err != nil {
		return err
	}

	rs.lockCurrentState.Lock()
	defer rs.lockCurrentState.Unlock()

	if !rs.Finished && m.GetData() == "recording_complete" {
		rs.Finished = true
	}

	return nil
}

// NewRunningSimulation initializes a new RunningSimulation identified by the given groupID that will run for a maximum
// amount of maxSimSeconds seconds and will be valid for the duration given in validFor.
func NewRunningSimulation(sim simulations.Simulation) *RunningSimulation {
	launchedAt := time.Now()
	if sim.GetLaunchedAt() != nil {
		launchedAt = *sim.GetLaunchedAt()
	}

	return &RunningSimulation{
		GroupID:       sim.GetGroupID(),
		currentState:  stateUnknown,
		publishing:    false,
		CreatedAt:     launchedAt,
		MaxValidUntil: launchedAt.Add(sim.GetValidFor()),
		Finished:      false,
	}
}
