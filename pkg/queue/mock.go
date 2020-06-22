package queue

import (
	"github.com/stretchr/testify/mock"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

type Mock struct {
	mock.Mock
}

// Get returns the entire list of items from the queue.
// If `offset` and `limit` are not nil, it will return up to `limit` results from the provided `offset`.
func (m *Mock) Get(offset, limit *int) ([]interface{}, *ign.ErrMsg) {
	args := m.Called(offset, limit)

	var result []interface{}
	result = append(result, args.Get(0))

	var err *ign.ErrMsg
	err = args.Get(1).(*ign.ErrMsg)

	return result, err
}

// Enqueue enqueues a groupID on the queue.
// Returns the groupID that was pushed.
func (m *Mock) Enqueue(entity interface{}) interface{} {
	args := m.Called(entity)
	return args.Get(0)
}

// Dequeue returns the next groupID from the queue.
func (m *Mock) Dequeue() (interface{}, *ign.ErrMsg) {
	args := m.Called()

	var err *ign.ErrMsg
	err = args.Get(1).(*ign.ErrMsg)

	return args.Get(0), err
}

// DequeueOrWait returns the next groupID from the queue or waits until there is one available.
func (m *Mock) DequeueOrWait() (interface{}, *ign.ErrMsg) {
	args := m.Called()

	var err *ign.ErrMsg
	err = args.Get(1).(*ign.ErrMsg)

	return args.Get(0), err
}

// MoveToFront moves a target groupID to the front of the queue.
func (m *Mock) MoveToFront(target interface{}) (interface{}, *ign.ErrMsg) {
	args := m.Called(target)

	var err *ign.ErrMsg
	err = args.Get(1).(*ign.ErrMsg)

	return args.Get(0), err
}

// MoveToBack moves a target element to the front of the queue.
func (m *Mock) MoveToBack(target interface{}) (interface{}, *ign.ErrMsg) {
	args := m.Called(target)

	var err *ign.ErrMsg
	err = args.Get(1).(*ign.ErrMsg)

	return args.Get(0), err
}

// Swap switch places between groupID A and groupID B.
func (m *Mock) Swap(a interface{}, b interface{}) (interface{}, *ign.ErrMsg) {
	args := m.Called(a, b)

	var err *ign.ErrMsg
	err = args.Get(1).(*ign.ErrMsg)

	return args.Get(0), err
}

// Remove removes a groupID from the queue.
func (m *Mock) Remove(id interface{}) (interface{}, *ign.ErrMsg) {
	args := m.Called(id)

	var err *ign.ErrMsg
	err = args.Get(1).(*ign.ErrMsg)

	return args.Get(0), err
}

// Count returns the length of the underlying queue's slice
func (m *Mock) Count() int {
	args := m.Called()
	return args.Int(0)
}

// Initializes a new mock Queue implementation
func NewMock() Queue {
	var q Queue
	q = &Mock{}
	return q
}
