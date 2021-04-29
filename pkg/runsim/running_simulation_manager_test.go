package runsim

import (
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	ignws "gitlab.com/ignitionrobotics/web/cloudsim/pkg/transport/ign"
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
		runningSimulations: make(map[simulations.GroupID]*RunningSimulation),
	}
}

func (s *managerTestSuite) TestAdd() {
	// Before adding, underlying maps should be empty.
	s.Require().Len(s.manager.runningSimulations, 0)

	gid := simulations.GroupID("aaaa-bbbb-dddd-eeee")
	rs := RunningSimulation{}
	t := ignws.NewPubSubTransporterMock()

	// Adding a running simulation should not return an error
	err := s.manager.Add(gid, &rs, t)
	s.Assert().NoError(err)

	// The underlying maps should have 1 element
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
	// We add a running simulation before running tests.
	rs := RunningSimulation{
		SimTimeSeconds:       0,
		SimWarmupSeconds:     5,
		SimMaxAllowedSeconds: 90,
		MaxValidUntil:        time.Now().Add(1 * time.Hour),
	}
	s.manager.runningSimulations["test"] = &rs

	// The running simulation isn't not expired yet
	s.Assert().Len(s.manager.ListExpiredSimulations(), 0)

	// We force the running simulation to be expired
	rs.SimTimeSeconds += 100

	// Now listing expired simulations returns an entry.
	s.Assert().Len(s.manager.ListExpiredSimulations(), 1)
}

func (s *managerTestSuite) TestListFinishedSimulations() {
	// We add a running simulation before running tests.
	rs := RunningSimulation{
		SimTimeSeconds:       0,
		SimWarmupSeconds:     5,
		SimMaxAllowedSeconds: 90,
		Finished:             false,
		MaxValidUntil:        time.Now().Add(1 * time.Hour),
	}
	s.manager.runningSimulations["test"] = &rs

	// The running simulation hasn't finished yet
	s.Assert().Len(s.manager.ListFinishedSimulations(), 0)

	// We mark the running simulation as finished
	rs.Finished = true

	// Listing running simulations that have finished now returns an entry.
	s.Assert().Len(s.manager.ListFinishedSimulations(), 1)
}

func (s *managerTestSuite) TestGetTransporter() {
	t := ignws.NewPubSubTransporterMock()

	s.manager.runningSimulations["test"] = &RunningSimulation{Transport: t}

	output := s.manager.GetTransporter("test")
	s.Assert().Equal(t, output)

	output = s.manager.GetTransporter("test2")
	s.Assert().Nil(output)
}

func (s *managerTestSuite) TestFree() {
	t := ignws.NewPubSubTransporterMock()
	rs := RunningSimulation{publishing: true, Transport: t}

	// First returns true
	t.On("IsConnected").Once().Return(true)

	// After the transporter gets disconnected, return false.
	t.On("IsConnected").Once().Return(false)

	// Disconnect should be called only once.
	t.On("Disconnect").Once()

	s.manager.runningSimulations["test"] = &rs

	s.manager.Free("test")

	s.Assert().False(rs.publishing)

	// Don't panic unless Disconnect is being called again (that should not happen) if the ws client has disconnected.
	s.NotPanics(func() {
		s.manager.Free("test")
	})
}

func (s *managerTestSuite) TestRemove() {
	// Add data before running tests
	t := ignws.NewPubSubTransporterMock()
	rs := RunningSimulation{Transport: t}

	s.manager.runningSimulations["test"] = &rs

	// We should not be able to remove a simulation that has a connection.
	t.On("IsConnected").Once().Return(true)
	err := s.manager.Remove("test")
	s.Assert().Error(err)

	// But if the simulation is not longer connected
	t.On("IsConnected").Once().Return(false)

	// We can safely remove it from the list of running simulations
	err = s.manager.Remove("test")
	s.Assert().NoError(err)

	// Although removing again or if the entry does not exist should return an error.
	err = s.manager.Remove("test")
	s.Assert().Error(err)
}
