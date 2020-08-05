package ec2

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
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
	session  *session.Session
	ec2API   ec2iface.EC2API
	machines cloud.Machines
}

func (s *ec2CreateMachinesTestSuite) SetupTest() {
	var err error
	s.NoError(err)
	config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials("YOUR-ACCESSKEYID", "YOUR-SECRETACCESSKEY", ""),
		Endpoint:         aws.String("https://localstack:4566"),
		Region:           aws.String(endpoints.UsEast1RegionID),
		S3ForcePathStyle: aws.Bool(true),
	}
	s.session = session.Must(session.NewSession(config))
	s.ec2API = ec2.New(s.session)
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
}

func (s *ec2CreateMachinesTestSuite) TestCreate_MinCountEqualsMaxCount() {
	input := []cloud.CreateMachinesInput{
		{
			DryRun:        false,
			KeyName:       "key-name",
			MinCount:      1,
			MaxCount:      1,
			FirewallRules: nil,
			SubnetID:      "subnet-06fe9fdb790aa78e7",
			Tags:          nil,
		},
	}
	_, err := s.machines.Create(input)
	s.NoError(err)
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
}

func (s *ec2CreateMachinesTestSuite) TestCreate_InvalidSubnet() {
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

	input = []cloud.CreateMachinesInput{
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
	_, err = s.machines.Create(input)
	s.Error(err)
	s.Equal(cloud.ErrInvalidSubnetID, err)

	input = []cloud.CreateMachinesInput{
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
	_, err = s.machines.Create(input)
	s.Error(err)
	s.Equal(cloud.ErrInvalidSubnetID, err)
}

func (s *ec2CreateMachinesTestSuite) TestCreate_ValidSubnet() {
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
