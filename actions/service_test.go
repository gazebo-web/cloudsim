package actions

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type ServiceTestStruct struct {
	PreHook  int
	Execute  int
	PostHook int
}

var serviceTestData = struct {
	// Errors
	errExecute  error
	errRollback error

	// Job Data
	// Job 1
	job1InputData ServiceTestStruct
	job1JobData   ServiceTestStruct
	// Job 2
	job2InputData ServiceTestStruct
	job2JobData   ServiceTestStruct
	// Job 3
	job3InputData ServiceTestStruct
	job3JobData   ServiceTestStruct
	// Rollback
	jobRollbackJobData ServiceTestStruct

	// Helper functions
	getJobDataCount         func(t *testing.T, db *gorm.DB, deployment *Deployment) int
	getDeploymentErrorCount func(t *testing.T, db *gorm.DB, deployment *Deployment) int
	createJobs              func(t *testing.T, rollbackHandlerCalls *[]int) Jobs
	execute                 func(t *testing.T, ctx Context, db *gorm.DB, service *Service, jobs Jobs,
		executeInput *ExecuteInput, jobInput interface{}, errorExpected bool) *Deployment
	processJobs func(t *testing.T, tr *TestResource, service *Service, executeInput *ExecuteInput,
		jobInput interface{}, jobs Jobs) (*ExecuteInput, error)
}{
	// Errors
	errExecute:  errors.New("execute"),
	errRollback: errors.New("rollback"),

	// Job Data
	// Job 1
	job1InputData: ServiceTestStruct{0, 0, 0},
	job1JobData:   ServiceTestStruct{1, 1, 1},
	// Job 2
	job2InputData: ServiceTestStruct{1, 1, 1},
	job2JobData:   ServiceTestStruct{2, 2, 2},
	// Job 3
	job3InputData: ServiceTestStruct{2, 2, 2},
	job3JobData:   ServiceTestStruct{3, 3, 3},
	// Rollback
	jobRollbackJobData: ServiceTestStruct{-1, -1, -1},

	// Helper functions
	// getJobDataCount returns the number of JobData entries in the database
	getJobDataCount: func(t *testing.T, db *gorm.DB, deployment *Deployment) int {
		var count int
		err := db.
			Model(&deploymentData{}).
			Where("deployment_id = ?", deployment.ID).
			Count(&count).Error
		require.NoError(t, err)

		return count
	},

	// getDeploymentErrorCount returns the number of deployment error entries in the database
	getDeploymentErrorCount: func(t *testing.T, db *gorm.DB, deployment *Deployment) int {
		var count int
		err := db.
			Model(&DeploymentError{}).
			Where("deployment_id = ?", deployment.ID).
			Count(&count).Error
		require.NoError(t, err)

		return count
	},

	// processJobs runs service.processJobs on a slice of jobs
	processJobs: func(t *testing.T, tr *TestResource, service *Service, executeInput *ExecuteInput,
		jobInput interface{}, jobs Jobs) (*ExecuteInput, error) {
		if executeInput == nil {
			executeInput = &ExecuteInput{}
		}

		// Create the action
		action := &Action{
			Jobs: jobs,
		}

		// Initialize the input
		require.NoError(t, executeInput.initialize(tr.db, action))

		// Process jobs
		err := service.processJobs(tr.ctx, tr.db, action, executeInput, jobInput)

		return executeInput, err
	},

	// createJob returns a test job that modifies and returns ServiceTestStruct type input
	// A pointer to an empty slice can be passed to rollbackHandlerCalls to record rollback calls of jobs created by
	// this function.
	createJobs: func(t *testing.T, rollbackHandlerCalls *[]int) Jobs {
		td := getTestData(t)
		createJob := func(jobIndex int, name string) *Job {
			// Prepare the rollback handler if necessary
			var rollbackHandler JobErrorHandler
			if rollbackHandlerCalls != nil {
				rollbackHandler = func(ctx Context, tx *gorm.DB, deployment *Deployment, value interface{},
					err error) (interface{}, error) {

					// The error received should be the test rollback error
					require.Equal(t, errors.New("rollback"), err)

					// Register the rollback handler call
					*rollbackHandlerCalls = append(*rollbackHandlerCalls, jobIndex)

					data, err := deployment.GetJobData(tx, nil, deploymentJobData)
					jobData := data.(*ServiceTestStruct)

					require.Equal(t, jobData.PreHook, jobIndex)
					require.Equal(t, jobData.Execute, jobIndex)
					require.Equal(t, jobData.PostHook, jobIndex)

					// Update the job data
					jobRollbackDataJob := &ServiceTestStruct{-1, -1, -1}
					require.NoError(t, deployment.SetJobData(tx, nil, deploymentJobData, jobRollbackDataJob))

					return nil, nil
				}
			}

			return &Job{
				Name: name,
				PreHooks: []JobFunc{
					func(ctx Context, tx *gorm.DB, deployment *Deployment, value interface{}) (interface{}, error) {
						input := value.(*ServiceTestStruct)

						input.PreHook++
						require.Equal(t, jobIndex, input.PreHook)

						return input, nil
					},
				},

				Execute: func(ctx Context, tx *gorm.DB, deployment *Deployment, value interface{}) (interface{}, error) {
					input := value.(*ServiceTestStruct)

					input.Execute++
					require.Equal(t, jobIndex, input.Execute)

					return input, nil
				},

				PostHooks: []JobFunc{
					func(ctx Context, tx *gorm.DB, deployment *Deployment, value interface{}) (interface{}, error) {
						input := value.(*ServiceTestStruct)

						input.PostHook++
						require.Equal(t, jobIndex, input.PostHook)

						// Create a job data entry
						require.NoError(t, deployment.SetJobData(tx, nil, deploymentJobData, input))

						return input, nil
					},
				},
				RollbackHandler: rollbackHandler,
				InputType:       GetJobDataType(&ServiceTestStruct{}),
				OutputType:      NilJobDataType,
			}
		}

		return Jobs{
			createJob(1, td.jobName1),
			createJob(2, td.jobName2),
			createJob(3, td.jobName3),
		}
	},

	// execute performs setup operations and calls service.Execute
	execute: func(t *testing.T, ctx Context, db *gorm.DB, service *Service, jobs Jobs, executeInput *ExecuteInput,
		jobInput interface{}, errorExpected bool) *Deployment {

		td := getTestData(t)

		// Create and register the action
		action, err := NewAction(jobs)
		require.NoError(t, err)
		require.NoError(t, service.RegisterAction(nil, td.actionName, action))

		// Initialize the input
		if executeInput == nil {
			executeInput = &ExecuteInput{
				ActionName: td.actionName,
			}
			require.NoError(t, executeInput.initialize(db, action))
		}
		deployment := executeInput.getDeployment()

		// Execute the action
		err = service.Execute(ctx, db, executeInput, jobInput)
		if errorExpected {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}

		return deployment
	},
}

// newTestService creates a new service containing the test action
func newTestService(t *testing.T) *Service {
	td := getTestData(t)

	service := NewService()

	if err := service.RegisterAction(&td.applicationName, td.actionName, td.action); err != nil {
		panic(fmt.Sprintf("failed to register action %s", td.actionName))
	}

	return service
}

func TestServiceGenerateApplicationName(t *testing.T) {
	td := getTestData(t)

	// Without Application
	expected := td.actionName
	value, _ := generateApplicationActionName(nil, td.actionName)
	assert.Equal(t, expected, value)

	// With application
	expected = td.applicationName + td.actionName
	value, _ = generateApplicationActionName(&td.applicationName, td.actionName)
	assert.Equal(t, expected, value)
}

func TestServiceRegisterAction(t *testing.T) {
	service := newTestService(t)

	applicationName := "App"
	actionName := "Action"
	action := Action{
		Jobs: nil,
	}

	// Generic action
	assert.NoError(t, service.RegisterAction(nil, actionName, &action))
	_, err := service.getAction(nil, actionName)
	assert.NoError(t, err)

	// Application-specific action
	assert.NoError(t, service.RegisterAction(&applicationName, actionName, &action))
	_, err = service.getAction(nil, actionName)
	assert.NoError(t, err)
}

func TestServiceRegisterExistingAction(t *testing.T) {
	td := getTestData(t)

	service := newTestService(t)

	assert.Error(t, service.RegisterAction(&td.applicationName, td.actionName, td.action))
}

func TestServiceRegisterActionNoAction(t *testing.T) {
	service := newTestService(t)

	td := getTestData(t)

	assert.Error(t, service.RegisterAction(&td.applicationName, td.actionName, nil))
}

func TestServiceRegisterActionNoActionName(t *testing.T) {
	td := getTestData(t)

	service := newTestService(t)

	assert.Error(t, service.RegisterAction(&td.applicationName, "", td.action))
}

func TestServiceGetExistingAction(t *testing.T) {
	td := getTestData(t)

	service := newTestService(t)

	// Get the test action
	action, err := service.getAction(&td.applicationName, td.actionName)
	assert.NotNil(t, action)
	assert.NoError(t, err)
}

func TestServiceGetNonexistentAction(t *testing.T) {
	service := newTestService(t)
	actionName := "NonexistentAction"

	_, err := service.getAction(nil, actionName)

	assert.Error(t, err)
}

func TestProcessJobs(t *testing.T) {
	tr := setupTest(t)
	td := getTestData(t)
	std := serviceTestData
	service := newTestService(t)

	// Reset the job data type registry
	jobDataTypeRegistry = newDataTypeRegistry()

	// Create the set of jobs
	jobs := std.createJobs(t, nil)
	jobCount := len(jobs)
	// Register job types
	for _, job := range jobs {
		job.registerTypes(jobDataTypeRegistry)
	}

	// checkJobData validates the job data stored by each job
	checkJobData := func(t *testing.T, tr *TestResource, deployment *Deployment, job string, dataType deploymentDataType,
		inputData *ServiceTestStruct) {
		out, err := deployment.GetJobData(tr.db, &job, dataType)
		require.NoError(t, err)
		require.Equal(t, *inputData, *out.(*ServiceTestStruct))
	}

	test := func(deployment *Deployment, jobInput *ServiceTestStruct) {
		// Verify that the values were processed as expected
		if jobInput != nil {
			require.Equal(t, jobCount, jobInput.PreHook)
			require.Equal(t, jobCount, jobInput.Execute)
			require.Equal(t, jobCount, jobInput.PostHook)
		}

		// Check that the job data was recorded
		// The number of job data entries should be twice the number of jobs (input + job)
		require.Equal(t, std.getJobDataCount(t, tr.db, deployment), jobCount*2)

		// Job 1
		checkJobData(t, tr, deployment, td.jobName1, deploymentJobInput, &std.job1InputData)
		checkJobData(t, tr, deployment, td.jobName1, deploymentJobData, &std.job1JobData)
		// Job 2
		checkJobData(t, tr, deployment, td.jobName2, deploymentJobInput, &std.job2InputData)
		checkJobData(t, tr, deployment, td.jobName2, deploymentJobData, &std.job2JobData)
		// Job 3
		checkJobData(t, tr, deployment, td.jobName3, deploymentJobInput, &std.job3InputData)
		checkJobData(t, tr, deployment, td.jobName3, deploymentJobData, &std.job3JobData)
	}

	// Process the jobs from scratch
	t.Run("Process jobs from scratch", func(t *testing.T) {
		jobInput := &ServiceTestStruct{}

		input, err := std.processJobs(t, tr, service, nil, jobInput, jobs)
		require.NoError(t, err)
		deployment := input.getDeployment()

		test(deployment, jobInput)
	})

	// Process jobs resuming from the second job
	t.Run("Process jobs resuming from second job", func(t *testing.T) {
		deployment, err := newDeployment(tr.db, &Action{Jobs: jobs})
		require.NoError(t, err)
		require.NoError(t, deployment.setJob(tr.db, td.jobName2, nil))

		executeInput := &ExecuteInput{
			Deployment: deployment,
			index:      1,
		}

		// Set job data
		require.NoError(t, deployment.SetJobData(tr.db, &td.jobName1, deploymentJobInput, &std.job1InputData))
		require.NoError(t, deployment.SetJobData(tr.db, &td.jobName1, deploymentJobData, &std.job1JobData))
		require.NoError(t, deployment.SetJobData(tr.db, &td.jobName2, deploymentJobInput, &std.job2InputData))

		input, err := std.processJobs(t, tr, service, executeInput, nil, jobs)
		require.NoError(t, err)
		deployment = input.getDeployment()

		test(deployment, nil)
	})
}

func TestProcessJobsFailedJobErrorHandler(t *testing.T) {
	tr := setupTest(t)
	td := getTestData(t)
	std := serviceTestData
	service := newTestService(t)

	// Job
	testErr := errors.New("posthooks")
	job := Job{
		Name: td.jobName1,
		PreHooks: []JobFunc{
			WrapErrorHandler(
				// Fun
				func(ctx Context, tx *gorm.DB, deployment *Deployment, value interface{}) (interface{}, error) {
					return value, errors.New("prehooks")
				},
				// Handler
				func(ctx Context, tx *gorm.DB, deployment *Deployment, value interface{}, err error) (interface{}, error) {
					input := value.(*ServiceTestStruct)

					input.PreHook++

					return input, nil
				},
			),
		},
		Execute: WrapErrorHandler(
			// Fn
			func(ctx Context, tx *gorm.DB, deployment *Deployment, value interface{}) (interface{}, error) {
				return value, errors.New("execute")
			},
			// Handler
			func(ctx Context, tx *gorm.DB, deployment *Deployment, value interface{}, err error) (interface{}, error) {
				input := value.(*ServiceTestStruct)

				input.Execute++

				return input, nil
			},
		),
		PostHooks: []JobFunc{
			WrapErrorHandler(
				// Fun
				func(ctx Context, tx *gorm.DB, deployment *Deployment, value interface{}) (interface{}, error) {
					return value, testErr
				},
				// Handler
				func(ctx Context, tx *gorm.DB, deployment *Deployment, value interface{}, err error) (interface{}, error) {
					input := value.(*ServiceTestStruct)

					input.PostHook++

					// Fail to handle the error
					return input, err
				},
			),
		},
	}
	jobCount := 1

	// Process the jobs
	jobInput := &ServiceTestStruct{}
	input, err := std.processJobs(t, tr, service, nil, jobInput, Jobs{&job})
	require.Error(t, testErr, err)
	deployment := input.getDeployment()

	// Verify that the values were processed as expected
	require.Equal(t, jobCount, jobInput.PreHook)
	require.Equal(t, jobCount, jobInput.Execute)
	require.Equal(t, jobCount, jobInput.PostHook)

	// A job data entry for the input should have been stored
	require.Equal(t, std.getJobDataCount(t, tr.db, deployment), jobCount)

	// There should be 3 errors registered for the job
	errs, err := deployment.GetErrors(tr.db, nil)
	require.NoError(t, err)
	require.Len(t, errs, 4)
}

func TestProcessJobsNilOutput(t *testing.T) {
	tr := setupTest(t)
	std := serviceTestData
	service := newTestService(t)
	jobInput := &ServiceTestStruct{}

	test := func(job *Job, expectedErr error) {
		// Process the jobs
		input, err := std.processJobs(t, tr, service, nil, jobInput, Jobs{job})
		require.Error(t, expectedErr, err)
		deployment := input.getDeployment()

		// There should be only 1 error registered for the job
		errs, err := deployment.GetErrors(tr.db, nil)
		require.NoError(t, err)
		require.Len(t, errs, 1)
	}

	// Job that returns nil and no error
	jobNil := &Job{
		Execute: func(ctx Context, tx *gorm.DB, deployment *Deployment, value interface{}) (interface{}, error) {
			return nil, nil
		},
	}
	test(jobNil, ErrJobNilOutput)

	// Job that returns nil and an error
	testErr := errors.New("test")
	jobTestErr := &Job{
		Execute: func(ctx Context, tx *gorm.DB, deployment *Deployment, value interface{}) (interface{}, error) {
			return nil, testErr
		},
	}
	test(jobTestErr, testErr)
}

// Execute
//   Helper functions
func testServiceValidateExecute(t *testing.T, db *gorm.DB, deployment *Deployment, jobCount int) {
	td := getTestData(t)
	std := serviceTestData

	// Check that the job data was recorded
	// The number of job data entries should be twice the number of jobs (input + job).
	require.Equal(t, jobCount*2, std.getJobDataCount(t, db, deployment))
	// There should be one deployment error for the failed post-hook and one for the rollback handler.
	require.Equal(t, 0, std.getDeploymentErrorCount(t, db, deployment))

	// Validate job data
	// checkJobData validates the job data stored by each job
	checkJobData := func(t *testing.T, db *gorm.DB, deployment *Deployment, job string,
		dataType deploymentDataType, inputData *ServiceTestStruct) {
		out, err := deployment.GetJobData(db, &job, dataType)
		require.NoError(t, err)
		require.Equal(t, *inputData, *out.(*ServiceTestStruct))
	}
	// Job 1
	checkJobData(t, db, deployment, td.jobName1, deploymentJobInput, &std.job1InputData)
	checkJobData(t, db, deployment, td.jobName1, deploymentJobData, &std.job1JobData)
	// Job 2
	checkJobData(t, db, deployment, td.jobName2, deploymentJobInput, &std.job2InputData)
	checkJobData(t, db, deployment, td.jobName2, deploymentJobData, &std.job2JobData)
	// Job 3
	checkJobData(t, db, deployment, td.jobName3, deploymentJobInput, &std.job3InputData)
	checkJobData(t, db, deployment, td.jobName3, deploymentJobData, &std.job3JobData)

	// Check that the deployment is marked as finished
	require.True(t, deployment.isFinished())
}

func testServiceValidateRollbackExecute(t *testing.T, db *gorm.DB, deployment *Deployment, jobCount int) {
	td := getTestData(t)
	std := serviceTestData

	// Check that the job data was recorded
	// The number of job data entries should be twice the number of jobs (input + job).
	require.Equal(t, jobCount*2, std.getJobDataCount(t, db, deployment))
	// There should be one deployment error for the failed post-hook and one for the rollback handler.
	require.Equal(t, 2, std.getDeploymentErrorCount(t, db, deployment))

	// Validate job data
	// checkJobData validates the job data stored by each job
	checkJobData := func(t *testing.T, db *gorm.DB, deployment *Deployment, job string,
		dataType deploymentDataType, inputData *ServiceTestStruct) {
		out, err := deployment.GetJobData(db, &job, dataType)
		require.NoError(t, err)
		require.Equal(t, *inputData, *out.(*ServiceTestStruct))
	}
	// Job 1
	checkJobData(t, db, deployment, td.jobName1, deploymentJobInput, &std.job1InputData)
	checkJobData(t, db, deployment, td.jobName1, deploymentJobData, &std.jobRollbackJobData)
	// Job 2
	checkJobData(t, db, deployment, td.jobName2, deploymentJobInput, &std.job2InputData)
	checkJobData(t, db, deployment, td.jobName2, deploymentJobData, &std.job2JobData)
	// Job 3
	checkJobData(t, db, deployment, td.jobName3, deploymentJobInput, &std.job3InputData)
	checkJobData(t, db, deployment, td.jobName3, deploymentJobData, &std.jobRollbackJobData)

	// Check that the deployment is marked as finished
	require.True(t, deployment.isFinished())
}

func testServiceUpdateJobsForRollback(t *testing.T, jobs Jobs) {
	std := serviceTestData
	jobCount := len(jobs)

	// Make the last posthook fail
	hookFn := jobs[jobCount-1].PostHooks[0]
	jobs[jobCount-1].PostHooks = []JobFunc{
		func(ctx Context, tx *gorm.DB, deployment *Deployment, value interface{}) (interface{}, error) {
			// Execute the posthook logic as usual and return an error
			_, err := hookFn(ctx, tx, deployment, value)
			require.NoError(t, err)

			return nil, std.errRollback
		},
	}

	// Make the rollback handler from Job 1 fail
	rollbackHandler := jobs[0].RollbackHandler
	jobs[0].RollbackHandler = func(ctx Context, tx *gorm.DB, deployment *Deployment, value interface{},
		err error) (interface{}, error) {

		_, err = rollbackHandler(ctx, tx, deployment, value, err)
		require.NoError(t, err)
		return nil, std.errRollback
	}

	// Remove the rollback handler from Job 2
	jobs[1].RollbackHandler = nil
}

//   Tests
func TestExecuteInvalidAction(t *testing.T) {
	tr := setupTest(t)

	service := newTestService(t)
	executeInput := &ExecuteInput{
		ActionName: "invalid_action",
	}
	err := service.Execute(tr.ctx, tr.db, executeInput, nil)
	require.Error(t, ErrActionNotFound, err)
}

func TestExecute(t *testing.T) {
	tr := setupTest(t)
	std := serviceTestData
	service := newTestService(t)

	// Reset the job data type registry
	jobDataTypeRegistry = newDataTypeRegistry()

	// Create the set of jobs
	jobs := std.createJobs(t, nil)
	jobCount := len(jobs)

	// Execute the action
	jobInput := std.job1InputData
	deployment := std.execute(t, tr.ctx, tr.db, service, jobs, nil, &jobInput, false)

	// Verify that the values were processed as expected
	require.Equal(t, 3, jobInput.PreHook)
	require.Equal(t, 3, jobInput.Execute)
	require.Equal(t, 3, jobInput.PostHook)

	testServiceValidateExecute(t, tr.db, deployment, jobCount)
}

func TestExecuteResumeAction(t *testing.T) {
	tr := setupTest(t)
	td := getTestData(t)
	std := serviceTestData
	service := newTestService(t)

	// Reset the job data type registry
	jobDataTypeRegistry = newDataTypeRegistry()

	// Create the set of jobs
	jobs := std.createJobs(t, nil)
	jobCount := len(jobs)

	// Create the action
	action, err := NewAction(jobs)
	require.NoError(t, err)

	// Create ExecuteInput
	executeInput := &ExecuteInput{
		ActionName: td.actionName,
	}
	require.NoError(t, executeInput.initialize(tr.db, action))

	// Create job data and update the deployment to start at the second job
	deployment := executeInput.getDeployment()
	require.NoError(t, deployment.SetJobData(tr.db, &td.jobName1, deploymentJobInput, &std.job1InputData))
	require.NoError(t, deployment.SetJobData(tr.db, &td.jobName1, deploymentJobData, &std.job1JobData))
	require.NoError(t, deployment.SetJobData(tr.db, &td.jobName2, deploymentJobInput, &std.job2InputData))
	require.NoError(t, deployment.setJob(tr.db, td.jobName2, nil))

	deployment = std.execute(t, tr.ctx, tr.db, service, jobs, executeInput, nil, false)

	testServiceValidateExecute(t, tr.db, deployment, jobCount)
}

func TestExecuteRollback(t *testing.T) {
	tr := setupTest(t)
	std := serviceTestData
	service := newTestService(t)

	// Reset the job data type registry
	jobDataTypeRegistry = newDataTypeRegistry()

	// rollbackHandlerCalls keeps track of called rollback handlers
	rollbackHandlerCalls := make([]int, 0)

	// Create the set of jobs
	jobs := std.createJobs(t, &rollbackHandlerCalls)
	jobCount := len(jobs)
	// Modify jobs to trigger and handle a rollback
	testServiceUpdateJobsForRollback(t, jobs)

	// Execute the action
	jobInput := std.job1InputData
	deployment := std.execute(t, tr.ctx, tr.db, service, jobs, nil, &jobInput, true)

	// Verify that the values were processed as expected
	require.Equal(t, 3, jobInput.PreHook)
	require.Equal(t, 3, jobInput.Execute)
	require.Equal(t, 3, jobInput.PostHook)

	testServiceValidateRollbackExecute(t, tr.db, deployment, jobCount)
}

func TestExecuteResumeRollback(t *testing.T) {
	tr := setupTest(t)
	td := getTestData(t)
	std := serviceTestData
	service := newTestService(t)

	// Reset the job data type registry
	jobDataTypeRegistry = newDataTypeRegistry()

	// rollbackHandlerCalls keeps track of the rollback handlers called
	rollbackHandlerCalls := make([]int, 0)

	// Create the set of jobs
	jobs := std.createJobs(t, &rollbackHandlerCalls)
	jobCount := len(jobs)
	// Modify jobs to trigger and handle a rollback
	testServiceUpdateJobsForRollback(t, jobs)

	// The first job's functions should not run
	jobs[0].Execute = func(ctx Context, tx *gorm.DB, deployment *Deployment, value interface{}) (interface{}, error) {
		t.Log("Job 1 Execute function was called instead of rolling back.")
		t.Fail()
		return nil, nil
	}

	// Create the test action
	action, err := NewAction(jobs)
	require.NoError(t, err)

	// Create ExecuteInput
	executeInput := &ExecuteInput{
		ActionName: td.actionName,
	}
	require.NoError(t, executeInput.initialize(tr.db, action))

	// Create job and rollback data and update to rollback from the first stage
	deployment := executeInput.getDeployment()
	require.NoError(t, deployment.setJob(tr.db, td.jobName1, nil))
	require.NoError(t, deployment.setRollbackStatus(tr.db, std.errRollback))
	// Job data
	require.NoError(t, deployment.SetJobData(tr.db, &td.jobName1, deploymentJobInput, &std.job1InputData))
	require.NoError(t, deployment.SetJobData(tr.db, &td.jobName1, deploymentJobData, &std.job1JobData))
	require.NoError(t, deployment.SetJobData(tr.db, &td.jobName2, deploymentJobInput, &std.job2InputData))
	require.NoError(t, deployment.SetJobData(tr.db, &td.jobName2, deploymentJobData, &std.job2JobData))
	require.NoError(t, deployment.SetJobData(tr.db, &td.jobName3, deploymentJobInput, &std.job3InputData))
	require.NoError(t, deployment.SetJobData(tr.db, &td.jobName3, deploymentJobData, &std.job3JobData))
	// Job errors
	require.NoError(t, deployment.addJobError(tr.db, &td.jobName3, std.errExecute))
	// Rollback data
	require.NoError(t, deployment.SetJobData(tr.db, &td.jobName3, deploymentJobData, &std.jobRollbackJobData))

	// Execute the action
	deployment = std.execute(t, tr.ctx, tr.db, service, jobs, executeInput, nil, true)

	testServiceValidateRollbackExecute(t, tr.db, deployment, jobCount)
}
