package ec2

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
