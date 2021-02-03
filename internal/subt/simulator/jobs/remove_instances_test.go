package jobs

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud/fake"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/platform"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"testing"
)

type removeInstancesTestSuite struct {
	suite.Suite
	Machines     *fake.Machines
	Platform     platform.Platform
	Logger       ign.Logger
	API          ec2iface.EC2API
	GroupID      simulations.GroupID
	InitialState *state.StopSimulation
	Store        actions.Store
}

func TestRemoveInstances(t *testing.T) {
	suite.Run(t, new(removeInstancesTestSuite))
}

func (s *removeInstancesTestSuite) SetupTest() {
	s.Logger = ign.NewLoggerNoRollbar("TestRemoveInstances", ign.VerbosityDebug)

	s.Machines = fake.NewMachines()

	s.Platform = platform.NewPlatform(platform.Components{
		Machines: s.Machines,
	})

	s.GroupID = "aaaa-bbbb-cccc-dddd"

	s.InitialState = state.NewStopSimulation(s.Platform, nil, s.GroupID)

	s.Store = actions.NewStore(s.InitialState)
}

func (s *removeInstancesTestSuite) TestRemoveInstancesFails() {
	s.Machines.On("Terminate", cloud.TerminateMachinesInput{
		Filters: map[string][]string{
			fmt.Sprintf("tag:%s", instanceTagGroupID): {
				s.GroupID.String(),
			},
		},
	}).Return(errors.New("some test error"))

	_, err := RemoveInstances.Run(s.Store, nil, nil, s.InitialState)
	s.Assert().Error(err)
}

func (s *removeInstancesTestSuite) TestRemoveInstancesSuccess() {
	s.Machines.On("Terminate", cloud.TerminateMachinesInput{
		Filters: map[string][]string{
			fmt.Sprintf("tag:%s", instanceTagGroupID): {
				s.GroupID.String(),
			},
		},
	}).Return(error(nil))

	_, err := RemoveInstances.Run(s.Store, nil, nil, s.InitialState)
	s.Assert().NoError(err)
}
