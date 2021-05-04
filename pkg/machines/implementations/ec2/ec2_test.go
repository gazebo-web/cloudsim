package ec2

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cycler"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/machines"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"testing"
	"time"
)

func TestEC2MachinesSuite(t *testing.T) {
	suite.Run(t, &EC2MachinesTestSuite{})
}

type EC2MachinesTestSuite struct {
	suite.Suite
	zones      []Zone
	zoneCycler cycler.Cycler
	m          *ec2Machines
}

func (s *EC2MachinesTestSuite) SetupTest() {
	s.zones = []Zone{
		{
			Zone:     "zone-0",
			SubnetID: "subnet-0",
		},
		{
			Zone:     "zone-1",
			SubnetID: "subnet-1",
		},
	}

	var err error
	s.zoneCycler, err = cycler.NewCyclerFromSlice(s.zones)
	s.Require().NoError(err)

	s.m = &ec2Machines{
		zones: s.zoneCycler,
	}
}

func (s *EC2MachinesTestSuite) TestNewMachines() {
	session, err := session.NewSession(nil)
	s.Require().NoError(err)
	ec := ec2.New(session)
	logger := ign.NewLoggerNoRollbar("TestNewMachines", ign.VerbosityDebug)
	m, err := NewMachines(&NewInput{
		API:    ec,
		Logger: logger,
		Zones:  s.zones,
	})
	s.Require().NoError(err)
	e, ok := m.(*ec2Machines)
	s.Require().True(ok)
	s.Assert().NotNil(e.API)
}

func (s *EC2MachinesTestSuite) TestIsValidKeyName() {
	s.Assert().False(s.m.isValidKeyName(""))
	s.Assert().True(s.m.isValidKeyName("testKey"))
}

func (s *EC2MachinesTestSuite) TestIsValidMachineCount() {
	s.Assert().False(s.m.isValidMachineCount(-1, -1))
	s.Assert().False(s.m.isValidMachineCount(-1, 0))
	s.Assert().False(s.m.isValidMachineCount(0, -1))
	s.Assert().False(s.m.isValidMachineCount(0, 0))
	s.Assert().False(s.m.isValidMachineCount(1, 0))
	s.Assert().False(s.m.isValidMachineCount(0, 1))
	s.Assert().False(s.m.isValidMachineCount(20, 1))
	s.Assert().True(s.m.isValidMachineCount(1, 1))
	s.Assert().True(s.m.isValidMachineCount(1, 5))
}

func (s *EC2MachinesTestSuite) TestIsValidSubnetID() {
	s.Assert().False(s.m.isValidSubnetID(aws.String("")))
	s.Assert().False(s.m.isValidSubnetID(aws.String("test")))
	s.Assert().False(s.m.isValidSubnetID(aws.String("test-1234")))
	s.Assert().False(s.m.isValidSubnetID(aws.String("tested-0dae7657")))
	s.Assert().False(s.m.isValidSubnetID(aws.String("tested-06fe9fdb790aa78e7")))

	s.Assert().False(s.m.isValidSubnetID(aws.String("subnet-0de7657")))
	s.Assert().True(s.m.isValidSubnetID(aws.String("subnet-0dae7657")))
	s.Assert().False(s.m.isValidSubnetID(aws.String("subnet-0dae765712")))

	s.Assert().False(s.m.isValidSubnetID(aws.String("subnet-06fe9fdb790aa7")))
	s.Assert().True(s.m.isValidSubnetID(aws.String("subnet-06fe9fdb790aa78e7")))
	s.Assert().False(s.m.isValidSubnetID(aws.String("subnet-06fe9fdb790aa78e71234778")))
}

func (s *EC2MachinesTestSuite) TestSleep0Seconds() {
	before := time.Now()
	s.m.sleepNSecondsBeforeMaxRetries(0, 10)
	now := time.Now()
	s.Assert().Equal(before.Second(), now.Second())
}

func (s *EC2MachinesTestSuite) TestSleep1Seconds() {
	before := time.Now()
	s.m.sleepNSecondsBeforeMaxRetries(1, 10)
	now := time.Now()
	s.Assert().Equal(before.Second()+1, now.Second())
}

func (s *EC2MachinesTestSuite) TestSleep0SecondsWhenIsMax() {
	before := time.Now()
	s.m.sleepNSecondsBeforeMaxRetries(10, 10)
	now := time.Now()
	s.Assert().Equal(before.Second(), now.Second())
}

func (s *EC2MachinesTestSuite) TestNewRunInstanceInput() {
	instanceProfile := "arn"
	script := "bash"
	out := s.m.newRunInstancesInput(machines.CreateMachinesInput{
		InstanceProfile: &instanceProfile,
		KeyName:         "key-name",
		Type:            "t2.large",
		Image:           "docker-image",
		MinCount:        1,
		MaxCount:        2,
		FirewallRules:   []string{"first-rule", "second-rule"},
		SubnetID:        &s.zones[1].SubnetID,
		Zone:            &s.zones[1].Zone,
		Tags: []machines.Tag{
			{
				Resource: "instance",
				Map: map[string]string{
					"key": "value",
				},
			},
		},
		InitScript: &script,
		Retries:    1,
	})

	s.Assert().NotNil(out.IamInstanceProfile.Arn)
	s.Assert().Equal("arn", *out.IamInstanceProfile.Arn)

	s.Assert().NotNil(out.KeyName)
	s.Assert().Equal("key-name", *out.KeyName)

	s.Assert().NotNil(out.InstanceType)
	s.Assert().Equal("t2.large", *out.InstanceType)

	s.Assert().NotNil(out.MinCount)
	s.Assert().NotNil(out.MaxCount)
	s.Assert().Equal(int64(1), *out.MinCount)
	s.Assert().Equal(int64(2), *out.MaxCount)

	s.Assert().NotNil(out.SubnetId)
	s.Assert().Equal(s.zones[1].SubnetID, *out.SubnetId)

	s.Assert().NotNil(out.Placement.AvailabilityZone)
	s.Assert().Equal(s.zones[1].Zone, *out.Placement.AvailabilityZone)

	s.Assert().NotNil(out.UserData)
	s.Assert().Equal("bash", *out.UserData)
}

func (s *EC2MachinesTestSuite) TestCreateTags() {
	tags := []machines.Tag{
		{
			Resource: "instance",
			Map: map[string]string{
				"key": "value",
			},
		},
	}

	tagSpec := s.m.createTags(tags)

	s.Assert().Len(tagSpec, 1)
	s.Assert().Len(tagSpec[0].Tags, 1)
	s.Assert().Equal("key", *tagSpec[0].Tags[0].Key)
	s.Assert().Equal("value", *tagSpec[0].Tags[0].Value)
}

func (s *EC2MachinesTestSuite) TestCreateFilters() {
	output := s.m.createFilters(map[string][]string{
		"instance-state-name": {
			"pending",
			"running",
		},
	})

	s.Assert().NotNil(output[0].Name)
	s.Assert().Equal("instance-state-name", *output[0].Name)
	s.Assert().NotNil(output[0].Values[0])
	s.Assert().Equal("pending", *output[0].Values[0])
	s.Assert().NotNil(output[0].Values[1])
	s.Assert().Equal("running", *output[0].Values[1])
}

func (s *EC2MachinesTestSuite) TestParseRunInstanceError() {
	err := s.m.parseRunInstanceError(errors.New("internal error"))
	s.Assert().True(errors.Is(err, machines.ErrUnknown))

	err = s.m.parseRunInstanceError(awserr.New(ErrCodeInsufficientInstanceCapacity, "test", nil))
	s.Assert().Equal(machines.ErrInsufficientMachines, err)

	err = s.m.parseRunInstanceError(awserr.New(ErrCodeRequestLimitExceeded, "test", nil))
	s.Assert().Equal(machines.ErrRequestsLimitExceeded, err)
}

func (s *EC2MachinesTestSuite) TestMachines_checkAvailableMachines() {
	// If limit is set to -1, always return true.
	m := &ec2Machines{
		limit: -1,
		zones: s.zoneCycler,
	}
	inputs := []machines.CreateMachinesInput{{MaxCount: 1}}
	s.Assert().True(m.checkAvailableMachines(inputs))

	// If limit is set, should return true if there are enough machines available.
	mockCounter := &mockEC2Count{
		ReturnMachines: true, // Returns 1 machines
	}
	m = &ec2Machines{
		limit:  2,
		zones:  s.zoneCycler,
		API:    mockCounter,
		Logger: ign.NewLoggerNoRollbar("TestMachines_checkAvailableMachines", ign.VerbosityDebug),
	}
	s.Assert().True(m.checkAvailableMachines(inputs))

	// If limit is set to the total amount of machines created at a certain moment, it should return false.
	m = &ec2Machines{
		limit:  1,
		API:    mockCounter,
		Logger: ign.NewLoggerNoRollbar("TestMachines_checkAvailableMachines", ign.VerbosityDebug),
	}
	s.Assert().False(m.checkAvailableMachines(inputs))
}
