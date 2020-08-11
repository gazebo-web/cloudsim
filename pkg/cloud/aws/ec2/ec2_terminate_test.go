package ec2

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/ign-go"
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
	logger := ign.NewLoggerNoRollbar("ec2TerminateMachinesTestSuite", ign.VerbosityDebug)
	s.machines = NewMachines(s.ec2API, logger)
}

func (s *ec2TerminateMachinesTestSuite) TestTerminate_ErrorWhenNilMachineNames() {
	err := s.machines.Terminate(cloud.TerminateMachinesInput{
		Names:   nil,
		Retries: 10,
	})
	s.Error(err)
	s.Equal(cloud.ErrMissingMachineNames, err)
}

func (s *ec2TerminateMachinesTestSuite) TestTerminate_ErrorWhenEmptyMachineNames() {
	err := s.machines.Terminate(cloud.TerminateMachinesInput{
		Names:   []string{},
		Retries: 10,
	})
	s.Error(err)
	s.Equal(cloud.ErrMissingMachineNames, err)
}

func (s *ec2TerminateMachinesTestSuite) TestTerminate_ErrorWithDryRunMode() {
	s.ec2API.InternalError = errors.New("test error")
	err := s.machines.Terminate(cloud.TerminateMachinesInput{
		Names:   []string{"machines-id"},
		Retries: 3,
	})
	s.Error(err)
	s.Equal(cloud.ErrUnknown, err)
}

func (s *ec2TerminateMachinesTestSuite) TestTerminate_ValidWithDryRunMode() {
	mock := &mockEC2TerminateDryRun{}
	logger := ign.NewLoggerNoRollbar("ec2TerminateMachinesTestSuite", ign.VerbosityDebug)
	s.machines = NewMachines(mock, logger)
	err := s.machines.Terminate(cloud.TerminateMachinesInput{
		Names:   []string{"machine-id"},
		Retries: 5,
	})
	s.NoError(err)
	s.Equal(3, mock.TerminateInstancesCalls)

}

type mockEC2Terminate struct {
	ec2iface.EC2API
	InternalError error
}

// TerminateInstances mocks EC2 TerminateInstances method.
func (m *mockEC2Terminate) TerminateInstances(input *ec2.TerminateInstancesInput) (*ec2.TerminateInstancesOutput, error) {
	if m.InternalError != nil {
		return nil, m.InternalError
	}
	return nil, nil
}

type mockEC2TerminateDryRun struct {
	ec2iface.EC2API
	TerminateInstancesCalls int
}

// TerminateInstances mocks EC2 TerminateInstances method.
func (m *mockEC2TerminateDryRun) TerminateInstances(input *ec2.TerminateInstancesInput) (*ec2.TerminateInstancesOutput, error) {
	if m.TerminateInstancesCalls == 0 {
		m.TerminateInstancesCalls++
		return nil, errors.New("test error")
	}
	if m.TerminateInstancesCalls == 1 {
		m.TerminateInstancesCalls++
		return nil, awserr.New(ErrCodeDryRunOperation, "dry run operation", errors.New("dry run error"))
	}
	m.TerminateInstancesCalls++
	return nil, nil
}
