package ec2

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/stretchr/testify/suite"
	cloud "gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud/aws"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/machines"
	"gitlab.com/ignitionrobotics/web/ign-go/v6"
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
	workerGroupName := "test"
	s.ec2API = &mockEC2Count{
		WorkerGroupName: workerGroupName,
	}
	logger := ign.NewLoggerNoRollbar("ec2CountMachinesTestSuite", ign.VerbosityDebug)
	var err error
	s.machines, err = NewMachines(&NewInput{
		API:             s.ec2API,
		CostCalculator:  cloud.NewCostCalculatorEC2(nil),
		Logger:          logger,
		WorkerGroupName: workerGroupName,
		Zones: []Zone{
			{
				Zone:     "test",
				SubnetID: "test",
			},
		},
	})
	s.Require().NoError(err)
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
	result := s.machines.Count(machines.CountMachinesInput{})
	s.Equal(3, result)
}

func (s *ec2CountMachinesTestSuite) TestCount_GetMachinesWithFilters() {
	s.ec2API.InternalError = nil
	s.ec2API.ReturnMachines = true
	result := s.machines.Count(machines.CountMachinesInput{
		Filters: map[string][]string{
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
	InternalError   error
	ReturnMachines  bool
	WorkerGroupName string
}

func (m *mockEC2Count) DescribeInstances(input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	if m.InternalError != nil {
		return nil, m.InternalError
	}

	// Verify that the Machines identification tag is included
	tagFound := false
	for i, filter := range input.Filters {
		if filter != nil && filter.Name != nil && *filter.Name == "tag:cloudsim-simulation-worker" {
			tagFound = true
			// Remove the tag from the filters to avoid breaking the following mocks
			input.Filters = append(input.Filters[0:i], input.Filters[i+1:len(input.Filters)]...)
			break
		}
	}
	if !tagFound {
		return nil, errors.New("call to DescribeInstances is missing worker-group-name tag")
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
