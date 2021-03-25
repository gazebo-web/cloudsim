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

func TestCreateMachines(t *testing.T) {
	suite.Run(t, new(ec2CreateMachinesTestSuite))
}

type ec2CreateMachinesTestSuite struct {
	suite.Suite
	ec2API   *mockEC2Create
	machines cloud.Machines
}

func (s *ec2CreateMachinesTestSuite) SetupTest() {
	s.ec2API = &mockEC2Create{}
	logger := ign.NewLoggerNoRollbar("ec2CreateMachinesTestSuite", ign.VerbosityDebug)
	s.machines = NewMachines(s.ec2API, logger)
}

func (s *ec2CreateMachinesTestSuite) TestCreate_MissingKeyName() {
	input := []cloud.CreateMachinesInput{
		{
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

func (s *ec2CreateMachinesTestSuite) TestCreate_InvalidClusterID() {
	input := []cloud.CreateMachinesInput{
		{
			KeyName:       "key-name",
			MinCount:      1,
			MaxCount:      99,
			FirewallRules: nil,
			SubnetID:      "subnet-06fe9fdb790aa78e7",
			Tags:          nil,
			Retries:       0,
			ClusterID:     "",
		},
	}
	_, err := s.machines.Create(input)
	s.Assert().Error(err)
	s.Assert().Equal(cloud.ErrInvalidClusterID, err)
	s.Assert().Equal(0, s.ec2API.RunInstancesCalls)
}

func (s *ec2CreateMachinesTestSuite) TestCreate_InvalidCountBothZero() {
	input := []cloud.CreateMachinesInput{
		{
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
			KeyName:       "key-name",
			MinCount:      1,
			MaxCount:      99,
			FirewallRules: nil,
			SubnetID:      "subnet-06fe9fdb790aa78e7",
			Tags:          nil,
			Retries:       0,
			ClusterID:     "cluster-name",
		},
	}
	_, err := s.machines.Create(input)
	s.NoError(err)
	s.Equal(1, s.ec2API.RunInstancesCalls)
}

func (s *ec2CreateMachinesTestSuite) TestCreate_ValidWithDryRunMode() {
	mock := &mockEC2CreateDryRunMode{}
	logger := ign.NewLoggerNoRollbar("ec2TerminateMachinesTestSuite", ign.VerbosityDebug)
	s.machines = NewMachines(mock, logger)
	input := []cloud.CreateMachinesInput{
		{
			KeyName:       "key-name",
			MinCount:      1,
			MaxCount:      99,
			FirewallRules: nil,
			SubnetID:      "subnet-06fe9fdb790aa78e7",
			Tags:          nil,
			Retries:       3,
			ClusterID:     "cluster-name",
		},
	}
	_, err := s.machines.Create(input)
	s.NoError(err)
	s.Equal(3, mock.RunInstancesCalls)
}

func (s *ec2CreateMachinesTestSuite) TestCreate_ErrorWithDryRunMode() {
	s.ec2API.InternalError = errors.New("force error on dry run mode")
	input := []cloud.CreateMachinesInput{
		{
			KeyName:       "key-name",
			MinCount:      1,
			MaxCount:      99,
			FirewallRules: nil,
			SubnetID:      "subnet-06fe9fdb790aa78e7",
			Tags:          nil,
			Retries:       3,
			ClusterID:     "cluster-name",
		},
	}
	_, err := s.machines.Create(input)
	s.Require().Error(err)
	s.Assert().True(errors.Is(err, cloud.ErrUnknown))
}

type mockEC2Create struct {
	ec2iface.EC2API
	RunInstancesCalls int
	InternalError     error
}

// RunInstances mocks EC2 RunInstances method.
func (m *mockEC2Create) RunInstances(input *ec2.RunInstancesInput) (*ec2.Reservation, error) {
	if m.InternalError != nil {
		return nil, m.InternalError
	}
	m.RunInstancesCalls++
	return &ec2.Reservation{}, nil
}

type mockEC2CreateDryRunMode struct {
	ec2iface.EC2API
	RunInstancesCalls int
}

// RunInstances mocks EC2 RunInstances method.
func (m *mockEC2CreateDryRunMode) RunInstances(input *ec2.RunInstancesInput) (*ec2.Reservation, error) {
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
