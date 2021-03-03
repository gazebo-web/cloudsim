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

func TestCountMachines(t *testing.T) {
	suite.Run(t, new(ec2CountMachinesTestSuite))
}

type ec2CountMachinesTestSuite struct {
	suite.Suite
	ec2API   *mockEC2Count
	machines machines.Machines
}

func (s *ec2CountMachinesTestSuite) SetupTest() {
	s.ec2API = &mockEC2Count{}
	logger := ign.NewLoggerNoRollbar("ec2CountMachinesTestSuite", ign.VerbosityDebug)
	s.machines = NewMachines(s.ec2API, logger)
}

func (s *ec2CountMachinesTestSuite) TestCount_ReturnZeroWhenThereAreNoMachines() {
	s.ec2API.InternalError = nil
	s.ec2API.ReturnMachines = false
	result := s.machines.Count(machines.CountMachinesInput{
		Filters: nil,
	})
	s.Equal(0, result)
}

func (s *ec2CountMachinesTestSuite) TestCount_ReturnErrorWhenThereIsAnInternalAWSError() {
	s.ec2API.InternalError = errors.New("test error")
	result := s.machines.Count(machines.CountMachinesInput{
		Filters: nil,
	})
	s.Equal(-1, result)
}

func (s *ec2CountMachinesTestSuite) TestCount_GetAllMachines() {
	s.ec2API.InternalError = nil
	s.ec2API.ReturnMachines = true
	result := s.machines.Count(machines.CountMachinesInput{
		Filters: nil,
	})
	s.Equal(3, result)
}

func (s *ec2CountMachinesTestSuite) TestCount_GetMachinesWithFilters() {
	s.ec2API.InternalError = nil
	s.ec2API.ReturnMachines = true
	result := s.machines.Count(machines.CountMachinesInput{
		Filters: map[string][]string{
			"tag:cloudsim-simulation-worker": {
				"name-prefix",
			},
			"instance-state-name": {
				"pending",
				"running",
			},
		},
	})
	s.Equal(1, result)
}

type mockEC2Count struct {
	ec2iface.EC2API
	InternalError  error
	ReturnMachines bool
}

func (m *mockEC2Count) DescribeInstances(input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	if m.InternalError != nil {
		return nil, m.InternalError
	}
	if !m.ReturnMachines {
		return &ec2.DescribeInstancesOutput{
			NextToken: aws.String("next-token"),
			Reservations: []*ec2.Reservation{
				{
					Instances: []*ec2.Instance{},
				},
			},
		}, nil
	}
	if len(input.Filters) > 0 {
		return &ec2.DescribeInstancesOutput{
			NextToken: aws.String("next-token"),
			Reservations: []*ec2.Reservation{
				{
					Instances: []*ec2.Instance{
						{
							InstanceId: aws.String("instance-a"),
						},
					},
				},
			},
		}, nil
	}
	return &ec2.DescribeInstancesOutput{
		NextToken: aws.String("next-token"),
		Reservations: []*ec2.Reservation{
			{
				Instances: []*ec2.Instance{
					{
						InstanceId: aws.String("instance-a"),
					},
					{
						InstanceId: aws.String("instance-b"),
					},
					{
						InstanceId: aws.String("instance-c"),
					},
				},
			},
		},
	}, nil
}
