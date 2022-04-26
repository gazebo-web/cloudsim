package ec2

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/stretchr/testify/suite"
	cloud "gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud/aws"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/machines"
	"gitlab.com/ignitionrobotics/web/ign-go/v5"
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
		API:            s.ec2API,
		Logger:         logger,
		CostCalculator: cloud.NewCostCalculatorEC2(nil),
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
		API:            mock,
		CostCalculator: cloud.NewCostCalculatorEC2(nil),
		Logger:         logger,
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

	s.ec2API.ResetInstanceCallsToZero = true

	for i := 0; i < 5; i++ {
		_, err := s.machines.(*ec2Machines).create(machines.CreateMachinesInput{
			KeyName:       "key-name",
			MinCount:      1,
			MaxCount:      99,
			FirewallRules: nil,
			Tags:          nil,
			Retries:       0,
			ClusterID:     "cluster-name",
		})
		s.Require().NoError(err)
	}

	// Request a new machine on zone 1, it will rotate to another zone.
	_, _ = s.machines.(*ec2Machines).create(machines.CreateMachinesInput{
		KeyName:       "key-name",
		MinCount:      1,
		MaxCount:      99,
		FirewallRules: nil,
		Tags:          nil,
		Retries:       0,
		ClusterID:     "cluster-name",
	})

	// After an error occurred with the first zone, a rotation step is performed, making zone's cycler
	// start on the next available zone.
	after := s.machines.(*ec2Machines).zones.Get().(Zone)
	s.Assert().NotEqual(before.Zone, after.Zone)
}

func (s *ec2CreateMachinesTestSuite) TestCreate_RotateAvailabilityDepletedZones() {
	before := s.machines.(*ec2Machines).zones.Get().(Zone)

	// Request 10 machines on zone 1, this will cause all zones to be depleted (on the mock implementation).
	s.ec2API.ResetInstanceCallsToZero = false

	for i := 0; i < 5; i++ {
		_, err := s.machines.(*ec2Machines).create(machines.CreateMachinesInput{
			KeyName:       "key-name",
			MinCount:      1,
			MaxCount:      99,
			FirewallRules: nil,
			Tags:          nil,
			Retries:       0,
			ClusterID:     "cluster-name",
		})
		s.Require().NoError(err)
	}

	// After all zones are depleted, it should return an error.
	_, err := s.machines.(*ec2Machines).create(machines.CreateMachinesInput{
		KeyName:       "key-name",
		MinCount:      1,
		MaxCount:      99,
		FirewallRules: nil,
		Tags:          nil,
		Retries:       0,
		ClusterID:     "cluster-name",
	})
	s.Assert().Error(err)

	// After all zones returned an error, the cycler returns to the first zone.
	after := s.machines.(*ec2Machines).zones.Get().(Zone)
	s.Assert().Equal(before.Zone, after.Zone)
	s.Assert().Equal("test1", after.Zone)
}

type mockEC2Create struct {
	ec2iface.EC2API
	RunInstancesCalls        int
	InternalError            error
	ResetInstanceCallsToZero bool
}

// RunInstances mocks EC2 RunInstances method.
func (m *mockEC2Create) RunInstances(input *ec2.RunInstancesInput) (*ec2.Reservation, error) {
	if m.InternalError != nil {
		return nil, m.InternalError
	}
	if m.RunInstancesCalls >= 5 {
		if m.ResetInstanceCallsToZero {
			m.RunInstancesCalls = 0
		}
		return nil, awserr.New(ErrCodeInsufficientInstanceCapacity, "error instance capacity", errors.New("test error"))
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
