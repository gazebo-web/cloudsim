package actions

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewAction(t *testing.T) {
	type TestStruct struct{ I int }

	// Create some jobs
	jobs := make(Jobs, 5)
	for i := 0; i < len(jobs); i++ {
		jobs[i] = &Job{
			Name: fmt.Sprintf("job_%d", i+1),
			Execute: func(store Store, tx *gorm.DB, deployment *Deployment, value interface{}) (interface{}, error) {
				return nil, nil
			},
			InputType:  GetJobDataType(&TestStruct{}),
			OutputType: NilJobDataType,
		}
	}

	// Reset the job data registry
	jobDataTypeRegistry = newDataTypeRegistry()

	// Action with all fields
	_, err := NewAction(jobs)
	assert.NoError(t, err)
	// Check that the job data registry contains the expected types
	assert.Len(t, jobDataTypeRegistry, 2)
	types := []string{
		"*actions.TestStruct",
		"struct {}",
	}
	for _, typeName := range types {
		_, ok := jobDataTypeRegistry[typeName]
		assert.True(t, ok)
	}

	// Action with jobs missing required fields
	for i := 0; i < len(jobs); i++ {
		jobs[i].InputType = nil
		jobs[i].OutputType = nil
	}
	_, err = NewAction(jobs)
	assert.Error(t, err)
}

func TestGetJobNames(t *testing.T) {
	job1Name := "test"
	job2Name := "job with space"
	job3Name := "job_3"

	jobs := Jobs{
		{Name: job1Name},
		{Name: job2Name},
		{Name: job3Name},
	}

	action := &Action{Jobs: jobs}
	require.Equal(t, job1Name+","+job2Name+","+job3Name, *action.getJobNames())
}
