package ec2

import (
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"testing"
)

func TestTerminateMachines(t *testing.T) {
	suite.Run(t, new(ec2TerminateMachinesTestSuite))
}

type ec2TerminateMachinesTestSuite struct {
	suite.Suite
	ec2API   *mockEC2Terminate
	machines cloud.Machines
}

func (s *ec2TerminateMachinesTestSuite) SetupTest() {
	s.ec2API = &mockEC2Terminate{}
	s.machines = NewMachines(s.ec2API)
}

func (s *ec2TerminateMachinesTestSuite) Test_ErrorWhenNilMachineNames() {
	err := s.machines.Terminate(cloud.TerminateMachinesInput{
		Names:  nil,
		DryRun: false,
	})
	s.Error(err)
	s.Equal(cloud.ErrMissingMachineNames, err)
}

func (s *ec2TerminateMachinesTestSuite) Test_ErrorWhenEmptyMachineNames() {
	err := s.machines.Terminate(cloud.TerminateMachinesInput{
		Names:  []string{},
		DryRun: false,
	})
	s.Error(err)
	s.Equal(cloud.ErrMissingMachineNames, err)
}

type mockEC2Terminate struct {
	ec2iface.EC2API
	InternalError error
}
