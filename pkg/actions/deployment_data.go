package actions

import (
	"encoding/json"
	"errors"
	"github.com/jinzhu/gorm"
	"reflect"
)

// deploymentDataType is the type of data being stored for a job in a deploymentData entry.
type deploymentDataType string

const (
	// DeploymentJobInput entries contain the data used as input for the job.
	DeploymentJobInput = deploymentDataType("input")
	// DeploymentJobData entries contain data stored by a job for future use.
	// This data is used by jobs to handle errors and rollback.
	DeploymentJobData = deploymentDataType("job")
)

var (
	// ErrDeploymentDataNoData is returned when job data was requested by the entry contained no data
	ErrDeploymentDataNoData = errors.New("an entry for the type and job was found, but there is no data")
)

// deploymentData contains data related to an action deployment's job.
// This information is used to give context for debugging, resume an interrupted action (e.g. due to a server restart),
// or to have context when handling errors (e.g. knowing what machines to terminate).
type deploymentData struct {
	gorm.Model
	// Deployment contains a reference to the deployment this data is for
	Deployment *Deployment
	// Deployment contains the ID of the deployment this data is for
	DeploymentID int `gorm:"not null"`
	// Job contains the job this data is for
	Job string `gorm:"not null"`
	// Type contains the type of data stored for the job
	Type deploymentDataType `gorm:"not null"`
	// DataType contains the data type of the value stored.
	// This is used in tandem with a dataTypeRegistry to automatically marshal and unmarshal data from storage.
	DataType string `gorm:"not null"`
	// Data contains the job's data. Marshalling and unmarshalling of this data is performed automatically when using
	// `setDeploymentData` and `getDeploymentData`, respectively.
	Data *string `gorm:"not null;type:longtext"`
}

// Value returns the unmarshalled value stored in this deploymentData instance.
func (dsd *deploymentData) Value() (interface{}, error) {
	// Get the data type
	dataType, err := jobDataTypeRegistry.getType(dsd.DataType)
	if err != nil {
		return nil, err
	}
	if dataType == nil {
		return nil, nil
	}

	// Create a pointer of the expected data type
	out := reflect.New(dataType).Interface()

	// Unmarshall the data field
	if err := json.Unmarshal([]byte(*dsd.Data), out); err != nil {
		return nil, err
	}

	// Return the value of the pointer
	return reflect.ValueOf(out).Elem().Interface(), nil
}

// DeploymentDataSet is a slice of deploymentData pointers.
type DeploymentDataSet []*deploymentData

// TableName sets the database table name for deploymentData
func (dsd deploymentData) TableName() string {
	return "action_deployments_data"
}

// newDeploymentData creates a new deploymentData instance and returns a pointer to it.
func newDeploymentData(deployment *Deployment, job string, deploymentDataType deploymentDataType, dataTypeName string,
	data *string) *deploymentData {
	return &deploymentData{
		Deployment: deployment,
		Job:        job,
		Type:       deploymentDataType,
		DataType:   dataTypeName,
		Data:       data,
	}
}

// setDeploymentData stores deployment data for a job in persistent storage.
// If there is no previous data, a new storage entry will be created.
// If there is previous data, the storage entry will be replaced with the new data.
func setDeploymentData(tx *gorm.DB, deployment *Deployment, job *string, deploymentDataType deploymentDataType,
	data interface{}) error {
	if job == nil {
		job = &deployment.CurrentJob
	}

	// Marshal the data to a string
	var dataBytes []byte
	var err error
	if dataBytes, err = json.Marshal(data); err != nil {
		return err
	}
	dataStr := string(dataBytes)

	// Create or update the storage entry
	dataTypeName := GetJobDataTypeName(data)
	deploymentData := newDeploymentData(deployment, *job, deploymentDataType, dataTypeName, &dataStr)
	err = tx.
		Where("deployment_id = ?", deployment.ID).
		Where("job = ?", *job).
		Where("type = ?", deploymentDataType).
		Assign(*deploymentData).
		FirstOrCreate(deploymentData).
		Error
	if err != nil {
		return err
	}

	return nil
}

// getDeploymentData gets a deployment job deploymentData entry from storage.
// Returns an error if there is no data of the selected type for the deployment job.
// `dataType` defines the deploymentDataType of job data returned.
func getDeploymentData(tx *gorm.DB, deployment *Deployment, job *string,
	deploymentDataType deploymentDataType) (interface{}, error) {
	// If job is nil, use the deployment's current job
	if job == nil {
		job = &deployment.CurrentJob
	}

	// Get the job data database entry
	deploymentData := deploymentData{}
	err := tx.
		Where("deployment_id = ?", deployment.ID).
		Where("job = ?", *job).
		Where("type = ?", deploymentDataType).
		First(&deploymentData).
		Error
	if err != nil {
		return nil, err
	}

	// Return an error indicating no data was stored
	if deploymentData.Data == nil || *deploymentData.Data == "null" {
		return nil, ErrDeploymentDataNoData
	}

	// Get the value from the job data entry
	out, err := deploymentData.Value()
	if err != nil {
		return nil, err
	}

	return out, nil
}
