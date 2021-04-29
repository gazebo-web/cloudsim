package ec2

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/machines"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"testing"
)

func TestListMachines(t *testing.T) {
	suite.Run(t, new(ec2ListMachinesTestSuite))
}

type ec2ListMachinesTestSuite struct {
	suite.Suite
	ec2API   *mockEC2List
	machines machines.Machines
}

type mockEC2List struct {
	ec2iface.EC2API
	DescribeInstancesStatusCalls int
	InternalError                error
}

func (mock *mockEC2List) DescribeInstanceStatus(input *ec2.DescribeInstanceStatusInput) (*ec2.DescribeInstanceStatusOutput, error) {
	defer func() {
		mock.DescribeInstancesStatusCalls++
	}()

	if mock.InternalError != nil {
		return nil, mock.InternalError
	}

	return &ec2.DescribeInstanceStatusOutput{
		InstanceStatuses: []*ec2.InstanceStatus{
			{
				InstanceId: aws.String("test"),
				InstanceState: &ec2.InstanceState{
					Code: aws.Int64(16),
					Name: aws.String("running"),
				},
			},
		},
	}, nil
}

func (s *ec2ListMachinesTestSuite) SetupTest() {
	s.ec2API = &mockEC2List{}
	logger := ign.NewLoggerNoRollbar("ec2ListMachinesTestSuite", ign.VerbosityDebug)
	var err error
	s.machines, err = NewMachines(&NewInput{
		API: s.ec2API,
		Logger: logger,
		Zones: []Zone{
			{
				Zone: "test",
				SubnetID: "test",
			},
		},
	})
	s.Require().NoError(err)
}

func (s *ec2ListMachinesTestSuite) TestList_WithError() {
	s.ec2API.InternalError = errors.New("test error")
	_, err := s.machines.List(machines.ListMachinesInput{})
	s.Assert().Error(err)
}

func (s *ec2ListMachinesTestSuite) TestList_Success() {
	out, err := s.machines.List(machines.ListMachinesInput{})
	s.Require().NoError(err)

	s.Assert().Len(out.Instances, 1)
	s.Assert().Equal("test", out.Instances[0].InstanceID)
	s.Assert().Equal("running", out.Instances[0].State)
}
