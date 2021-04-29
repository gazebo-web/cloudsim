package actions

import (
	"github.com/jinzhu/gorm"
)

// DeploymentError contains a single deployment job error. A job can contain multiple error entries.
type DeploymentError struct {
	gorm.Model
	DeploymentID int
	Deployment   *Deployment
	Job          *string
	Error        *string `gorm:"type:text"`
}

// TableName sets the database table name for DeploymentError.
func (de DeploymentError) TableName() string {
	return "action_deployments_errors"
}

// DeploymentErrors is a slice of DeploymentError pointers.
type DeploymentErrors []*DeploymentError

// newDeploymentError creates a new DeploymentError entry in persistent storage and returns a pointer to it.
// If `job` is nil, the current job of the deployment is used
func newDeploymentError(tx *gorm.DB, deployment *Deployment, job *string, err error) (*DeploymentError, error) {
	if job == nil {
		job = &deployment.CurrentJob
	}

	errMsg := err.Error()
	deploymentErr := &DeploymentError{
		Deployment: deployment,
		Job:        job,
		Error:      &errMsg,
	}

	// Create the persistent storage entry
	if err := tx.Create(deploymentErr).Error; err != nil {
		return nil, err
	}

	return deploymentErr, nil
}

// getDeploymentErrors returns a slice of DeploymentError entries for a Deployment found in persistent storage.
// `job` is not nil, this will return the errors for a single job,
// otherwise this will return errors for all jobs.
func getDeploymentErrors(tx *gorm.DB, deployment *Deployment, job *string) (DeploymentErrors, error) {
	var deploymentErrs DeploymentErrors

	query := tx.Where("deployment_id = ?", deployment.ID)

	// Optionally filter by job
	if job != nil {
		query = query.Where("job = ?", *job)
	}

	err := query.Find(&deploymentErrs).Error
	if err != nil {
		return nil, err
	}

	return deploymentErrs, nil
}
