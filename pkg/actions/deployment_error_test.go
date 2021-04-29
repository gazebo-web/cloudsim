package actions

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"testing"
)

var deploymentErrorTestData = struct {
	// Helper functions
	newDeploymentError func(t *testing.T, db *gorm.DB, deployment *Deployment, job *string,
		err error) *DeploymentError
	getDeploymentErrorCount func(t *testing.T, db *gorm.DB, deployment *Deployment) int
}{
	// Helper functions
	newDeploymentError: func(t *testing.T, db *gorm.DB, deployment *Deployment, job *string,
		err error) *DeploymentError {
		deploymentErr, err := newDeploymentError(db, deployment, job, err)
		require.NoError(t, err)

		return deploymentErr
	},

	getDeploymentErrorCount: func(t *testing.T, db *gorm.DB, deployment *Deployment) int {
		var deploymentErrCount int
		err := db.
			Model(&DeploymentError{}).
			Where("deployment_id = ?", deployment.ID).
			Count(&deploymentErrCount).Error
		require.NoError(t, err)

		return deploymentErrCount
	},
}

func TestNewDeploymentErrorAndGetDeploymentErrors(t *testing.T) {
	tr := setupTest(t)
	defer tr.db.Close()

	td := getTestData(t)
	detd := deploymentErrorTestData

	// Deployment
	deployment, err := newDeployment(tr.db, td.action, uuid.NewV4().String())
	require.NoError(t, err)

	// Check that there are no errors
	require.Equal(t, 0, detd.getDeploymentErrorCount(t, tr.db, deployment))

	// Create DeploymentErrors
	detd.newDeploymentError(t, tr.db, deployment, nil, errors.New("1A"))
	detd.newDeploymentError(t, tr.db, deployment, &td.jobName2, errors.New("2A"))
	detd.newDeploymentError(t, tr.db, deployment, &td.jobName2, errors.New("2B"))
	detd.newDeploymentError(t, tr.db, deployment, &td.jobName2, errors.New("2C"))
	detd.newDeploymentError(t, tr.db, deployment, &td.jobName3, errors.New("3A"))

	// Check that there are 5 errors total created for the deployment
	deploymentErrs, err := getDeploymentErrors(tr.db, deployment, nil)
	require.NoError(t, err)
	require.Equal(t, 5, len(deploymentErrs))

	// Check that there is 1 error created for the current job of the deployment
	deploymentErrs, err = getDeploymentErrors(tr.db, deployment, &deployment.CurrentJob)
	require.NoError(t, err)
	require.Equal(t, 1, len(deploymentErrs))

	// Check that there are 2 errors created for second job of the deployment
	deploymentErrs, err = getDeploymentErrors(tr.db, deployment, &td.jobName2)
	require.NoError(t, err)
	require.Equal(t, 3, len(deploymentErrs))

	// Check that there are no errors for a new deployment
	newDeployment, err := newDeployment(tr.db, td.action, uuid.NewV4().String())
	require.NoError(t, err)
	newDeploymentErrs, err := getDeploymentErrors(tr.db, newDeployment, &td.jobName1)
	require.NoError(t, err)
	require.Equal(t, 0, len(newDeploymentErrs))

}
