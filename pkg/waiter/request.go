package waiter

import (
	"k8s.io/apimachinery/pkg/util/wait"
	"time"
)

var (
	// ErrRequestTimeout is an error returned when the Wait method timeouts.
	ErrRequestTimeout = wait.ErrWaitTimeout
)

// Waiter is used to wait for kubernetes nodes and pods to be in a certain state.
type Waiter interface {
	Wait(timeout time.Duration, frequency time.Duration) error
}

// request is a Waiter implementation that will be used to wait for a job to succeed.
type request struct {
	job func() (bool, error)
}

// Wait executes a job in regular time intervals given by a certain frequency.
// If there is an error in the job or the
func (r request) Wait(timeout time.Duration, frequency time.Duration) error {
	return wait.PollImmediate(frequency, timeout, r.job)
}

// NewWaitRequest creates a new Waiter implementation.
func NewWaitRequest(job func() (bool, error)) Waiter {
	return &request{
		job: job,
	}
}
