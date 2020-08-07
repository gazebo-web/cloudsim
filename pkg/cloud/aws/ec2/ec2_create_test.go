package ec2

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"testing"
)

func TestCreateMachines(t *testing.T) {
	suite.Run(t, new(ec2CreateMachinesTestSuite))
}

type ec2CreateMachinesTestSuite struct {
	suite.Suite
	ec2API   *mockEC2
	machines cloud.Machines
}

func (s *ec2CreateMachinesTestSuite) SetupTest() {
	s.ec2API = &mockEC2{}
	s.machines = NewMachines(s.ec2API)
}

func (s *ec2CreateMachinesTestSuite) TestNewMachines() {
	e, ok := s.machines.(*machines)
	s.True(ok)
	s.NotNil(e.API)
}

func (s *ec2CreateMachinesTestSuite) TestCreate_MissingKeyName() {
	input := []cloud.CreateMachinesInput{
		{
			DryRun:        false,
			KeyName:       "",
			MinCount:      1,
			MaxCount:      10,
			FirewallRules: nil,
			SubnetID:      "subnet-06fe9fdb790aa78e7",
			Tags:          nil,
		},
	}
	_, err := s.machines.Create(input)
	s.Error(err)
	s.Equal(cloud.ErrMissingKeyName, err)
	s.Equal(0, s.ec2API.RunInstancesCalls)
}

func (s *ec2CreateMachinesTestSuite) TestCreate_InvalidCountBothZero() {
	input := []cloud.CreateMachinesInput{
		{
			DryRun:        false,
			KeyName:       "key-name",
			MinCount:      0,
			MaxCount:      0,
			FirewallRules: nil,
			SubnetID:      "subnet-06fe9fdb790aa78e7",
			Tags:          nil,
		},
	}
	_, err := s.machines.Create(input)
	s.Error(err)
	s.Equal(cloud.ErrInvalidMachinesCount, err)
	s.Equal(0, s.ec2API.RunInstancesCalls)
}

func (s *ec2CreateMachinesTestSuite) TestCreate_InvalidCountMinCountZero() {
	input := []cloud.CreateMachinesInput{
		{
			DryRun:        false,
			KeyName:       "key-name",
			MinCount:      0,
			MaxCount:      1,
			FirewallRules: nil,
			SubnetID:      "subnet-06fe9fdb790aa78e7",
			Tags:          nil,
		},
	}
	_, err := s.machines.Create(input)
	s.Error(err)
	s.Equal(cloud.ErrInvalidMachinesCount, err)
	s.Equal(0, s.ec2API.RunInstancesCalls)
}

func (s *ec2CreateMachinesTestSuite) TestCreate_InvalidCountMaxCountZero() {
	input := []cloud.CreateMachinesInput{
		{
			DryRun:        false,
			KeyName:       "key-name",
			MinCount:      1,
			MaxCount:      0,
			FirewallRules: nil,
			SubnetID:      "subnet-06fe9fdb790aa78e7",
			Tags:          nil,
		},
	}
	_, err := s.machines.Create(input)
	s.Error(err)
	s.Equal(cloud.ErrInvalidMachinesCount, err)
	s.Equal(0, s.ec2API.RunInstancesCalls)
}

func (s *ec2CreateMachinesTestSuite) TestCreate_MinCountGreaterThanMaxCount() {
	input := []cloud.CreateMachinesInput{
		{
			DryRun:        false,
			KeyName:       "key-name",
			MinCount:      99,
			MaxCount:      1,
			FirewallRules: nil,
			SubnetID:      "subnet-06fe9fdb790aa78e7",
			Tags:          nil,
		},
	}
	_, err := s.machines.Create(input)
	s.Error(err)
	s.Equal(cloud.ErrInvalidMachinesCount, err)
	s.Equal(0, s.ec2API.RunInstancesCalls)
}

func (s *ec2CreateMachinesTestSuite) TestCreate_NegativeCount() {
	input := []cloud.CreateMachinesInput{
		{
			DryRun:        false,
			KeyName:       "key-name",
			MinCount:      -100,
			MaxCount:      -25,
			FirewallRules: nil,
			SubnetID:      "subnet-06fe9fdb790aa78e7",
			Tags:          nil,
		},
	}
	_, err := s.machines.Create(input)
	s.Error(err)
	s.Equal(cloud.ErrInvalidMachinesCount, err)
	s.Equal(0, s.ec2API.RunInstancesCalls)
}

func (s *ec2CreateMachinesTestSuite) TestCreate_InvalidSubnet() {
	input := []cloud.CreateMachinesInput{
		{
			DryRun:        false,
			KeyName:       "key-name",
			MinCount:      1,
			MaxCount:      99,
			FirewallRules: nil,
			SubnetID:      "subnet-1234",
			Tags:          nil,
		},
	}
	_, err := s.machines.Create(input)
	s.Error(err)
	s.Equal(cloud.ErrInvalidSubnetID, err)

	s.Equal(0, s.ec2API.RunInstancesCalls)
}

func (s *ec2CreateMachinesTestSuite) TestCreate_ValidWithoutDryRunMode() {
	input := []cloud.CreateMachinesInput{
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
	_, err := s.machines.Create(input)
	s.NoError(err)
	s.Equal(1, s.ec2API.RunInstancesCalls)
}

func (s *ec2CreateMachinesTestSuite) TestCreate_ValidWithDryRunMode() {
	mock := &mockEC2DryRunMode{}
	s.machines = NewMachines(mock)
	input := []cloud.CreateMachinesInput{
		{
			DryRun:        true,
			KeyName:       "key-name",
			MinCount:      1,
			MaxCount:      99,
			FirewallRules: nil,
			SubnetID:      "subnet-06fe9fdb790aa78e7",
			Tags:          nil,
			Retries:       3,
		},
	}
	_, err := s.machines.Create(input)
	s.NoError(err)
	s.Equal(3, mock.RunInstancesCalls)
}

type mockEC2 struct {
	ec2iface.EC2API
	RunInstancesCalls int
}

// RunInstances mocks EC2 RunInstances method.
func (m *mockEC2) RunInstances(input *ec2.RunInstancesInput) (*ec2.Reservation, error) {
	m.RunInstancesCalls++
	return &ec2.Reservation{}, nil
}

type mockEC2DryRunMode struct {
	ec2iface.EC2API
	RunInstancesCalls int
}

// RunInstances mocks EC2 RunInstances method.
func (m *mockEC2DryRunMode) RunInstances(input *ec2.RunInstancesInput) (*ec2.Reservation, error) {
	if m.RunInstancesCalls == 0 {
		m.RunInstancesCalls++
		return nil, awserr.New(ErrCodeRequestLimitExceeded, "request limit exceeded", errors.New("test error"))
	}
	if m.RunInstancesCalls == 1 {
		m.RunInstancesCalls++
		return nil, awserr.New(ErrCodeDryRunOperation, "dry run operation", errors.New("dry run error"))
	}
	m.RunInstancesCalls++
	return &ec2.Reservation{}, nil
}
