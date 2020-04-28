package queue

import "gitlab.com/ignitionrobotics/web/ign-go"

type Mock struct {
	GetMock func(offset, limit *int) ([]interface{}, *ign.ErrMsg)
	EnqueueMock func(entity interface{}) interface{}
	DequeueMock func() (interface{}, *ign.ErrMsg)
	DequeueOrWaitMock func() (interface{}, *ign.ErrMsg)
	MoveToFrontMock func(target interface{}) (interface{}, *ign.ErrMsg)
	MoveToBackMock func(target interface{}) (interface{}, *ign.ErrMsg)
	SwapMock func(a interface{}, b interface{}) (interface{}, *ign.ErrMsg)
	RemoveMock func(id interface{}) (interface{}, *ign.ErrMsg)
	CountMock func() int
}

func (m *Mock) Get(offset, limit *int) ([]interface{}, *ign.ErrMsg) {
	return m.GetMock(offset, limit)
}

func (m *Mock) Enqueue(entity interface{}) interface{} {
	return m.EnqueueMock(entity)
}

func (m *Mock) Dequeue() (interface{}, *ign.ErrMsg) {
	return m.DequeueMock()
}

func (m *Mock) DequeueOrWait() (interface{}, *ign.ErrMsg) {
	return m.DequeueOrWaitMock()
}

func (m *Mock) MoveToFront(target interface{}) (interface{}, *ign.ErrMsg) {
	return m.MoveToFrontMock(target)
}

func (m *Mock) MoveToBack(target interface{}) (interface{}, *ign.ErrMsg) {
	return m.MoveToBackMock(target)
}

func (m *Mock) Swap(a interface{}, b interface{}) (interface{}, *ign.ErrMsg) {
	return m.SwapMock(a, b)
}

func (m *Mock) Remove(id interface{}) (interface{}, *ign.ErrMsg) {
	return m.RemoveMock(id)
}

func (m *Mock) Count() int {
	return m.CountMock()
}

func NewMock() IQueue {
	var q IQueue
	q = &Mock{}
	return q
}