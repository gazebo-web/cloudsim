package waiter

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestNewWaitRequest(t *testing.T) {
	job := func() (bool, error) {
		fmt.Println("test")
		return true, nil
	}
	wr := NewWaitRequest(job)
	assert.NotNil(t, wr)
	assert.IsType(t, &request{}, wr)
}

func TestRequest_WaitFailsWhenErrIsReturned(t *testing.T) {
	returnedErr := errors.New("test error")
	job := func() (bool, error) {
		return true, returnedErr
	}
	wr := NewWaitRequest(job)

	var wg sync.WaitGroup
	var err error
	wg.Add(1)
	go func() {
		err = wr.Wait(3*time.Second, 1*time.Microsecond)
		wg.Done()
	}()
	wg.Wait()
	assert.Error(t, err)
	assert.Equal(t, returnedErr, err)
}

func TestRequest_WaitTimeouts(t *testing.T) {
	job := func() (bool, error) {
		return false, nil
	}
	wr := NewWaitRequest(job)

	var wg sync.WaitGroup
	var err error
	wg.Add(1)
	go func() {
		err = wr.Wait(3*time.Second, 1*time.Microsecond)
		wg.Done()
	}()
	wg.Wait()
	assert.Error(t, err)
	assert.Equal(t, ErrRequestTimeout, err)
}
