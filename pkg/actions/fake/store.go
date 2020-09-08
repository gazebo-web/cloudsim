package fake

import (
	"errors"
	"fmt"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	"reflect"
	"sync"
)

// store is an action.Store implementation.
type store struct {
	value     interface{}
	lockValue sync.RWMutex
	valueType reflect.Type
	persisted bool
	data      interface{}
}

func (s *store) Get() interface{} {
	s.lockValue.RLock()
	defer s.lockValue.RUnlock()
	return s.value
}

func (s *store) set(value interface{}) error {
	inputType := reflect.TypeOf(value)

	if !inputType.AssignableTo(s.valueType) {
		return errors.New("input type is not assignable to underlying data type")
	}

	s.value = value

	return nil
}

func (s *store) Set(value interface{}) error {
	s.lockValue.Lock()
	defer s.lockValue.Unlock()

	err := s.set(value)
	if err != nil {
		return fmt.Errorf("invalid error data type received on Set value. Error: %+v", err)
	}

	return nil
}

func (s *store) Load() error {
	s.lockValue.Lock()
	defer s.lockValue.Unlock()
	err := s.set(s.data)
	if err != nil {
		return err
	}
	return nil
}

func (s *store) Save() error {
	s.lockValue.RLock()
	defer s.lockValue.RUnlock()
	s.persisted = true
	s.data = s.value
	return nil
}

func NewFakeStore(value interface{}) actions.Store {
	return &store{
		value:     value,
		valueType: reflect.TypeOf(value),
		data:      value,
	}
}
