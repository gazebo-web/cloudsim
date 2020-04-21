package tools

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestSptr(t *testing.T) {
	var text string
	text = "Test"
	result := Sptr(text)
	var resultType *string
	assert.IsType(t, resultType, result)
	assert.NotNil(t, result)
	assert.Equal(t, "Test", *result)
}

func TestSleep(t *testing.T) {
	var waitTime time.Duration
	var wg sync.WaitGroup

	waitTime = 5 * time.Millisecond


	slept := false

	wg.Add(1)

	go func() {
		Sleep(waitTime)
		slept = true
		wg.Done()
	}()

	wg.Wait()

	assert.Equal(t, true, slept)
}

func TestIsECR(t *testing.T) {
	assert.False(t, IsECR("not-ecr"))
	assert.True(t, IsECR("1111111111.dkr.ecr.us-west-1.amazonaws.com/osrf:test"))
}

func TestGetLocalIPAddress(t *testing.T) {
	result, err := GetLocalIPAddress()
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Contains(t, result, ".")
	assert.Regexp(t, "^(?:[0-9]{1,3}\\.){3}[0-9]{1,3}$", result)
}

func TestIntptr(t *testing.T) {
	var value int
	value = 5
	result := Intptr(value)
	assert.NotNil(t, result)
	assert.IsType(t, &value, result)
	assert.Equal(t, value, *result)
}

func TestInt32ptr(t *testing.T) {
	var value int32
	value = 5
	result := Int32ptr(value)
	assert.NotNil(t, result)
	assert.IsType(t, &value, result)
	assert.Equal(t, value, *result)
}

func TestInt64ptr(t *testing.T) {
	var value int64
	value = 5
	result := Int64ptr(value)
	assert.NotNil(t, result)
	assert.IsType(t, &value, result)
	assert.Equal(t, value, *result)
}