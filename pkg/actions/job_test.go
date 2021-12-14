package actions

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestDefaultExecutePanics(t *testing.T) {
	j := &Job{}

	// The default job should panic
	assert.Panics(t, func() { _, _ = j.Execute(nil, nil, nil, struct{}{}) })
}

func TestDefaultRunPanics(t *testing.T) {
	j := &Job{}

	// The default job should panic
	assert.Panics(t, func() { _, _ = j.Run(nil, nil, nil, struct{}{}) })
}

func TestRegisterTypes(t *testing.T) {
	type TestStruct struct{}

	// Reset the job data type registry
	jobDataTypeRegistry = newDataTypeRegistry()

	job1 := Job{
		InputType:  GetJobDataType(TestStruct{}),
		OutputType: NilJobDataType,
	}
	job2 := Job{
		OutputType: NilJobDataType,
		InputType:  GetJobDataType(&TestStruct{}),
	}

	job1.registerTypes(dataTypeRegistry{})
	job2.registerTypes(dataTypeRegistry{})

	require.Len(t, jobDataTypeRegistry, 3)

	types := []string{
		"actions.TestStruct",
		"*actions.TestStruct",
		"struct {}",
	}
	for _, typeName := range types {
		_, ok := jobDataTypeRegistry[typeName]
		assert.True(t, ok)
	}
}

func TestProcessHooksPass(t *testing.T) {
	j := &Job{}

	// TestResource hooks
	j.PreHooks = []JobFunc{
		func(store Store, tx *gorm.DB, deployment *Deployment, value interface{}) (interface{}, error) {
			return value.(int) + 1, nil
		},
		func(store Store, tx *gorm.DB, deployment *Deployment, value interface{}) (interface{}, error) {
			return value.(int) + 2, nil
		},
		func(store Store, tx *gorm.DB, deployment *Deployment, value interface{}) (interface{}, error) {
			return value.(int) + 3, nil
		},
	}

	value, err := j.processHooks(nil, nil, nil, 0, &j.PreHooks)
	assert.NoError(t, err)
	assert.Equal(t, value, 6)
}

func TestProcessHooksFail(t *testing.T) {
	j := &Job{}

	// Test hooks
	j.PreHooks = []JobFunc{
		func(store Store, tx *gorm.DB, deployment *Deployment, value interface{}) (interface{}, error) {
			return value.(int) + 1, nil
		},
		func(store Store, tx *gorm.DB, deployment *Deployment, value interface{}) (interface{}, error) {
			return value.(int) + 2, assert.AnError
		},
		func(store Store, tx *gorm.DB, deployment *Deployment, value interface{}) (interface{}, error) {
			return value.(int) + 3, nil
		},
	}

	value, err := j.processHooks(nil, nil, nil, 0, &j.PreHooks)
	assert.Nil(t, value)
	assert.EqualError(t, err, assert.AnError.Error())
}

func TestCallJobFunc(t *testing.T) {
	valueFunc := func(store Store, tx *gorm.DB, deployment *Deployment, value interface{}) (interface{}, error) {
		return value, nil
	}
	nilFunc := func(store Store, tx *gorm.DB, deployment *Deployment, value interface{}) (interface{}, error) {
		return nil, nil
	}

	test := func(jobFunc JobFunc, value interface{}, error bool) {
		value, err := callJobFunc(jobFunc, nil, nil, nil, value)
		if error {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}
	}

	test(valueFunc, nil, false)
	test(valueFunc, 1, false)
	test(nilFunc, nil, false)
	test(nilFunc, 1, false)
}

func TestTestJobExecute(t *testing.T) {
	valueFunc := func(store Store, tx *gorm.DB, deployment *Deployment, value interface{}) (interface{}, error) {
		return value, nil
	}
	nilFunc := func(store Store, tx *gorm.DB, deployment *Deployment, value interface{}) (interface{}, error) {
		return nil, nil
	}

	test := func(execute JobFunc, postHook JobFunc, value interface{}, error bool) {
		// Prepare job
		j := &Job{
			Execute: execute,
		}

		// Optionally add post-hooks
		if postHook != nil {
			j.PostHooks = []JobFunc{postHook}
		}

		value, err := j.Execute(nil, nil, nil, value)

		if error {
			assert.NoError(t, err)
		} else {
			// The test job should not fail and return the same value it receives
			assert.NoError(t, err)
			assert.Equal(t, value.(int), value)
		}
	}

	test(valueFunc, nil, 1, false)
	test(nilFunc, nil, 1, true)
	test(valueFunc, valueFunc, 1, false)
	test(valueFunc, nilFunc, 1, true)
}

func TestTestJobRun(t *testing.T) {
	val := 1
	j := &Job{
		PreHooks: []JobFunc{
			// Multiply the input value by two
			func(store Store, tx *gorm.DB, deployment *Deployment, value interface{}) (interface{}, error) {
				return value.(int) * 2, nil
			},
			// Check that the input value is now two times val
			func(store Store, tx *gorm.DB, deployment *Deployment, value interface{}) (interface{}, error) {
				var err error
				if value.(int) != val*2 {
					err = assert.AnError
				}
				return value, err
			},
		},
		// Check that the input value is two times val
		Execute: func(store Store, tx *gorm.DB, deployment *Deployment, value interface{}) (interface{}, error) {
			var err error
			if value.(int) != val*2 {
				err = assert.AnError
			}
			return value, err
		},
		// Divide output value by two
		PostHooks: []JobFunc{
			// Check that the output value is two times val
			func(store Store, tx *gorm.DB, deployment *Deployment, value interface{}) (interface{}, error) {
				var err error
				if value.(int) != val*2 {
					err = assert.AnError
				}
				return value, err
			},
			// Divide the output value by two
			func(store Store, tx *gorm.DB, deployment *Deployment, value interface{}) (interface{}, error) {
				return value.(int) / 2, nil
			},
			// Check that the output value is now val
			func(store Store, tx *gorm.DB, deployment *Deployment, value interface{}) (interface{}, error) {
				var err error
				if value.(int) != val {
					err = assert.AnError
				}
				return value, err
			},
		},
	}

	value, err := j.Run(nil, nil, nil, val)

	// The test job should not return an error and should return the same
	// value it receives
	assert.NoError(t, err)
	assert.Equal(t, value.(int), val)
}

func TestGetJobDataType(t *testing.T) {
	test := func(value interface{}) {
		// Check value
		require.Equal(t, reflect.TypeOf(value), GetJobDataType(value))
		// Check reflect.Type
		require.Equal(t, reflect.TypeOf(&value), GetJobDataType(reflect.TypeOf(&value)))
	}

	// Basic types
	test("")
	test(123)
	test(true)

	// Struct
	type TestInput struct{}
	test(TestInput{})

	// Nil
	test(nil)
}

func TestGetJobDataTypeName(t *testing.T) {
	test := func(value interface{}, typeName string) {
		// Check value
		require.Equal(t, typeName, GetJobDataTypeName(value))
		// Check reflect.Type
		require.Equal(t, typeName, GetJobDataTypeName(reflect.TypeOf(value)))
	}

	// Basic types
	str := ""
	test(str, "string")
	test(&str, "*string")
	int := 123
	test(int, "int")
	test(&int, "*int")
	b := false
	test(b, "bool")
	test(&b, "*bool")

	// Struct
	type TestInput struct{}
	input := TestInput{}
	test(input, "actions.TestInput")
	test(&input, "*actions.TestInput")

	// Nil
	test(nil, "nil")
}

func TestJobsRegisterTypes(t *testing.T) {
	type TestStruct struct{}

	// Reset the job data type registry
	jobDataTypeRegistry = newDataTypeRegistry()

	jobs := Jobs{
		{
			InputType:  GetJobDataType(TestStruct{}),
			OutputType: NilJobDataType,
		},
		{
			OutputType: NilJobDataType,
			InputType:  GetJobDataType(&TestStruct{}),
		},
	}

	jobs.registerTypes(dataTypeRegistry{})

	require.Len(t, jobDataTypeRegistry, 3)

	types := []string{
		"actions.TestStruct",
		"*actions.TestStruct",
		"struct {}",
	}
	for _, typeName := range types {
		_, ok := jobDataTypeRegistry[typeName]
		assert.True(t, ok)
	}
}

func TestJobsValidate(t *testing.T) {
	// Check that an empty Job fails
	emptyJob := &Job{}
	require.Error(t, emptyJob.validate())

	// Check that all required fields are filled
	data := "test"
	job := &Job{
		Name: "test",
		Execute: func(store Store, tx *gorm.DB, deployment *Deployment, value interface{}) (interface{}, error) {
			return nil, nil
		},
		// Nil type input
		InputType:  GetJobDataType(NilJobDataType),
		OutputType: GetJobDataType(data),
	}
	require.NoError(t, job.validate())
}

func TestJobsValidateFailNoJobs(t *testing.T) {
	// Check that an empty Job fails
	emptyJob := &Job{}
	require.Error(t, emptyJob.validate())

	jobs := &Jobs{}
	require.Equal(t, ErrJobsEmptySequence, jobs.notEmpty())
}

func TestJobsValidateFailInvalidSingleJob(t *testing.T) {
	// Check that an empty Job fails
	emptyJob := &Job{}
	require.Error(t, emptyJob.validate())

	jobs := &Jobs{emptyJob}
	require.NotNil(t, jobs.jobsAreValid())
}

func TestJobsNoRepeatedJobNames(t *testing.T) {
	// Check that an empty Job fails
	emptyJob := &Job{}
	require.Error(t, emptyJob.validate())

	// Prepare a valid job
	job := &Job{
		Name: "test",
		Execute: func(store Store, tx *gorm.DB, deployment *Deployment, value interface{}) (interface{}, error) {
			return nil, nil
		},
		// Nil type input
		InputType:  GetJobDataType(NilJobDataType),
		OutputType: GetJobDataType("test"),
	}
	assert.Nil(t, job.validate())

	jobs := &Jobs{job, job}
	require.True(t, errors.Is(jobs.jobNamesAreUnique(), ErrJobsNamesNotUnique))
}

func TestJobValidate(t *testing.T) {
	// Check that an empty Job fails
	emptyJob := &Job{}
	require.Error(t, emptyJob.validate())

	// Check that all required fields are filled
	data := "test"
	job := &Job{
		Name: "test",
		Execute: func(store Store, tx *gorm.DB, deployment *Deployment, value interface{}) (interface{}, error) {
			return nil, nil
		},
		// Nil type input
		InputType:  GetJobDataType(NilJobDataType),
		OutputType: GetJobDataType(data),
	}
	require.NoError(t, job.validate())
}

func TestExtendJob(t *testing.T) {
	jobName := "test_job"
	jobVar := &Job{
		Name: jobName,
		Execute: func(store Store, tx *gorm.DB, deployment *Deployment, value interface{}) (interface{}, error) {
			return jobName, nil
		},
		RollbackHandler: func(store Store, tx *gorm.DB, deployment *Deployment, value interface{}, err error) (interface{}, error) {
			return jobName, nil
		},
	}

	getFuncPtr := func(fn interface{}) uintptr {
		return reflect.ValueOf(fn).Pointer()
	}

	// Replacing the execute function should panic

	require.Panics(t, func() {
		extension := Job{
			Execute: func(store Store, tx *gorm.DB, deployment *Deployment, value interface{}) (interface{}, error) {
				return true, nil
			},
		}
		jobVar.Extend(extension)
	})

	// Create the extension
	hook := func(store Store, tx *gorm.DB, deployment *Deployment, value interface{}) (interface{}, error) {
		return nil, nil
	}
	rollback := JobErrorHandler(func(store Store, tx *gorm.DB, deployment *Deployment, value interface{}, err error) (interface{}, error) {
		return fmt.Sprintf("%s-test", jobName), nil
	})
	extendedJob := jobVar.Extend(Job{
		PreHooks:        []JobFunc{hook},
		PostHooks:       []JobFunc{hook},
		RollbackHandler: rollback,
	})

	// Ensure that the name and Execute functions remain the same
	require.Equal(t, jobVar.Name, extendedJob.Name)
	require.Equal(t, getFuncPtr(jobVar.Execute), getFuncPtr(extendedJob.Execute))

	// Ensure the original job was not modified
	require.Nil(t, jobVar.PreHooks)
	require.Nil(t, jobVar.PostHooks)
	require.NotNil(t, jobVar.RollbackHandler)

	// Check that the extended job contains extended properties
	require.NotNil(t, extendedJob.PreHooks)
	require.Equal(t, 1, len(extendedJob.PreHooks))
	require.Equal(t, getFuncPtr(hook), getFuncPtr(extendedJob.PreHooks[0]))

	require.NotNil(t, extendedJob.PostHooks)
	require.Equal(t, 1, len(extendedJob.PostHooks))
	require.Equal(t, getFuncPtr(hook), getFuncPtr(extendedJob.PostHooks[0]))

	require.NotNil(t, jobVar.RollbackHandler)
	require.Equal(t, getFuncPtr(rollback), getFuncPtr(jobVar.RollbackHandler))
}
