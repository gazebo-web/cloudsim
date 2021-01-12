package simulations

import (
	"context"
	"fmt"
	"gitlab.com/ignitionrobotics/web/cloudsim/ign-transport/proto/ignition/msgs"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/transport"
	ignws "gitlab.com/ignitionrobotics/web/cloudsim/pkg/transport/ign"
	"gitlab.com/ignitionrobotics/web/ign-go"
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
	// The websocket transport mechanism to communicate to the ign websocket server
	websocketTransportNode ignws.PubSubWebsocketTransporter
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
	stateUnknown        gazeboState = "unknown"
	stateRun            gazeboState = "run"
	statePause          gazeboState = "pause"
	stdoutSkipStatsMsgs             = 100
)

// NewRunningSimulation creates a new running simulation.
// The worldStatsTopic arg is the topic to subscribe to get notifications about the
// simulation state (eg. /world/default/stats). The optional worldWarmupTopic
// is used to get notifications about the time when the Simulation actually started.
func NewRunningSimulation(ctx context.Context, dep *SimulationDeployment, t ignws.PubSubWebsocketTransporter, worldStatsTopic string,
	worldWarmupTopic string, maxSimSeconds int) (*RunningSimulation, error) {
	groupID := *dep.GroupID
	logger(ctx).Info(fmt.Sprintf("Creating new RunningSimulation for groupID[%s] with topics stats[%s] and maxSimSeconds[%d]", groupID, worldStatsTopic, maxSimSeconds))

	// Backward compatibility: we assume 30 minutes by default.
	var validFor time.Duration
	validFor, _ = time.ParseDuration(*dep.ValidFor)

	s := RunningSimulation{
		GroupID:                groupID,
		Owner:                  *dep.Owner,
		currentState:           stateUnknown,
		lockCurrentState:       sync.RWMutex{},
		lockDesiredState:       sync.RWMutex{},
		publishing:             false,
		SimCreatedAtTime:       time.Now(),
		MaxValidUntil:          time.Now().Add(validFor),
		SimMaxAllowedSeconds:   int64(maxSimSeconds),
		websocketTransportNode: t,
	}
	// create a new specific logger for this running simulation
	reqID := fmt.Sprintf("RunningSimulation-sim-%s", groupID)
	newLogger := logger(ctx).Clone(reqID)
	// Override logger
	ctx = ign.NewContextWithLogger(ctx, newLogger)

	// subscribe to the stats topic to know the play/pause status
	err := s.websocketTransportNode.Subscribe(worldStatsTopic, func(message transport.Message) {
		s.callbackWorldStats(ctx, message)
	})
	if err != nil {
		return nil, err
	}

	err = s.websocketTransportNode.Subscribe(worldWarmupTopic, func(message transport.Message) {
		s.callbackWarmup(ctx, message)
	})
	if err != nil {
		return nil, err
	}

	return &s, nil
}

// Free releases the resources of this running simulation.
func (s *RunningSimulation) Free(ctx context.Context) {
	logger(ctx).Info(fmt.Sprintf("RunningSimulation groupID[%s] Free() invoked", s.GroupID))
	s.publishing = false
	if s.websocketTransportNode != nil && s.websocketTransportNode.IsConnected() {
		s.websocketTransportNode.Disconnect()
	}
	s.websocketTransportNode = nil
}

// IsExpired returns true is the RunningSimulation is expired.
func (s *RunningSimulation) IsExpired() bool {
	secondsExpired := false
	// If SimMaxAllowedSeconds is 0 then there is no limit for Sim seconds
	if s.SimMaxAllowedSeconds > 0 {
		secondsExpired = (s.SimTimeSeconds - s.SimWarmupSeconds) > s.SimMaxAllowedSeconds
	}
	return secondsExpired || time.Now().After(s.MaxValidUntil)
}

// ResumeSimulation request to resume the gazebo simulation from pause state.
// Dev note: To do it, this node will send `resume` messages to Gazebo until the node receives
// a message saying the simulation is running.
func (s *RunningSimulation) ResumeSimulation(ctx context.Context) error {
	logger(ctx).Info(fmt.Sprintf("RunningSimulation groupID[%s]- ResumeSimulation invoked", s.GroupID))

	s.lockCurrentState.RLock()
	defer s.lockCurrentState.RUnlock()
	if s.currentState == stateRun || s.publishing {
		return nil
	}

	s.runSetGazeboState(ctx, stateRun)
	return nil
}

// PauseSimulation request Gazebo to pause the simulation.
func (s *RunningSimulation) PauseSimulation(ctx context.Context) error {
	logger(ctx).Info(fmt.Sprintf("RunningSimulation groupID[%s]- PauseSimulation invoked", s.GroupID))

	s.lockCurrentState.RLock()
	defer s.lockCurrentState.RUnlock()
	if s.currentState == statePause || s.publishing {
		return nil
	}

	s.runSetGazeboState(ctx, statePause)
	return nil
}

// runSetGazeboState is an internal loop in a separate go-routine to publish
// "set state" messages to gazebo until we can confirm we are in the desired state.
func (s *RunningSimulation) runSetGazeboState(ctx context.Context, newState gazeboState) {
	// If already publishing then return.
	if s.publishing {
		return
	}

	// Create a new go routine to publish the new desired state to gazebo
	go func() {
		for s.publishing {
			// Re-Read current and desired states to get latest values
			s.lockDesiredState.RLock()
			s.lockCurrentState.RLock()
			desired := s.desiredState
			current := s.currentState
			exit := false
			if current == desired {
				s.publishing = false
				exit = true
			}
			s.lockDesiredState.RUnlock()
			s.lockCurrentState.RUnlock()
			if exit {
				break
			}

			// Publish msg to gazebo
			if s.websocketTransportNode != nil {
				msg, msgType := buildGazeboSetStateMessage(ctx, desired)
				// TODO we need to update this with a call to ign' service "/world/default/control"
				pubMsg := ignws.NewPublicationMessage("/foo", msgType, msg)
				_ = s.websocketTransportNode.Publish(pubMsg)
			}
			// Wait some time before re checking
			Sleep(50 * time.Millisecond)
		}
	}()
}

func buildGazeboSetStateMessage(ctx context.Context, state gazeboState) (msg, msgType string) {
	msgType = "ignition.msgs.StringMsg"
	if state == stateRun {
		msg = "play"
	} else {
		msg = "pause"
	}
	return
}

// callbackWorldStats is the callback passed to the websocket client. It will be invoked
// each time a message is received in the topic associated to this node's groupID.
func (s *RunningSimulation) callbackWorldStats(ctx context.Context, msg transport.Message) {
	ws := msgs.WorldStatistics{}
	if err := msg.GetPayload(&ws); err != nil {
		// do nothing . Just log it
		logger(ctx).Error(fmt.Sprintf("RunningSimulation groupID[%s]- error while unmarshalling WorldStats msg. Got Msg[%s]", s.GroupID, msg), err)
		return
	}

	// Simple attempt to control throttling while printing debug messages to stdout
	s.stdoutSkipStatsMsgsCount++
	if s.stdoutSkipStatsMsgsCount > stdoutSkipStatsMsgs {
		s.stdoutSkipStatsMsgsCount = 0
		logger(ctx).Debug(fmt.Sprintf("RunningSimulation groupID[%s]- WorldStats message received. Parsed struct: [%s]", s.GroupID, ws.String()))
	}

	s.lockCurrentState.Lock()
	defer s.lockCurrentState.Unlock()
	if ws.Paused {
		s.currentState = statePause
	} else {
		s.currentState = stateRun
	}

	// Also update the reported Sim time
	s.SimTimeSeconds = ws.SimTime.Sec
}

// callbackWarmup is the callback passed to the websocket client that will be invoked each time
// a message is received at the /warmup/ready topic.
func (s *RunningSimulation) callbackWarmup(ctx context.Context, msg transport.Message) {
	wup := msgs.StringMsg{}
	if err := msg.GetPayload(&wup); err != nil {
		// do nothing . Just log it
		logger(ctx).Error(fmt.Sprintf("RunningSimulation groupID[%s]- error while unmarshalling Warmup msg. Got Msg[%s]", s.GroupID, msg), err)
		return
	}

	if wup.Data == "started" {
		// We only act the first time we receive this message
		if s.SimWarmupSeconds == 0 {
			logger(ctx).Info(fmt.Sprintf("RunningSimulation groupID[%s]- Warmup message received. Parsed struct: [%v]", s.GroupID, wup.String()))

			s.SimWarmupSeconds = s.SimTimeSeconds
		}
	} else if !s.Finished && wup.Data == "finished" {
		logger(ctx).Info(fmt.Sprintf("RunningSimulation groupID[%s]- Finished message received. Parsed struct: [%v]", s.GroupID, wup.String()))

		s.Finished = true
	}
}

// SendMessage publishes a string message to an specific topic.
func (s *RunningSimulation) SendMessage(ctx context.Context, topic, msg, msgType string) {
	logger(ctx).Info(fmt.Sprintf("RunningSimulation groupID[%s]- publish msg [%s] to topic [%s] with type [%s]", s.GroupID, msg, topic, msgType))
	if s.websocketTransportNode != nil {
		pubMsg := ignws.NewPublicationMessage(topic, msgType, msg)
		_ = s.websocketTransportNode.Publish(pubMsg)
	}
}
