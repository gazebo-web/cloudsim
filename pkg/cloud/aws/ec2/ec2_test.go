package ec2

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"sync"
	"testing"
	"time"
)

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

	assert.False(t, m.isValidSubnetID("subnet-0de7657"))
	assert.True(t, m.isValidSubnetID("subnet-0dae7657"))
	assert.False(t, m.isValidSubnetID("subnet-0dae765712"))

	assert.False(t, m.isValidSubnetID("subnet-06fe9fdb790aa7"))
	assert.True(t, m.isValidSubnetID("subnet-06fe9fdb790aa78e7"))
	assert.False(t, m.isValidSubnetID("subnet-06fe9fdb790aa78e71234778"))
}

func TestSleep0Seconds(t *testing.T) {
	m := &machines{}

	var wg sync.WaitGroup
	wg.Add(1)
	before := time.Now()
	func() {
		m.sleepNSecondsBeforeMaxRetries(0, 10)
		wg.Done()
	}()
	wg.Wait()
	now := time.Now()
	assert.Equal(t, before.Second(), now.Second())
}

func TestSleep1Seconds(t *testing.T) {
	m := &machines{}

	var wg sync.WaitGroup
	wg.Add(1)
	before := time.Now()
	func() {
		m.sleepNSecondsBeforeMaxRetries(1, 10)
		wg.Done()
	}()
	wg.Wait()
	now := time.Now()
	assert.Equal(t, before.Second()+1, now.Second())
}

func TestSleep0SecondsWhenIsMax(t *testing.T) {
	m := &machines{}

	var wg sync.WaitGroup
	wg.Add(1)
	before := time.Now()
	func() {
		m.sleepNSecondsBeforeMaxRetries(10, 10)
		wg.Done()
	}()
	wg.Wait()
	now := time.Now()
	assert.Equal(t, before.Second(), now.Second())
}

func TestNewRunInstanceInput(t *testing.T) {
	m := &machines{}
	out := m.newRunInstancesInput(cloud.CreateMachinesInput{
		ResourceName:  "arn",
		DryRun:        true,
		KeyName:       "key-name",
		Type:          "t2.large",
		Image:         "docker-image",
		MinCount:      1,
		MaxCount:      2,
		FirewallRules: []string{"first-rule", "second-rule"},
		SubnetID:      "subnet-id",
		Zone:          "zone-a",
		Tags: map[string]map[string]string{
			"namespace": {
				"key": "value",
			},
		},
		InitScript: "bash",
		Retries:    1,
	})

	assert.NotNil(t, out.IamInstanceProfile.Arn)
	assert.Equal(t, "arn", *out.IamInstanceProfile.Arn)

	assert.NotNil(t, *out.DryRun)
	assert.True(t, *out.DryRun)

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

	tagSpec := m.createTags(map[string]map[string]string{
		"namespace": {
			"key": "value",
		},
	})

	assert.Len(t, tagSpec, 1)
	assert.Len(t, tagSpec[0].Tags, 1)
	assert.Equal(t, "key", *tagSpec[0].Tags[0].Key)
	assert.Equal(t, "value", *tagSpec[0].Tags[0].Value)
}
