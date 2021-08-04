package ec2

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/machines"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"testing"
)

func TestCreateMachines(t *testing.T) {
	suite.Run(t, new(ec2CreateMachinesTestSuite))
}

type ec2CreateMachinesTestSuite struct {
	suite.Suite
	ec2API   *mockEC2Create
	machines machines.Machines
}

func (s *ec2CreateMachinesTestSuite) SetupTest() {
	s.ec2API = &mockEC2Create{}
	logger := ign.NewLoggerNoRollbar("ec2CreateMachinesTestSuite", ign.VerbosityDebug)
	var err error
	s.machines, err = NewMachines(&NewInput{
		API:    s.ec2API,
		Logger: logger,
		Zones: []Zone{
			{
				Zone:     "test1",
				SubnetID: "test1",
			},
			{
				Zone:     "test2",
				SubnetID: "test2",
			},
			{
				Zone:     "test3",
				SubnetID: "test3",
			},
		},
	})
	s.Require().NoError(err)
}

func (s *ec2CreateMachinesTestSuite) TestCreate_MissingKeyName() {
	input := []machines.CreateMachinesInput{
		{
			KeyName:       "",
			MinCount:      1,
			MaxCount:      10,
			FirewallRules: nil,
			SubnetID:      aws.String("subnet-06fe9fdb790aa78e7"),
			Tags:          nil,
		},
	}
	_, err := s.machines.Create(input)
	s.Error(err)
	s.Equal(machines.ErrMissingKeyName, err)
	s.Equal(0, s.ec2API.RunInstancesCalls)
}

func (s *ec2CreateMachinesTestSuite) TestCreate_InvalidClusterID() {
	input := []machines.CreateMachinesInput{
		{
			KeyName:       "key-name",
			MinCount:      1,
			MaxCount:      99,
			FirewallRules: nil,
			SubnetID:      aws.String("subnet-06fe9fdb790aa78e7"),
			Tags:          nil,
			Retries:       0,
			ClusterID:     "",
		},
	}
	_, err := s.machines.Create(input)
	s.Assert().Error(err)
	s.Assert().Equal(machines.ErrInvalidClusterID, err)
	s.Assert().Equal(0, s.ec2API.RunInstancesCalls)
}

func (s *ec2CreateMachinesTestSuite) TestCreate_InvalidCountBothZero() {
	input := []machines.CreateMachinesInput{
		{
			KeyName:       "key-name",
			MinCount:      0,
			MaxCount:      0,
			FirewallRules: nil,
			SubnetID:      aws.String("subnet-06fe9fdb790aa78e7"),
			Tags:          nil,
		},
	}
	_, err := s.machines.Create(input)
	s.Error(err)
	s.Equal(machines.ErrInvalidMachinesCount, err)
	s.Equal(0, s.ec2API.RunInstancesCalls)
}

func (s *ec2CreateMachinesTestSuite) TestCreate_InvalidCountMinCountZero() {
	input := []machines.CreateMachinesInput{
		{
			KeyName:       "key-name",
			MinCount:      0,
			MaxCount:      1,
			FirewallRules: nil,
			SubnetID:      aws.String("subnet-06fe9fdb790aa78e7"),
			Tags:          nil,
		},
	}
	_, err := s.machines.Create(input)
	s.Error(err)
	s.Equal(machines.ErrInvalidMachinesCount, err)
	s.Equal(0, s.ec2API.RunInstancesCalls)
}

func (s *ec2CreateMachinesTestSuite) TestCreate_InvalidCountMaxCountZero() {
	input := []machines.CreateMachinesInput{
		{
			KeyName:       "key-name",
			MinCount:      1,
			MaxCount:      0,
			FirewallRules: nil,
			SubnetID:      aws.String("subnet-06fe9fdb790aa78e7"),
			Tags:          nil,
		},
	}
	_, err := s.machines.Create(input)
	s.Error(err)
	s.Equal(machines.ErrInvalidMachinesCount, err)
	s.Equal(0, s.ec2API.RunInstancesCalls)
}

func (s *ec2CreateMachinesTestSuite) TestCreate_MinCountGreaterThanMaxCount() {
	input := []machines.CreateMachinesInput{
		{
			KeyName:       "key-name",
			MinCount:      99,
			MaxCount:      1,
			FirewallRules: nil,
			SubnetID:      aws.String("subnet-06fe9fdb790aa78e7"),
			Tags:          nil,
		},
	}
	_, err := s.machines.Create(input)
	s.Error(err)
	s.Equal(machines.ErrInvalidMachinesCount, err)
	s.Equal(0, s.ec2API.RunInstancesCalls)
}

func (s *ec2CreateMachinesTestSuite) TestCreate_NegativeCount() {
	input := []machines.CreateMachinesInput{
		{
			KeyName:       "key-name",
			MinCount:      -100,
			MaxCount:      -25,
			FirewallRules: nil,
			SubnetID:      aws.String("subnet-06fe9fdb790aa78e7"),
			Tags:          nil,
		},
	}
	_, err := s.machines.Create(input)
	s.Error(err)
	s.Equal(machines.ErrInvalidMachinesCount, err)
	s.Equal(0, s.ec2API.RunInstancesCalls)
}

func (s *ec2CreateMachinesTestSuite) TestCreate_InvalidSubnet() {
	input := []machines.CreateMachinesInput{
		{
			KeyName:       "key-name",
			MinCount:      1,
			MaxCount:      99,
			FirewallRules: nil,
			SubnetID:      aws.String("subnet-1234"),
			Tags:          nil,
		},
	}
	_, err := s.machines.Create(input)
	s.Error(err)
	s.Equal(machines.ErrInvalidSubnetID, err)

	s.Equal(0, s.ec2API.RunInstancesCalls)
}

func (s *ec2CreateMachinesTestSuite) TestCreate_ValidWithoutDryRunMode() {
	input := []machines.CreateMachinesInput{
		{
			KeyName:       "key-name",
			MinCount:      1,
			MaxCount:      99,
			FirewallRules: nil,
			SubnetID:      aws.String("subnet-06fe9fdb790aa78e7"),
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
	var err error
	s.machines, err = NewMachines(&NewInput{
		API:    mock,
		Logger: logger,
		Zones: []Zone{
			{
				Zone:     "test",
				SubnetID: "test",
			},
		},
	})
	s.Require().NoError(err)
	input := []machines.CreateMachinesInput{
		{
			KeyName:       "key-name",
			MinCount:      1,
			MaxCount:      99,
			FirewallRules: nil,
			SubnetID:      aws.String("subnet-06fe9fdb790aa78e7"),
			Tags:          nil,
			Retries:       3,
			ClusterID:     "cluster-name",
		},
	}
	_, err = s.machines.Create(input)
	s.NoError(err)
	s.Equal(3, mock.RunInstancesCalls)
}

func (s *ec2CreateMachinesTestSuite) TestCreate_ErrorWithDryRunMode() {
	s.ec2API.InternalError = errors.New("force error on dry run mode")
	input := []machines.CreateMachinesInput{
		{
			KeyName:       "key-name",
			MinCount:      1,
			MaxCount:      99,
			FirewallRules: nil,
			SubnetID:      aws.String("subnet-06fe9fdb790aa78e7"),
			Tags:          nil,
			Retries:       3,
			ClusterID:     "cluster-name",
		},
	}
	_, err := s.machines.Create(input)
	s.Require().Error(err)
	s.Assert().True(errors.Is(err, machines.ErrUnknown))
}

func (s *ec2CreateMachinesTestSuite) TestCreate_RotateAvailabilityZones() {
	before := s.machines.(*ec2Machines).zones.Get().(Zone)

	output, err := s.machines.(*ec2Machines).create(machines.CreateMachinesInput{
		KeyName:       "key-name",
		MinCount:      1,
		MaxCount:      99,
		FirewallRules: nil,
		SubnetID:      aws.String("subnet-06fe9fdb790aa78e7"),
		Tags:          nil,
		Retries:       0,
		ClusterID:     "cluster-name",
	})
	s.Require().NoError(err)

	// If everything went well, the first zone should have been used first.
	s.Assert().Equal(before.Zone, output.Zone)

	// After the request, a rotation step is performed, making zone's cycler start on the next available zone.
	after := s.machines.(*ec2Machines).zones.Get().(Zone)
	s.Assert().NotEqual(before.Zone, after.Zone)

	// Repeat call to make sure the next zone is being used
	output, err = s.machines.(*ec2Machines).create(machines.CreateMachinesInput{
		KeyName:       "key-name",
		MinCount:      1,
		MaxCount:      99,
		FirewallRules: nil,
		SubnetID:      aws.String("subnet-06fe9fdb790aa78e7"),
		Tags:          nil,
		Retries:       0,
		ClusterID:     "cluster-name",
	})
	s.Require().NoError(err)

	s.Assert().Equal(after.Zone, output.Zone)
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
