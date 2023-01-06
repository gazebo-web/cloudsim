package actions

import (
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var executeInputTestData = struct {
	// Helper functions
	createTestDeployment func(job string) *Deployment
	getDeploymentCount   func(t *testing.T, db *gorm.DB) int
}{
	// Helper functions
	createTestDeployment: func(job string) *Deployment {
		return &Deployment{
			Action:     "test_action",
			CurrentJob: job,
		}
	},

	getDeploymentCount: func(t *testing.T, db *gorm.DB) int {
		var deploymentCount int
		require.NoError(t, db.Model(&Deployment{}).Count(&deploymentCount).Error)

		return deploymentCount
	},
}

type TestExecuteInput struct {
	*ExecuteInput
}

func TestNewExecuteInput(t *testing.T) {
	tr := setupTest(t)
	defer tr.db.Close()

	td := getTestData(t)

	var input *ExecuteInput
	var err error
	input, err = newExecuteInput(tr.db, td.action)
	require.NoError(t, err)
	// Check that a deployment was created
	require.NotNil(t, input.Deployment)
}

// TestGetExecuteInput checks that any struct that composes ExecuteInput automatically implements ExecuteInputer
func TestGetExecuteInput(t *testing.T) {
	executeInput := TestExecuteInput{
		ExecuteInput: &ExecuteInput{},
	}

	assert.Implements(t, (*ExecuteInputer)(nil), executeInput)
	assert.IsType(t, &ExecuteInput{}, executeInput.getExecuteInput())
}

// TestInitializeNew tests that a new ExecuteInput (with no deployment) is initialized properly
func TestInitializeNew(t *testing.T) {
	tr := setupTest(t)
	defer tr.db.Close()

	td := getTestData(t)
	eitd := executeInputTestData

	input := &ExecuteInput{}
	require.Nil(t, input.Deployment)

	// Get the total number of deployment in the database
	deploymentCount := eitd.getDeploymentCount(t, tr.db)

	// Initialize the input
	require.NoError(t, input.initialize(tr.db, td.action))

	// Check that a deployment was created when initializing the input
	require.NotNil(t, input.Deployment)
	require.NotNil(t, input.Deployment.UUID)

	// Check that the total number of deployments in the database has increased
	require.Equal(t, deploymentCount+1, eitd.getDeploymentCount(t, tr.db))
}

// TestInitializeNew tests that a new ExecuteInput with a deployment is initialized properly
func TestInitializeWithDeployment(t *testing.T) {
	tr := setupTest(t)
	defer tr.db.Close()

	td := getTestData(t)
	eitd := executeInputTestData

	input := &ExecuteInput{
		Deployment: eitd.createTestDeployment(td.jobName2),
	}
	require.NotNil(t, input.Deployment)

	// Get the total number of deployment in the database
	deploymentCount := eitd.getDeploymentCount(t, tr.db)

	// Initialize the input
	require.NoError(t, input.initialize(tr.db, td.action))

	// Check that the total number of deployments in the database has not increased
	require.Equal(t, deploymentCount, eitd.getDeploymentCount(t, tr.db))
}

func TestRestore(t *testing.T) {
	td := getTestData(t)
	eitd := executeInputTestData

	executeInput := TestExecuteInput{
		ExecuteInput: &ExecuteInput{
			Deployment: eitd.createTestDeployment(td.jobName2),
		},
	}
	assert.NoError(t, executeInput.restore(td.action))

	// Since the deployment was in its second job, the index should be 1
	require.Equal(t, 1, executeInput.index)
}

func TestRestoreFailsWhenNoAction(t *testing.T) {
	executeInput := TestExecuteInput{
		ExecuteInput: &ExecuteInput{},
	}
	require.Equal(t, ErrExecuteInputRestoreNoAction, executeInput.restore(nil))
}

func TestRestoreFailsWhenNoDeployment(t *testing.T) {
	td := getTestData(t)

	executeInput := TestExecuteInput{
		ExecuteInput: &ExecuteInput{},
	}
	require.Equal(t, ErrExecuteInputRestoreNoDeployment, executeInput.restore(td.action))
}

func TestExecuteInputRestoreIndex(t *testing.T) {
	td := getTestData(t)
	eitd := executeInputTestData
	input := ExecuteInput{Deployment: eitd.createTestDeployment(td.jobName2)}

	// Check that the index is restored with no error
	assert.NoError(t, input.restoreIndex(td.action))
	// The index should be 1 to match the second job in the slice
	assert.Equal(t, 1, input.index)
}

func TestExecuteInputRestoreIndexFailsWhenNoDeployment(t *testing.T) {
	td := getTestData(t)
	eitd := executeInputTestData

	invalidJobName := "invalid"
	testExecuteInput := TestExecuteInput{
		ExecuteInput: &ExecuteInput{Deployment: eitd.createTestDeployment(invalidJobName)},
	}
	assert.Error(t, testExecuteInput.restoreIndex(td.action))
}

// TestModifyExecuteInput checks that the values of the composed ExecuteInput is modified
func TestModifyExecuteInput(t *testing.T) {
	executeInput := TestExecuteInput{
		ExecuteInput: &ExecuteInput{},
	}

	assert.Equal(t, 0, executeInput.index)

	func(input ExecuteInputer) {
		executeInput := input.getExecuteInput()
		executeInput.index++
	}(
		executeInput,
	)

	assert.Equal(t, 1, executeInput.index)
}
