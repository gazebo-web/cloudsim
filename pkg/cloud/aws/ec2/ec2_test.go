package ec2

import (
	"github.com/stretchr/testify/assert"
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
