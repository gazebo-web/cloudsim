package fake

import (
	"errors"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"reflect"
	"sync"
)

// store is an action.Store implementation.
type store struct {
	value     interface{}
	lockValue sync.RWMutex
	valueType reflect.Type
	data      interface{}
}

func (s *store) Get() interface{} {
	s.lockValue.RLock()
	defer s.lockValue.RUnlock()
	return s.value
}

func (s *store) set(value interface{}) {
	s.value = value
}

func (s *store) isAssignable(value interface{}) bool {
	inputType := reflect.TypeOf(value)
	return inputType.AssignableTo(s.valueType)
}

func (s *store) Set(value interface{}) error {
	s.lockValue.Lock()
	defer s.lockValue.Unlock()

	if !s.isAssignable(value) {
		return errors.New("value type received on Set input is not assignable to the underlying value")
	}

	s.set(value)

	return nil
}

func (s *store) Load() error {
	s.lockValue.Lock()
	defer s.lockValue.Unlock()
	s.set(s.data)
	return nil
}

func NewFakeStore(value interface{}) actions.Store {
	return &store{
		value:     value,
		valueType: reflect.TypeOf(value),
		data:      value,
	}
}
