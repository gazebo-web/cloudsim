package runsim

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/cloudsim/ign-transport/proto/ignition/msgs"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
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

func (s *runningSimulationTestSuite) TestHasReachedMaxSeconds() {
	rs := RunningSimulation{
		GroupID:                  simulations.GroupID("aaaa-bbbb-cccc-dddd"),
		currentState:             stateUnknown,
		publishing:               true,
		SimTimeSeconds:           100,
		SimWarmupSeconds:         10,
		SimMaxAllowedSeconds:     100,
		CreatedAt:                time.Now(),
		MaxValidUntil:            time.Now().Add(1 * time.Hour),
		Finished:                 false,
		stdoutSkipStatsMsgsCount: 0,
	}
	s.Assert().False(rs.hasReachedMaxSimSeconds())

	rs.SimTimeSeconds += 10
	s.Assert().False(rs.hasReachedMaxSimSeconds())

	rs.SimTimeSeconds++
	s.Assert().True(rs.hasReachedMaxSimSeconds())

}

func (s *runningSimulationTestSuite) TestIsExpired() {
	rs := RunningSimulation{
		GroupID:              "aaaa-bbbb-cccc-dddd",
		currentState:         stateUnknown,
		publishing:           true,
		SimTimeSeconds:       101,
		SimWarmupSeconds:     0,
		SimMaxAllowedSeconds: 100,
		CreatedAt:            time.Now(),
		MaxValidUntil:        time.Now().Add(time.Minute),
	}

	s.Assert().True(rs.IsExpired())

	rs = RunningSimulation{
		GroupID:              "aaaa-bbbb-cccc-dddd",
		currentState:         stateUnknown,
		publishing:           true,
		SimTimeSeconds:       100,
		SimWarmupSeconds:     10,
		SimMaxAllowedSeconds: 0,
		CreatedAt:            time.Now().Add(-1 * time.Hour),
		MaxValidUntil:        time.Now().Add(-30 * time.Minute),
	}

	s.Assert().True(rs.IsExpired())
}

func (s *runningSimulationTestSuite) TestReadWarmupWhenStarted() {
	rs := RunningSimulation{
		GroupID:              "aaaa-bbbb-cccc-dddd",
		currentState:         stateUnknown,
		publishing:           true,
		SimTimeSeconds:       30,
		SimWarmupSeconds:     0,
		SimMaxAllowedSeconds: 500,
		CreatedAt:            time.Now(),
		MaxValidUntil:        time.Now().Add(time.Minute),
	}

	m := msgs.StringMsg{
		Data: "started",
	}
	b, err := proto.Marshal(&m)
	s.Require().NoError(err)

	msg := ign.Message{
		Payload: string(b),
	}

	err = rs.ReadWarmup(context.Background(), &msg)
	s.Require().NoError(err)

	s.Assert().Equal(int64(30), rs.SimWarmupSeconds)

	// A second started message is received by error
	err = rs.ReadWarmup(context.Background(), &msg)
	s.Require().NoError(err)

	// Warmup seconds should remain the same.
	s.Assert().Equal(int64(30), rs.SimWarmupSeconds)
}

func (s *runningSimulationTestSuite) TestReadWarmupWhenFinished() {
	rs := RunningSimulation{
		GroupID:              "aaaa-bbbb-cccc-dddd",
		currentState:         stateUnknown,
		publishing:           true,
		SimTimeSeconds:       30,
		SimWarmupSeconds:     0,
		SimMaxAllowedSeconds: 500,
		CreatedAt:            time.Now(),
		MaxValidUntil:        time.Now().Add(time.Minute),
	}

	m := msgs.StringMsg{
		Data: "recording_complete",
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
