package actions

import (
  "fmt"
	"errors"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// DeploymentStatus is a possible status for a Deployment.
type DeploymentStatus string

const (
	deploymentStatusRunning  DeploymentStatus = "Running"
	deploymentStatusFinished DeploymentStatus = "Finished"
	deploymentStatusRollback DeploymentStatus = "Rollback"
)

var (
	errRollbackErrIsNil = errors.New("nil rollback error provided when changing status")
)

// Deployment contains persistent state of a single action execution.
// The stored state is used to resume actions in case of interruption (server restarts, etc.).
type Deployment struct {
	gorm.Model
	// UUID contains the unique identifier of this deployment.
	UUID string `gorm:"not null"`
	// Action contains the action this deployment executes.
	Action string `gorm:"not null"`
	// Status contains the status of this deployment.
	Status DeploymentStatus `gorm:"not null"`
	// CurrentJob contains the current action job the deployment is executing.
	CurrentJob string `gorm:"not null"`
	// RollbackError contains the error that triggered the rollback of this deployment.
	RollbackError *string
}

// Deployments is a slice of Deployment pointers.
type Deployments []*Deployment

// TableName sets the database table name for Deployment.
func (Deployment) TableName() string {
	return "action_deployments"
}

// newDeployment creates a new Deployment entry in persistent storage and returns a pointer to it.
func newDeployment(tx *gorm.DB, action *Action) (*Deployment, error) {
	// Create the deployment
	deployment := &Deployment{
		UUID:       uuid.NewV4().String(),
		Action:     action.Name,
		CurrentJob: action.Jobs[0].Name,
		Status:     deploymentStatusRunning,
	}

	// Create the storage record
	// TODO: This should use an interface to allow swapping to different storages in the future.
	if err := tx.Create(deployment).Error; err != nil {
		return nil, err
	}

	return deployment, nil
}

// getDeployment gets the deployment for a given a uuid
func getDeployment(tx *gorm.DB, uuid *string) (*Deployment, error) {
	deployment := &Deployment{}

	err := tx.
		Where("UUID = ?", uuid).
		First(deployment).
		Error
	if err != nil {
		return nil, err
	}

	return deployment, nil
}

// GetRunningDeployments returns the set of deployments that are still running.
func GetRunningDeployments(tx *gorm.DB) (Deployments, error) {
	var deployments Deployments

	err := tx.
		Where("status = ?", deploymentStatusRunning).
		Find(&deployments).
		Error
	if err != nil {
		return nil, err
	}

	return deployments, nil
}

// setJob updates the current job of a deployment and creates an input entry if the value is not nil
func (d *Deployment) setJob(tx *gorm.DB, job string, inputData interface{}) error {
	// Update the current job if it changed
	if d.CurrentJob != job {
		d.CurrentJob = job
		if err := tx.Model(d).Save(d).Error; err != nil {
			return err
		}
	}

	// Set the job's input data
	if inputData != nil {
		if err := d.SetJobData(tx, &job, DeploymentJobInput, inputData); err != nil {
			return err
		}
	}

	return nil
}

// setStatus changes the status of the deployment.
func (d *Deployment) setStatus(tx *gorm.DB, status DeploymentStatus) error {
	// Update the status
	d.Status = status

	return tx.Model(d).Save(*d).Error
}

// SetJobData creates a job data entry of a specific type for a job in this deployment.
func (d *Deployment) SetJobData(tx *gorm.DB, job *string, dataType deploymentDataType, data interface{}) error {
  fmt.Printf("SetJobData\n\n")
	return setDeploymentData(tx, d, job, dataType, data)
}

// GetJobData gets job data entry of a specific type for a job in this deployment.
// The job data is written to `out`. `out` must be a pointer.
func (d *Deployment) GetJobData(tx *gorm.DB, job *string, dataType deploymentDataType) (interface{}, error) {
	return getDeploymentData(tx, d, job, dataType)
}

// addJobError adds a new deployment error entry for a deployment job.
// If `job` is nil, then the current job of the deployment will be used.
func (d *Deployment) addJobError(tx *gorm.DB, job *string, err error) error {
	_, err = newDeploymentError(tx, d, job, err)
	return err
}

// GetErrors returns a slice with all the deployment errors logged for this deployment.
// If `job` is not nil, this will return the errors for a single job, otherwise this will return errors for all jobs.
func (d *Deployment) GetErrors(tx *gorm.DB, job *string) (DeploymentErrors, error) {
	return getDeploymentErrors(tx, d, job)
}

// getRollbackError returns the error that triggered the deployment's rollback.
func (d *Deployment) getRollbackError() error {
	if d.RollbackError == nil {
		return nil
	}

	return errors.New(*d.RollbackError)
}

// setRollbackStatus sets the status of this deployment to Rollback and stores the error that triggered the rollback.
func (d *Deployment) setRollbackStatus(tx *gorm.DB, rollbackError error) error {
	// Make sure rollbackError is defined
	if rollbackError == nil {
		return errRollbackErrIsNil
	}

	// Set the rollbackError field
	errMsg := rollbackError.Error()
	d.RollbackError = &errMsg

	return d.setStatus(tx, deploymentStatusRollback)
}

// setFinishedStatus sets the status of this deployment to Finished.
func (d *Deployment) setFinishedStatus(tx *gorm.DB) error {
	return d.setStatus(tx, deploymentStatusFinished)
}

// isRunning indicates if the deployment is in the Running status.
func (d *Deployment) isRunning() bool {
	return d.Status == deploymentStatusRunning
}

// isRollingBack indicates if the deployment is in the Rollback status.
func (d *Deployment) isRollingBack() bool {
	return d.Status == deploymentStatusRollback
}

// isRunning indicates if the deployment is in the Finished status.
func (d *Deployment) isFinished() bool {
	return d.Status == deploymentStatusFinished
}
