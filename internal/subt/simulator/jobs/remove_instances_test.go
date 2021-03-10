package jobs

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/stretchr/testify/suite"
	subtapp "gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/subt/simulator/state"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud/fake"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/machines"
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
	Filters      map[string][]string
}

func TestRemoveInstances(t *testing.T) {
	suite.Run(t, new(removeInstancesTestSuite))
}

func (s *removeInstancesTestSuite) SetupSuite() {
	s.GroupID = "aaaa-bbbb-cccc-dddd"

	s.Filters = make(map[string][]string)
	tags := subtapp.GetTagsInstanceBase(s.GroupID)

	for _, tag := range tags {
		for k, v := range tag.Map {
			s.Filters[fmt.Sprintf("tag:%s", k)] = []string{v}
		}
	}
}

func (s *removeInstancesTestSuite) SetupTest() {
	s.Logger = ign.NewLoggerNoRollbar("TestRemoveInstances", ign.VerbosityDebug)

	s.Machines = fake.NewMachines()

	s.Platform = platform.NewPlatform(platform.Components{
		Machines: s.Machines,
	})

	s.InitialState = state.NewStopSimulation(s.Platform, nil, s.GroupID)

	s.Store = actions.NewStore(s.InitialState)
}

func (s *removeInstancesTestSuite) TestRemoveInstancesFails() {
	s.Machines.On("Terminate", machines.TerminateMachinesInput{
		Filters: s.Filters,
	}).Return(errors.New("some test error"))

	_, err := RemoveInstances.Run(s.Store, nil, nil, s.InitialState)

	// If terminating instances fails, we don't want to have an error. The reason for this is strictly related to
	// avoid having a termination process that has stopped in the middle, leaving some machines undeleted,
	// because a single termination request failed.
	s.Assert().NoError(err)
}

func (s *removeInstancesTestSuite) TestRemoveInstancesSuccess() {
	s.Machines.On("Terminate", machines.TerminateMachinesInput{
		Filters: s.Filters,
	}).Return(error(nil))

	_, err := RemoveInstances.Run(s.Store, nil, nil, s.InitialState)
	s.Assert().NoError(err)
}
