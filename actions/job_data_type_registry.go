package actions

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	// ErrJobDataTypeNotFound is returned when a job data type was not found.
	// This is likely indicative of a problem in a job's InputType or OutputType properties.
	ErrJobDataTypeNotFound = errors.New("job data type not found")
)

// dataTypeRegistry is a registry that maps type names to concrete job data types.
// It is used to automatically marshal and unmarshal data.
type dataTypeRegistry map[string]JobDataType

// newDataTypeRegistry creates and returns a new dataTypeRegistry
func newDataTypeRegistry() dataTypeRegistry {
	return make(map[string]JobDataType)
}

// register registers a new type of job data type in the registry.
func (dtr dataTypeRegistry) register(dataType JobDataType) {
	jobDataTypeRegistry[GetJobDataTypeName(dataType)] = dataType
}

// getType receives a job data type name and returns data type from the registry.
func (dtr dataTypeRegistry) getType(typeName string) (reflect.Type, error) {
	if dataType, ok := dtr[typeName]; ok {
		return dataType, nil
	}

	// Add the type not found to the error message
	err := fmt.Errorf("%s: %s", ErrJobDataTypeNotFound, typeName)

	return nil, err
}

// jobDataTypeRegistry is a global registry shared among different jobs.
var jobDataTypeRegistry = newDataTypeRegistry()
