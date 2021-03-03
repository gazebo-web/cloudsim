package ec2

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/machines"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"testing"
)

func TestTerminateMachines(t *testing.T) {
	suite.Run(t, new(ec2TerminateMachinesTestSuite))
}

type ec2TerminateMachinesTestSuite struct {
	suite.Suite
	ec2API   *mockEC2Terminate
	machines machines.Machines
}

func (s *ec2TerminateMachinesTestSuite) SetupTest() {
	s.ec2API = &mockEC2Terminate{}
	logger := ign.NewLoggerNoRollbar("ec2TerminateMachinesTestSuite", ign.VerbosityDebug)
	s.machines = NewMachines(s.ec2API, logger)
}

func (s *ec2TerminateMachinesTestSuite) TestTerminate_ErrorWhenNilMachineNames() {
	err := s.machines.Terminate(machines.TerminateMachinesInput{
		Instances: nil,
	})
	s.Error(err)
	s.Equal(machines.ErrMissingMachineNames, err)
	s.Equal(0, s.ec2API.TerminateInstancesCalls)
}

func (s *ec2TerminateMachinesTestSuite) TestTerminate_ErrorWhenEmptyMachineNames() {
	err := s.machines.Terminate(machines.TerminateMachinesInput{
		Instances: []string{},
	})
	s.Error(err)
	s.Equal(machines.ErrMissingMachineNames, err)
	s.Equal(0, s.ec2API.TerminateInstancesCalls)
}

func (s *ec2TerminateMachinesTestSuite) TestTerminate_Valid() {
	err := s.machines.Terminate(machines.TerminateMachinesInput{
		Instances: []string{"machine-id"},
	})
	s.NoError(err)
	s.Equal(1, s.ec2API.TerminateInstancesCalls)
}

type mockEC2Terminate struct {
	ec2iface.EC2API
	InternalError           error
	TerminateInstancesCalls int
}

// TerminateInstances mocks EC2 TerminateInstances method.
func (m *mockEC2Terminate) TerminateInstances(input *ec2.TerminateInstancesInput) (*ec2.TerminateInstancesOutput, error) {
	if m.InternalError != nil {
		return nil, m.InternalError
	}
	m.TerminateInstancesCalls++
	return nil, nil
}
