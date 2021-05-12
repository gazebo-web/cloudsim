package runsim

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/cloudsim/ign-transport/proto/ignition/msgs"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/transport/ign"
	"testing"
	"time"
)

func TestRunningSimulation(t *testing.T) {
	suite.Run(t, new(runningSimulationTestSuite))
}

type runningSimulationTestSuite struct {
	suite.Suite
}

func (s *runningSimulationTestSuite) TestIsExpired() {
	rs := RunningSimulation{
		GroupID:       "aaaa-bbbb-cccc-dddd",
		currentState:  stateUnknown,
		publishing:    true,
		Started:       false,
		CreatedAt:     time.Now(),
		MaxValidUntil: time.Now().Add(time.Minute),
	}

	s.Assert().False(rs.IsExpired())

	rs = RunningSimulation{
		GroupID:       "aaaa-bbbb-cccc-dddd",
		currentState:  stateUnknown,
		publishing:    true,
		Started:       true,
		CreatedAt:     time.Now().Add(-1 * time.Hour),
		MaxValidUntil: time.Now().Add(-30 * time.Minute),
	}

	s.Assert().True(rs.IsExpired())
}

func (s *runningSimulationTestSuite) TestReadWarmupWhenStarted() {
	rs := RunningSimulation{
		GroupID:       "aaaa-bbbb-cccc-dddd",
		currentState:  stateUnknown,
		publishing:    true,
		CreatedAt:     time.Now(),
		MaxValidUntil: time.Now().Add(time.Minute),
	}

	m := msgs.StringMsg{
		Header: &msgs.Header{Stamp: &msgs.Time{Sec: int64(time.Now().Second())}},
		Data:   "started",
	}
	b, err := proto.Marshal(&m)
	s.Require().NoError(err)

	msg := ign.Message{
		Payload: string(b),
	}

	err = rs.ReadWarmup(context.Background(), &msg)
	s.Assert().NoError(err)
	s.Assert().True(rs.Started)
}

func (s *runningSimulationTestSuite) TestReadWarmupWhenFinished() {
	rs := RunningSimulation{
		GroupID:       "aaaa-bbbb-cccc-dddd",
		currentState:  stateUnknown,
		publishing:    true,
		CreatedAt:     time.Now(),
		MaxValidUntil: time.Now().Add(time.Minute),
	}
	m := msgs.StringMsg{
		Header: &msgs.Header{Stamp: &msgs.Time{Sec: int64(time.Now().Second())}},
		Data:   "recording_complete",
	}
	b, err := proto.Marshal(&m)
	s.Require().NoError(err)

	msg := ign.Message{
		Payload: string(b),
	}

	err = rs.ReadWarmup(context.Background(), &msg)
	s.Require().NoError(err)

	s.Assert().True(rs.Finished)
}
