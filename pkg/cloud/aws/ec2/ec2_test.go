package ec2

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"testing"
)

func TestMachines(t *testing.T) {
	suite.Run(t, new(ec2MachinesTestSuite))
}

type ec2MachinesTestSuite struct {
	suite.Suite
	ec2API   *mockEC2
	machines cloud.Machines
}

func (s *ec2MachinesTestSuite) SetupTest() {
	s.ec2API = &mockEC2{
		Mock: new(mock.Mock),
	}
	s.machines = NewMachines(s.ec2API)
}

func (s *ec2MachinesTestSuite) TestNewMachines() {
	e, ok := s.machines.(*machines)
	s.True(ok)
	s.NotNil(e.API)
}

func (s *ec2MachinesTestSuite) TestCreate_MissingKeyName() {
	input := []cloud.CreateMachinesInput{
		{
			DryRun:        false,
			KeyName:       "",
			MinCount:      1,
			MaxCount:      10,
			FirewallRules: nil,
			SubnetID:      "subnet-1234",
			Tags:          nil,
		},
	}
	_, err := s.machines.Create(input)
	s.Error(err)
	s.Equal(cloud.ErrMissingKeyName, err)
}

func (s *ec2MachinesTestSuite) TestCreate_InvalidCountBothZero() {
	input := []cloud.CreateMachinesInput{
		{
			DryRun:        false,
			KeyName:       "key-name",
			MinCount:      0,
			MaxCount:      0,
			FirewallRules: nil,
			SubnetID:      "subnet-1234",
			Tags:          nil,
		},
	}
	_, err := s.machines.Create(input)
	s.Error(err)
	s.Equal(cloud.ErrInvalidMachinesCount, err)
}

func (s *ec2MachinesTestSuite) TestCreate_InvalidCountMinCountZero() {
	input := []cloud.CreateMachinesInput{
		{
			DryRun:        false,
			KeyName:       "key-name",
			MinCount:      0,
			MaxCount:      1,
			FirewallRules: nil,
			SubnetID:      "subnet-1234",
			Tags:          nil,
		},
	}
	_, err := s.machines.Create(input)
	s.Error(err)
	s.Equal(cloud.ErrInvalidMachinesCount, err)
}

func (s *ec2MachinesTestSuite) TestCreate_InvalidCountMaxCountZero() {
	input := []cloud.CreateMachinesInput{
		{
			DryRun:        false,
			KeyName:       "key-name",
			MinCount:      1,
			MaxCount:      0,
			FirewallRules: nil,
			SubnetID:      "subnet-1234",
			Tags:          nil,
		},
	}
	_, err := s.machines.Create(input)
	s.Error(err)
	s.Equal(cloud.ErrInvalidMachinesCount, err)
}

func (s *ec2MachinesTestSuite) TestCreate_MinCountGreaterThanMaxCount() {
	input := []cloud.CreateMachinesInput{
		{
			DryRun:        false,
			KeyName:       "key-name",
			MinCount:      99,
			MaxCount:      1,
			FirewallRules: nil,
			SubnetID:      "subnet-1234",
			Tags:          nil,
		},
	}
	_, err := s.machines.Create(input)
	s.Error(err)
	s.Equal(cloud.ErrInvalidMachinesCount, err)
}

func (s *ec2MachinesTestSuite) TestCreate_MinCountEqualsMaxCount() {
	input := []cloud.CreateMachinesInput{
		{
			DryRun:        false,
			KeyName:       "key-name",
			MinCount:      1,
			MaxCount:      1,
			FirewallRules: nil,
			SubnetID:      "subnet-1234",
			Tags:          nil,
		},
	}
	_, err := s.machines.Create(input)
	s.NoError(err)
}

func (s *ec2MachinesTestSuite) TestCreate_NegativeCount() {
	input := []cloud.CreateMachinesInput{
		{
			DryRun:        false,
			KeyName:       "key-name",
			MinCount:      -100,
			MaxCount:      -25,
			FirewallRules: nil,
			SubnetID:      "subnet-1234",
			Tags:          nil,
		},
	}
	_, err := s.machines.Create(input)
	s.Error(err)
	s.Equal(cloud.ErrInvalidMachinesCount, err)
}

func (s *ec2MachinesTestSuite) TestCreate_InvalidSubnet() {
	input := []cloud.CreateMachinesInput{
		{
			DryRun:        false,
			KeyName:       "key-name",
			MinCount:      1,
			MaxCount:      99,
			FirewallRules: nil,
			SubnetID:      "test-1234",
			Tags:          nil,
		},
	}
	_, err := s.machines.Create(input)
	s.Error(err)
	s.Equal(cloud.ErrInvalidSubnetID, err)
}

func (s *ec2MachinesTestSuite) TestCreate_ValidSubnet() {
	input := []cloud.CreateMachinesInput{
		{
			DryRun:        false,
			KeyName:       "key-name",
			MinCount:      1,
			MaxCount:      99,
			FirewallRules: nil,
			SubnetID:      "subnet-0dae7657",
			Tags:          nil,
		},
	}
	_, err := s.machines.Create(input)
	s.NoError(err)

	input = []cloud.CreateMachinesInput{
		{
			DryRun:        false,
			KeyName:       "key-name",
			MinCount:      1,
			MaxCount:      99,
			FirewallRules: nil,
			SubnetID:      "subnet-06fe9fdb790aa78e7",
			Tags:          nil,
		},
	}
	_, err = s.machines.Create(input)
	s.NoError(err)
}

type mockEC2 struct {
	ec2iface.EC2API
	*mock.Mock
}

// RunInstances mocks EC2 RunInstances method.
func (m *mockEC2) RunInstances(input *ec2.RunInstancesInput) (*ec2.Reservation, error) {
	args := m.Called(input)
	reservation := args.Get(0).(*ec2.Reservation)
	err := args.Error(1)
	return reservation, err
}
