package actions

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type deploymentTestStruct struct {
	I int
}

var deploymentTestData = struct {
	// Actions
	action *Action

	// Helper functions
	getDeploymentJobDataCount func(t *testing.T, db *gorm.DB) int
}{
	// Actions
	action: &Action{
		Name: "test_action",
		Jobs: Jobs{
			{
				Name: "job_1",
				Execute: func(store Store, tx *gorm.DB, deployment *Deployment, value interface{}) (interface{}, error) {
					return value, nil
				},
				InputType:  NilJobDataType,
				OutputType: GetJobDataType(&deploymentTestStruct{}),
			},
			{
				Name: "job_2",
				Execute: func(store Store, tx *gorm.DB, deployment *Deployment, value interface{}) (interface{}, error) {
					return value, nil
				},
				InputType:  NilJobDataType,
				OutputType: NilJobDataType,
			},
			{
				Name: "job_3",
				Execute: func(store Store, tx *gorm.DB, deployment *Deployment, value interface{}) (interface{}, error) {
					return value, nil
				},
				InputType:  NilJobDataType,
				OutputType: GetJobDataType(&deploymentTestStruct{}),
			},
		},
	},

	// Helper functions
	getDeploymentJobDataCount: func(t *testing.T, db *gorm.DB) int {
		var jobDataCount int
		require.NoError(t, db.Model(&deploymentData{}).Count(&jobDataCount).Error)

		return jobDataCount
	},
}

func TestNewDeploymentAndGetDeployment(t *testing.T) {
	tr := setupTest(t)

	// New Deployment
	deployment, err := newDeployment(tr.db, deploymentTestData.action)
	require.NoError(t, err)
	require.NotNil(t, deployment)
	require.NotNil(t, deployment.UUID)

	// Get Deployment
	dbDeployment, err := getDeployment(tr.db, &deployment.UUID)
	assert.NoError(t, err)
	assert.NotNil(t, dbDeployment)
	// Make CreatedAt and UpdatedAt fields equal to compare other fields
	dbDeployment.CreatedAt = deployment.CreatedAt
	dbDeployment.UpdatedAt = deployment.UpdatedAt
	assert.Equal(t, deployment, dbDeployment)
}

func TestGetRunningDeployments(t *testing.T) {
	tr := setupTest(t)

	createTestDeployment := func(db *gorm.DB, status DeploymentStatus) error {
		deployment := &Deployment{
			Status: status,
		}
		return db.Create(deployment).Error
	}

	// Create the test deployments
	assert.NoError(t, createTestDeployment(tr.db, deploymentStatusFinished))
	assert.NoError(t, createTestDeployment(tr.db, deploymentStatusRunning))
	assert.NoError(t, createTestDeployment(tr.db, deploymentStatusRollback))
	assert.NoError(t, createTestDeployment(tr.db, deploymentStatusRunning))
	assert.NoError(t, createTestDeployment(tr.db, deploymentStatusRunning))
	assert.NoError(t, createTestDeployment(tr.db, deploymentStatusFinished))
	assert.NoError(t, createTestDeployment(tr.db, deploymentStatusRollback))

	deployments, err := GetRunningDeployments(tr.db)
	assert.NoError(t, err)
	assert.Len(t, deployments, 3)
	// Check that all the deployments are not finished
	for _, deployment := range deployments {
		assert.Equal(t, deploymentStatusRunning, deployment.Status)
	}
}

func TestSetJob(t *testing.T) {
	tr := setupTest(t)
	td := getTestData(t)
	dtd := deploymentTestData

	// Reset the job data type registry
	jobDataTypeRegistry = newDataTypeRegistry()

	deployment, err := newDeployment(tr.db, td.action)
	require.NoError(t, err)
	require.Equal(t, td.jobName1, deployment.CurrentJob)

	// Check that there are no job data entries
	require.Equal(t, 0, dtd.getDeploymentJobDataCount(t, tr.db))

	// Prepare the job data
	type deploymentTestStruct struct {
		I int
	}
	testInputData := &deploymentTestStruct{I: 1}
	jobDataTypeRegistry.register(GetJobDataType(testInputData))

	// Update the deployment's job
	require.NoError(t, deployment.setJob(tr.db, td.jobName2, testInputData))
	require.Equal(t, td.jobName2, deployment.CurrentJob)

	// There should be a deployment job data
	require.Equal(t, 1, dtd.getDeploymentJobDataCount(t, tr.db))

	// Check that the job data is of input type
	out, err := deployment.GetJobData(tr.db, &td.jobName2, deploymentJobInput)
	require.NoError(t, err)
	require.Equal(t, testInputData.I, out.(*deploymentTestStruct).I)
}

func TestSetAndGetJobData(t *testing.T) {
	tr := setupTest(t)
	td := getTestData(t)

	// Reset the job data type registry
	jobDataTypeRegistry = newDataTypeRegistry()

	deployment, err := newDeployment(tr.db, td.action)
	require.NoError(t, err)

	// Prepare the job data
	testInputData := deploymentTestStruct{I: 1}
	jobDataTypeRegistry.register(GetJobDataType(testInputData))
	testJobData := deploymentTestStruct{I: 2}
	jobDataTypeRegistry.register(GetJobDataType(testJobData))

	// Set job data
	require.NoError(t, deployment.SetJobData(tr.db, &td.jobName1, deploymentJobInput, testInputData))
	require.NoError(t, deployment.SetJobData(tr.db, &td.jobName1, deploymentJobData, testJobData))

	// Modify job data
	testInputData = deploymentTestStruct{I: 100}
	require.NoError(t, deployment.SetJobData(tr.db, &td.jobName1, deploymentJobInput, testInputData))

	// Get job data
	compareTestData := func(job string, dataType deploymentDataType, expected deploymentTestStruct) {
		out, err := deployment.GetJobData(tr.db, &job, dataType)
		require.NoError(t, err)
		require.Equal(t, expected, out.(deploymentTestStruct))
	}
	compareTestData(td.jobName1, deploymentJobInput, testInputData)
	compareTestData(td.jobName1, deploymentJobData, testJobData)
}

func TestAddAndGetErrors(t *testing.T) {
	tr := setupTest(t)
	td := getTestData(t)

	deployment, err := newDeployment(tr.db, td.action)
	require.NoError(t, err)

	// Add the errors to the deployment
	testErrors := []error{
		errors.New("1"),
		errors.New("2"),
		errors.New("3"),
	}
	require.NoError(t, deployment.addJobError(tr.db, &td.jobName1, testErrors[0]))
	require.NoError(t, deployment.addJobError(tr.db, &td.jobName2, testErrors[1]))
	require.NoError(t, deployment.addJobError(tr.db, &td.jobName3, testErrors[2]))

	compareTestData := func(job *string, expected []error) {
		dbErrors, err := deployment.GetErrors(tr.db, job)
		require.NoError(t, err)
		require.Equal(t, len(expected), len(dbErrors))
		for i := 0; i < len(expected); i++ {
			require.Equal(t, expected[i].Error(), *dbErrors[i].Error)
		}
	}
	// Get errors for all jobs
	compareTestData(nil, testErrors)
	// Get errors for a specific job
	compareTestData(&td.jobName1, []error{testErrors[0]})
}

func TestSetStatusAndIsStatusAndGetRollbackError(t *testing.T) {
	tr := setupTest(t)
	td := getTestData(t)

	// The default status should be Running
	deployment, err := newDeployment(tr.db, td.action)
	require.NoError(t, err)
	require.Equal(t, deploymentStatusRunning, deployment.Status)

	dbDeployment, err := getDeployment(tr.db, &deployment.UUID)
	require.NoError(t, err)
	require.True(t, dbDeployment.isRunning())

	// Update the deployment status to Finished
	require.NoError(t, deployment.setFinishedStatus(tr.db))
	require.Equal(t, deploymentStatusFinished, deployment.Status)

	dbDeployment, err = getDeployment(tr.db, &deployment.UUID)
	require.NoError(t, err)
	require.True(t, dbDeployment.isFinished())

	// Update the deployment status to Rollback
	rollbackErr := errors.New("rollback")
	require.NoError(t, deployment.setRollbackStatus(tr.db, rollbackErr))
	require.Equal(t, deploymentStatusRollback, deployment.Status)
	require.Equal(t, rollbackErr.Error(), *deployment.RollbackError)

	dbDeployment, err = getDeployment(tr.db, &deployment.UUID)
	require.NoError(t, err)
	require.True(t, dbDeployment.isRollingBack())
	require.Equal(t, rollbackErr, dbDeployment.getRollbackError())
}
