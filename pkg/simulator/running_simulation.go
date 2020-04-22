package simulator

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/logger"
	igntransport "gitlab.com/ignitionrobotics/web/cloudsim/third_party/ign-transport"
	msgs "gitlab.com/ignitionrobotics/web/cloudsim/third_party/ign-transport/proto/ignition/msgs"
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

// NewRunningSimulationInput
type NewRunningSimulationInput struct {
	GroupID string
	Owner string
	MaxSeconds int64
	ValidFor	time.Duration
	worldStatsTopic string
	worldWarmupTopic string
}

// NewRunningSimulation
func NewRunningSimulation(ctx context.Context, input NewRunningSimulationInput) (*RunningSimulation, error) {
	s := RunningSimulation{
		GroupID:              input.GroupID,
		Owner:                input.Owner,
		currentState:         GazeboStateUnknown,
		lockCurrentState:     sync.RWMutex{},
		lockDesiredState:     sync.RWMutex{},
		publishing:           false,
		SimCreatedAtTime:     time.Now(),
		MaxValidUntil:        time.Now().Add(input.ValidFor),
		SimMaxAllowedSeconds: input.MaxSeconds,
	}
	var err error
	if s.ignTransportNode, err = igntransport.NewIgnTransportNode(&input.GroupID); err != nil {
		return nil, err
	}

	// TODO: Create a new logger from context

	// 	create a new specific logger for this running simulation
	//	reqID := fmt.Sprintf("RunningSimulation-sim-%s", groupID)
	//	newLogger := logger(ctx).Clone(reqID)
	//	Override logger
	//	ctx = ign.NewContextWithLogger(ctx, newLogger)

	_ = s.ignTransportNode.IgnTransportSubscribe(input.worldStatsTopic, func(msg []byte, msgType string) {
		s.callbackWorldStats(ctx, msg, msgType)
	})

	_ = s.ignTransportNode.IgnTransportSubscribe(input.worldWarmupTopic, func(msg []byte, msgType string) {
		s.callbackWarmup(ctx, msg, msgType)
	})

	return &s, nil
}

// gazeboState
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

// callbackWorldStats is the callback passed to ign-transport. It will be invoked
// each time a message is received in the topic associated to this node's groupID.
func (s *RunningSimulation) callbackWorldStats(ctx context.Context, msg []byte, msgType string) {

	ws := msgs.WorldStatistics{}
	var err error
	if err = proto.Unmarshal(msg, &ws); err != nil {
		// do nothing . Just log it
		logger.Logger(ctx).Error(fmt.Sprintf("RunningSimulation groupID[%s]- error while unmarshalling WorldStats msg. Got type[%s]. Msg[%s]", s.GroupID, msgType, msg), err)
		return
	}

	// Simple attempt to control throttling while printing debug messages to stdout
	s.stdoutSkipStatsMsgsCount++
	if s.stdoutSkipStatsMsgsCount > stdoutSkipStatsMsgs {
		s.stdoutSkipStatsMsgsCount = 0
		logger.Logger(ctx).Debug(fmt.Sprintf("RunningSimulation groupID[%s]- WorldStats message received. Parsed struct: [%v]", s.GroupID, ws))
	}

	s.lockCurrentState.Lock()
	defer s.lockCurrentState.Unlock()
	if ws.Paused {
		s.currentState = GazeboStatePause
	} else {
		s.currentState = GazeboStateRun
	}

	// Also update the reported Sim time
	s.SimTimeSeconds = ws.SimTime.Sec
}

// callbackWarmup is the callback passed to ign-transport that will be invoked each time
// a message is received at the /warmup/ready topic.
func (s *RunningSimulation) callbackWarmup(ctx context.Context, msg []byte, msgType string) {
	wup := msgs.StringMsg{}
	var err error
	if err = proto.Unmarshal(msg, &wup); err != nil {
		// do nothing . Just log it
		logger.Logger(ctx).Error(fmt.Sprintf("RunningSimulation groupID[%s]- error while unmarshalling Warmup msg. Got type[%s]. Msg[%s]", s.GroupID, msgType, msg), err)
		return
	}

	if wup.Data == "started" {
		// We only act the first time we receive this message
		if s.SimWarmupSeconds == 0 {
			logger.Logger(ctx).Info(fmt.Sprintf("RunningSimulation groupID[%s]- Warmup message received. Parsed struct: [%v]", s.GroupID, wup))

			s.SimWarmupSeconds = s.SimTimeSeconds
		}
	} else if !s.Finished && wup.Data == "finished" {
		logger.Logger(ctx).Info(fmt.Sprintf("RunningSimulation groupID[%s]- Finished message received. Parsed struct: [%v]", s.GroupID, wup))

		s.Finished = true
	}
}

// SendMessage publishes a string message to an specific topic.
func (s *RunningSimulation) SendMessage(ctx context.Context, topic, msg, msgType string) {
	logger.Logger(ctx).Info(fmt.Sprintf("RunningSimulation groupID[%s]- publish msg [%s] to topic [%s] with type [%s]", s.GroupID, msg, topic, msgType))
	if s.ignTransportNode != nil {
		_ = s.ignTransportNode.IgnTransportPublishStringMsg(topic, msg)
	}
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