package ec2

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"os"
	"testing"
)

func TestCreateMachines(t *testing.T) {
	suite.Run(t, new(ec2CreateMachinesTestSuite))
}

type ec2CreateMachinesTestSuite struct {
	suite.Suite
	session          *session.Session
	ec2API           ec2iface.EC2API
	machines         cloud.Machines
	machinesCount    int
	subnet           string
	availabilityZone string
	securityGroup    string
	arn              string
}

func (s *ec2CreateMachinesTestSuite) SetupSuite() {
	accessId := os.Getenv("AWS_ACCESS_KEY_ID")
	accessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	endpoint := os.Getenv("AWS_ENDPOINT")
	if len(endpoint) == 0 {
		endpoint = "http://localstack:4566"
	}

	config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(accessId, accessKey, ""),
		Endpoint:         aws.String(endpoint),
		Region:           aws.String(endpoints.UsEast1RegionID),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}
	var err error
	s.session, err = session.NewSession(config)
	if err != nil {
		s.FailNow(err.Error())
	}
	s.ec2API = ec2.New(s.session)
	s.machines = NewMachines(s.ec2API)
	s.subnet, s.availabilityZone, err = s.getDefaultSubnetAndAZ()
	if err != nil {
		s.FailNow(err.Error())
	}
	s.securityGroup, err = s.getDefaultSecurityGroup()
	if err != nil {
		s.FailNow(err.Error())
	}
}

func (s *ec2CreateMachinesTestSuite) SetupTest() {
	var err error
	s.machinesCount, err = s.countMachines()
	s.NoError(err)
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

func (s *ec2CreateMachinesTestSuite) TestCreate_Valid() {
	before := s.machinesCount
	fmt.Println("Subnet:", s.subnet)
	input := []cloud.CreateMachinesInput{
		{
			ResourceName:  nil,
			DryRun:        true,
			KeyName:       "key-name",
			Type:          "a1.medium",
			MinCount:      1,
			MaxCount:      1,
			FirewallRules: []string{s.securityGroup},
			SubnetID:      s.subnet,
			Zone:          s.availabilityZone,
			Tags: map[string]map[string]string{
				"test": {
					"owner": "osrf",
				},
			},
		},
	}
	_, err := s.machines.Create(input)
	s.NoError(err)
	after, err := s.countMachines()
	s.NoError(err)
	s.Equal(before+1, after)
	s.machinesCount = after
}

func (s *ec2CreateMachinesTestSuite) getDefaultSubnetAndAZ() (string, string, error) {
	out, err := s.ec2API.DescribeSubnets(&ec2.DescribeSubnetsInput{})
	if err != nil {
		return "", "", err
	}
	return *out.Subnets[0].SubnetId, *out.Subnets[0].AvailabilityZone, nil
}

func (s *ec2CreateMachinesTestSuite) getDefaultSecurityGroup() (string, error) {
	out, err := s.ec2API.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{})
	if err != nil {
		return "", err
	}
	return *out.SecurityGroups[0].GroupId, nil
}

func (s *ec2CreateMachinesTestSuite) countMachines() (int, error) {
	o, err := s.ec2API.DescribeInstances(&ec2.DescribeInstancesInput{MaxResults: aws.Int64(1000)})
	if err != nil {
		return -1, err
	}
	return len(o.Reservations), nil
}
