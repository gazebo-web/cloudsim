package actions

import (
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/require"
	"testing"
)

type testData struct {
	// Applications
	applicationName string

	// Job names
	jobName1 string
	jobName2 string
	jobName3 string

	// Actions
	actionName string
	action     *Action
}

func getTestData(t *testing.T) testData {
	// Applications
	applicationName := "test_app"

	// Job names
	jobName1 := "job_1"
	jobName2 := "job_2"
	jobName3 := "job_3"

	// Actions
	createJob := func(name string) *Job {
		return &Job{
			Name: name,
			Execute: func(ctx Context, tx *gorm.DB, deployment *Deployment, value interface{}) (interface{}, error) {
				return value, nil
			},
			InputType:  NilJobDataType,
			OutputType: NilJobDataType,
		}
	}

	actionName := "test_action"
	action, err := NewAction(
		Jobs{
			createJob(jobName1),
			createJob(jobName2),
			createJob(jobName3),
		},
	)
	require.NoError(t, err)
	action.Name = actionName

	return testData{
		// Applications
		applicationName: applicationName,

		// Job names
		jobName1: jobName1,
		jobName2: jobName2,
		jobName3: jobName3,

		// Actions
		actionName: actionName,
		action:     action,
	}
}
