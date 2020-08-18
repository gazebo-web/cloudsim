package ec2

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"testing"
	"time"
)

func TestNewMachines(t *testing.T) {
	s, err := session.NewSession(nil)
	assert.NoError(t, err)
	ec := ec2.New(s)
	logger := ign.NewLoggerNoRollbar("TestNewMachines", ign.VerbosityDebug)
	m := NewMachines(ec, logger)
	e, ok := m.(*machines)
	assert.True(t, ok)
	assert.NotNil(t, e.API)
}

func TestIsValidKeyName(t *testing.T) {
	m := &machines{}
	assert.False(t, m.isValidKeyName(""))
	assert.True(t, m.isValidKeyName("testKey"))
}

func TestIsValidMachineCount(t *testing.T) {
	m := &machines{}
	assert.False(t, m.isValidMachineCount(-1, -1))
	assert.False(t, m.isValidMachineCount(-1, 0))
	assert.False(t, m.isValidMachineCount(0, -1))
	assert.False(t, m.isValidMachineCount(0, 0))
	assert.False(t, m.isValidMachineCount(1, 0))
	assert.False(t, m.isValidMachineCount(0, 1))
	assert.False(t, m.isValidMachineCount(20, 1))
	assert.True(t, m.isValidMachineCount(1, 1))
	assert.True(t, m.isValidMachineCount(1, 5))
}

func TestIsValidSubnetID(t *testing.T) {
	m := &machines{}
	assert.False(t, m.isValidSubnetID(""))
	assert.False(t, m.isValidSubnetID("test"))
	assert.False(t, m.isValidSubnetID("test-1234"))
	assert.False(t, m.isValidSubnetID("tested-0dae7657"))
	assert.False(t, m.isValidSubnetID("tested-06fe9fdb790aa78e7"))

	assert.False(t, m.isValidSubnetID("subnet-0de7657"))
	assert.True(t, m.isValidSubnetID("subnet-0dae7657"))
	assert.False(t, m.isValidSubnetID("subnet-0dae765712"))

	assert.False(t, m.isValidSubnetID("subnet-06fe9fdb790aa7"))
	assert.True(t, m.isValidSubnetID("subnet-06fe9fdb790aa78e7"))
	assert.False(t, m.isValidSubnetID("subnet-06fe9fdb790aa78e71234778"))
}

func TestSleep0Seconds(t *testing.T) {
	m := &machines{}
	before := time.Now()
	m.sleepNSecondsBeforeMaxRetries(0, 10)
	now := time.Now()
	assert.Equal(t, before.Second(), now.Second())
}

func TestSleep1Seconds(t *testing.T) {
	m := &machines{}
	before := time.Now()
	m.sleepNSecondsBeforeMaxRetries(1, 10)
	now := time.Now()
	assert.Equal(t, before.Second()+1, now.Second())
}

func TestSleep0SecondsWhenIsMax(t *testing.T) {
	m := &machines{}
	before := time.Now()
	m.sleepNSecondsBeforeMaxRetries(10, 10)
	now := time.Now()
	assert.Equal(t, before.Second(), now.Second())
}

func TestNewRunInstanceInput(t *testing.T) {
	m := &machines{}
	instanceProfile := "arn"
	script := "bash"
	out := m.newRunInstancesInput(cloud.CreateMachinesInput{
		InstanceProfile: &instanceProfile,
		KeyName:         "key-name",
		Type:            "t2.large",
		Image:           "docker-image",
		MinCount:        1,
		MaxCount:        2,
		FirewallRules:   []string{"first-rule", "second-rule"},
		SubnetID:        "subnet-id",
		Zone:            "zone-a",
		Tags: []cloud.Tag{
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

	assert.NotNil(t, out.IamInstanceProfile.Arn)
	assert.Equal(t, "arn", *out.IamInstanceProfile.Arn)

	assert.NotNil(t, out.KeyName)
	assert.Equal(t, "key-name", *out.KeyName)

	assert.NotNil(t, out.InstanceType)
	assert.Equal(t, "t2.large", *out.InstanceType)

	assert.NotNil(t, out.MinCount)
	assert.NotNil(t, out.MaxCount)
	assert.Equal(t, int64(1), *out.MinCount)
	assert.Equal(t, int64(2), *out.MaxCount)

	assert.NotNil(t, out.SecurityGroups)
	assert.Len(t, out.SecurityGroups, 2)
	assert.Equal(t, "first-rule", *out.SecurityGroups[0])
	assert.Equal(t, "second-rule", *out.SecurityGroups[1])

	assert.NotNil(t, out.SubnetId)
	assert.Equal(t, "subnet-id", *out.SubnetId)

	assert.NotNil(t, out.Placement.AvailabilityZone)
	assert.Equal(t, "zone-a", *out.Placement.AvailabilityZone)

	assert.NotNil(t, out.UserData)
	assert.Equal(t, "bash", *out.UserData)
}

func TestCreateTags(t *testing.T) {
	m := &machines{}

	tags := []cloud.Tag{
		{
			Resource: "instance",
			Map: map[string]string{
				"key": "value",
			},
		},
	}

	tagSpec := m.createTags(tags)

	assert.Len(t, tagSpec, 1)
	assert.Len(t, tagSpec[0].Tags, 1)
	assert.Equal(t, "key", *tagSpec[0].Tags[0].Key)
	assert.Equal(t, "value", *tagSpec[0].Tags[0].Value)
}

func TestCreateFilters(t *testing.T) {
	m := &machines{}

	output := m.createFilters(map[string][]string{
		"instance-state-name": {
			"pending",
			"running",
		},
	})

	assert.NotNil(t, output[0].Name)
	assert.Equal(t, "instance-state-name", *output[0].Name)
	assert.NotNil(t, output[0].Values[0])
	assert.Equal(t, "pending", *output[0].Values[0])
	assert.NotNil(t, output[0].Values[1])
	assert.Equal(t, "running", *output[0].Values[1])
}

func TestParseRunInstanceError(t *testing.T) {
	m := &machines{}

	err := m.parseRunInstanceError(errors.New("internal error"))
	assert.Equal(t, cloud.ErrUnknown, err)

	err = m.parseRunInstanceError(awserr.New(ErrCodeInsufficientInstanceCapacity, "test", nil))
	assert.Equal(t, cloud.ErrInsufficientMachines, err)

	err = m.parseRunInstanceError(awserr.New(ErrCodeRequestLimitExceeded, "test", nil))
	assert.Equal(t, cloud.ErrRequestsLimitExceeded, err)
}
