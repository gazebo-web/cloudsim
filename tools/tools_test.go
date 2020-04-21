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