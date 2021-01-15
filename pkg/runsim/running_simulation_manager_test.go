package runsim

import (
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	ignws "gitlab.com/ignitionrobotics/web/cloudsim/transport/ign"
	"testing"
	"time"
)

func TestManager(t *testing.T) {
	suite.Run(t, new(managerTestSuite))
}

type managerTestSuite struct {
	suite.Suite
	manager *manager
}

func (s *managerTestSuite) SetupTest() {
	s.manager = &manager{
		transporters:       make(map[simulations.GroupID]ignws.PubSubWebsocketTransporter),
		runningSimulations: make(map[simulations.GroupID]*RunningSimulation),
	}
}

func (s *managerTestSuite) TestAdd() {
	// Before adding, underlying maps should be empty.
	s.Require().Len(s.manager.transporters, 0)
	s.Require().Len(s.manager.runningSimulations, 0)

	gid := simulations.GroupID("aaaa-bbbb-dddd-eeee")
	rs := RunningSimulation{}
	t := ignws.NewPubSubTransporterMock()

	// Adding a running simulation should not return an error
	err := s.manager.Add(gid, &rs, t)
	s.Assert().NoError(err)

	// The underlying maps should have 1 element
	s.Assert().Len(s.manager.transporters, 1)
	s.Assert().Len(s.manager.runningSimulations, 1)

	// Adding a running simulation with the same group id should return an error.
	err = s.manager.Add(gid, &rs, t)
	s.Assert().Error(err)

	// If rs and/or t are nil, it should return an error.
	err = s.manager.Add("test", nil, t)
	s.Assert().Error(err)

	err = s.manager.Add("test", &rs, nil)
	s.Assert().Error(err)

	err = s.manager.Add("test", nil, nil)
	s.Assert().Error(err)
}

func (s *managerTestSuite) TestListExpiredSimulations() {
	rs := RunningSimulation{
		SimTimeSeconds:       0,
		SimWarmupSeconds:     5,
		SimMaxAllowedSeconds: 90,
		MaxValidUntil:        time.Now().Add(1 * time.Hour),
	}
	s.manager.runningSimulations["test"] = &rs

	s.Assert().Len(s.manager.ListExpiredSimulations(), 0)

	rs.SimTimeSeconds += 100

	s.Assert().Len(s.manager.ListExpiredSimulations(), 1)
}

func (s *managerTestSuite) TestListFinishedSimulations() {
	rs := RunningSimulation{
		SimTimeSeconds:       0,
		SimWarmupSeconds:     5,
		SimMaxAllowedSeconds: 90,
		Finished:             false,
		MaxValidUntil:        time.Now().Add(1 * time.Hour),
	}
	s.manager.runningSimulations["test"] = &rs

	s.Assert().Len(s.manager.ListFinishedSimulations(), 0)

	rs.Finished = true

	s.Assert().Len(s.manager.ListFinishedSimulations(), 1)
}

func (s *managerTestSuite) TestGetTransporter() {
	t := ignws.NewPubSubTransporterMock()

	s.manager.transporters["test"] = t

	output := s.manager.GetTransporter("test")
	s.Assert().Equal(t, output)

	output = s.manager.GetTransporter("test2")
	s.Assert().Nil(output)
}

func (s *managerTestSuite) TestFree() {
	rs := RunningSimulation{publishing: true}
	t := ignws.NewPubSubTransporterMock()

	// First returns true
	t.On("IsConnected").Once().Return(true)

	// After the transporter gets disconnected, return false.
	t.On("IsConnected").Once().Return(false)

	t.On("Disconnect").Once()

	s.manager.runningSimulations["test"] = &rs
	s.manager.transporters["test"] = t

	s.manager.Free("test")

	s.Assert().False(rs.publishing)

	// Don't panic unless Disconnect is being called again (that should not happen)
	s.NotPanics(func() {
		s.manager.Free("test")
	})
}

func (s *managerTestSuite) TestRemove() {
	rs := RunningSimulation{publishing: true}
	t := ignws.NewPubSubTransporterMock()

	s.manager.runningSimulations["test"] = &rs
	s.manager.transporters["test"] = t

	t.On("IsConnected").Once().Return(true)
	err := s.manager.Remove("test")
	s.Assert().Error(err)

	t.On("IsConnected").Once().Return(false)

	err = s.manager.Remove("test")
	s.Assert().NoError(err)

	err = s.manager.Remove("test")
	s.Assert().Error(err)
}
