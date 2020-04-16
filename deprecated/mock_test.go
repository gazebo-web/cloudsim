package deprecated

import (
	"fmt"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"reflect"
)

// Mock operations can be configured to:
// - return "not implemented error" (default),
// - "pass through" to underlying implementation,
// - call a mock function.
//
// Calls to mock functions and their results are kept track in an OpTracker.
// If configured, a callback will be called after each Operation is called.
//
// You can also override these functions

// MockFunction is a function that returns a set of predictable values for a specific operation.
// It is meant for testing.
type MockFunction func(args ...interface{}) interface{}

// MockFunctionValidator is a function that validates that a mock function for a specific operation is in a valid
// state. The definition of a valid state depends on the mock function being validated.
type MockFunctionValidator func() error

// MockFunctionGenerator generates a mock function and validator.
type MockFunctionGenerator func(args ...interface{}) (MockFunction, MockFunctionValidator)

// MockCallback is a function called after calling an Operation.
type MockCallback func(args ...interface{})

type Mock struct {
	Tracker                         *OpTracker
	operationMockFunction           map[OpType]MockFunction
	operationMockFunctionValidators map[OpType]MockFunctionValidator
	notificationCallbacks           map[OpType]MockCallback
}

// PassThrough is returned by mock functions whenever the real implementation of a mocked function should be called.
const PassThrough = "PassThrough"

// Reset (re)initializes the mock and tracker, clearing all records, mock functions and callbacks.
func (m *Mock) Reset() {
	m.operationMockFunction = make(map[OpType]MockFunction, 0)
	m.operationMockFunctionValidators = make(map[OpType]MockFunctionValidator, 0)
	m.notificationCallbacks = make(map[OpType]MockCallback, 0)
	m.Tracker = NewOpTracker()
}

// SetMockFunction sets the mock function and validator for a specific OpType.
func (m *Mock) SetMockFunction(op OpType, f MockFunctionGenerator, args ...interface{}) {
	mockFunction, mockFunctionValidator := f(args...)
	m.operationMockFunction[op] = mockFunction
	m.operationMockFunctionValidators[op] = mockFunctionValidator
}

// GetMockResult gets the result of calling the mock function for a specific OpType.
func (m *Mock) GetMockResult(op OpType) interface{} {
	m.Tracker.TrackCall(op)
	value := m.operationMockFunction[op]()
	m.Tracker.TrackCallResult(value)

	return value
}

// ValidateMockStatus determines whether the mock functions are in a valid state.
// States considered valid depend on each mock function.
func (m *Mock) ValidateMockFunctions() bool {
	// Validate all mock functions
	for op, validator := range m.operationMockFunctionValidators {
		if err := validator(); err != nil {
			ign.NewLogger("Mock Validation Error", true, ign.VerbosityDebug).Error(op, err)
			return false
		}
	}

	return true
}

// SetCallback sets the callback for a specific OpType that will be called each time the operation mock function is
// called.
func (m *Mock) SetCallback(op OpType, cb MockCallback) {
	m.notificationCallbacks[op] = cb
}

// InvokeCallback calls the callback for a specific OpType.
func (m *Mock) InvokeCallback(op OpType, args ...interface{}) {
	cb := m.notificationCallbacks[op]
	if cb != nil {
		cb(args...)
	}
}

// FixedValues returns a mocked function that returns a value from a fixed sequence of values.
// If the mocked function is called more times than the number of values defined, it will panic.
// If the first parameter passed to this generator is true, then the generated validator requires
// the mock function to return all fixed values.
func FixedValues(res ...interface{}) (MockFunction, MockFunctionValidator) {
	if res == nil || len(res) == 0 {
		panic("FixedValues requires at least one parameter.")
	}

	// The first parameter is the strict flag
	strict := res[0].(bool)

	// The rest of the parameters are the sequence of fixed values returned
	res = res[1:]

	// Function: Returns a value from a fixed sequence of values each time it's called
	count := 0
	mockFunction := func(args ...interface{}) interface{} {
		value := res[count]
		count++
		return value
	}

	// Validator: If strict, require that all fixed values are returned
	mockFunctionValidator := func() error {
		result := !strict || count == len(res)
		if !result {
			return fmt.Errorf("fixedValues was not called enough times. It was called %d/%d times", count, len(res))
		}
		return nil
	}

	return mockFunction, mockFunctionValidator
}

// ReplaceValue replaces any value with a mock value. The replaced value
// must be contained in a variable if it is from another package, otherwise any
// value can be replaced. This function can be used to replace functions or
// plain values (ints, strings, structs, etc.), and returns a function that
// restores their original value. The returned function is meant to be deferred.
// Parameter value MUST contain a pointer.
// Parameter mock must be of the same type the pointer points to.
func ReplaceValue(value interface{}, mock interface{}) func() {
	valueValue := reflect.ValueOf(value)

	if valueValue.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("value being replaced must be a pointer not [%s]", valueValue.Kind()))
	}

	mockValue := reflect.ValueOf(mock)
	funcElem := valueValue.Elem()

	if funcElem.Type() != mockValue.Type() {
		s := fmt.Sprintf(
			"value being replaced and mock value types are not equal:\n  value: [%s]\n  mock:  [%s]",
			funcElem.Type(),
			mockValue.Type(),
		)
		panic(s)
	}

	originalValue := funcElem.Interface()
	funcElem.Set(mockValue)

	return func() {
		funcElem.Set(reflect.ValueOf(originalValue))
	}
}
