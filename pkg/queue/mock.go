package queue

import (
	"github.com/stretchr/testify/mock"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

type Mock struct {
	mock.Mock
}

func (m *Mock) Get(offset, limit *int) ([]interface{}, *ign.ErrMsg) {
	args := m.Called(offset, limit)

	var result []interface{}
	result = append(result, args.Get(0))

	var err *ign.ErrMsg
	err = args.Get(1).(*ign.ErrMsg)

	return result, err
}

func (m *Mock) Enqueue(entity interface{}) interface{} {
	args := m.Called(entity)
	return args.Get(0)
}

func (m *Mock) Dequeue() (interface{}, *ign.ErrMsg) {
	args := m.Called()

	var err *ign.ErrMsg
	err = args.Get(1).(*ign.ErrMsg)

	return args.Get(0), err
}

func (m *Mock) DequeueOrWait() (interface{}, *ign.ErrMsg) {
	args := m.Called()

	var err *ign.ErrMsg
	err = args.Get(1).(*ign.ErrMsg)

	return args.Get(0), err
}

func (m *Mock) MoveToFront(target interface{}) (interface{}, *ign.ErrMsg) {
	args := m.Called(target)

	var err *ign.ErrMsg
	err = args.Get(1).(*ign.ErrMsg)

	return args.Get(0), err
}

func (m *Mock) MoveToBack(target interface{}) (interface{}, *ign.ErrMsg) {
	args := m.Called(target)

	var err *ign.ErrMsg
	err = args.Get(1).(*ign.ErrMsg)

	return args.Get(0), err
}

func (m *Mock) Swap(a interface{}, b interface{}) (interface{}, *ign.ErrMsg) {
	args := m.Called(a, b)

	var err *ign.ErrMsg
	err = args.Get(1).(*ign.ErrMsg)

	return args.Get(0), err
}

func (m *Mock) Remove(id interface{}) (interface{}, *ign.ErrMsg) {
	args := m.Called(id)

	var err *ign.ErrMsg
	err = args.Get(1).(*ign.ErrMsg)

	return args.Get(0), err
}

func (m *Mock) Count() int {
	args := m.Called()
	return args.Int(0)
}

func NewMock() Queue {
	var q Queue
	q = &Mock{}
	return q
}