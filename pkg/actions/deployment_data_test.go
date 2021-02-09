package actions

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type DeploymentJobDataTestStruct struct {
	I       int
	S       string
	B       bool
	IPtr    *int
	SPtr    *string
	BPtr    *bool
	INilPtr *int
	SNilPtr *string
	BNilPtr *bool
}

var deploymentJobDataTestData = struct {
	// Helper functions
	createTestData  func(i int, s string, b bool, iPtr int, sPtr string, bPtr bool) DeploymentJobDataTestStruct
	getJobDataCount func(t *testing.T, db *gorm.DB, deployment *Deployment) int
	marshallJSON    func(t *testing.T, value interface{}) string
}{
	// Helper functions
	createTestData: func(i int, s string, b bool, iPtr int, sPtr string, bPtr bool) DeploymentJobDataTestStruct {
		return DeploymentJobDataTestStruct{
			I:       i,
			S:       s,
			B:       b,
			IPtr:    &iPtr,
			SPtr:    &sPtr,
			BPtr:    &bPtr,
			INilPtr: nil,
			SNilPtr: nil,
			BNilPtr: nil,
		}
	},

	getJobDataCount: func(t *testing.T, db *gorm.DB, deployment *Deployment) int {
		var count int
		err := db.
			Model(&deploymentData{}).
			Where("deployment_id = ?", deployment.ID).
			Count(&count).Error
		require.NoError(t, err)

		return count
	},

	marshallJSON: func(t *testing.T, value interface{}) string {
		out, err := json.Marshal(value)
		require.NoError(t, err)
		return string(out)
	},
}

func TestValue(t *testing.T) {
	type TestStruct struct {
		I int
		S []string
	}

	test := func(value interface{}) interface{} {
		// Register the types
		jobDataTypeRegistry = newDataTypeRegistry()
		jobDataTypeRegistry.register(GetJobDataType(value))
		require.Len(t, jobDataTypeRegistry, 1)

		// Marshal the data
		dataBytes, err := json.Marshal(value)
		require.NoError(t, err)
		dataStr := string(dataBytes)

		jobData := &deploymentData{
			DataType: GetJobDataTypeName(value),
			Data:     &dataStr,
		}

		// Get the value
		out, err := jobData.Value()
		require.NoError(t, err)

		return out
	}

	// Struct
	testStruct := TestStruct{
		I: 100,
		S: []string{"this", "is", "a", "test"},
	}
	require.Equal(t, testStruct, test(testStruct).(TestStruct))

	// Struct pointer
	testStructPtr := &testStruct
	require.Equal(t, testStructPtr, test(testStructPtr).(*TestStruct))

	// String
	testString := "test"
	require.Equal(t, testString, test(testString).(string))

	// String pointer
	testStringPtr := &testString
	require.Equal(t, testStringPtr, test(testStringPtr).(*string))

	// Nil
	require.Equal(t, nil, test(nil))
}

func TestSetDeploymentJobDataAndGetDeploymentJobData(t *testing.T) {
	tr := setupTest(t)
	defer tr.db.Close()

	td := getTestData(t)
	dsdtd := deploymentJobDataTestData

	// Register job types
	jobDataTypeRegistry = newDataTypeRegistry()
	jobDataTypeRegistry.register(GetJobDataType(DeploymentJobDataTestStruct{}))

	deployment, err := newDeployment(tr.db, td.action)
	require.NoError(t, err)

	// Get total count of entries
	require.Equal(t, 0, dsdtd.getJobDataCount(t, tr.db, deployment))

	// Job data
	testInputData := dsdtd.createTestData(123, "input", false, 321, "inputPtr", true)
	testJobData := dsdtd.createTestData(789, "job", true, 987, "jobPtr", false)

	// Create the job data entries
	require.NoError(t, setDeploymentData(tr.db, deployment, &td.jobName1, DeploymentJobInput, testInputData))
	require.NoError(t, setDeploymentData(tr.db, deployment, &td.jobName1, DeploymentJobData, testJobData))
	require.NoError(t, setDeploymentData(tr.db, deployment, &td.jobName2, DeploymentJobData, nil))
	// Check that two entries have been created
	assert.Equal(t, 3, dsdtd.getJobDataCount(t, tr.db, deployment))

	// Update an existing job data entry
	testInputData = dsdtd.createTestData(111, "modifiedInput", true, 999, "modifiedPtr", false)
	require.NoError(t, setDeploymentData(tr.db, deployment, &td.jobName1, DeploymentJobInput, testInputData))
	// Check that the number of entries remains the same
	assert.Equal(t, 3, dsdtd.getJobDataCount(t, tr.db, deployment))

	// Get the job data from the database
	compareWithDB := func(job string, dataType deploymentDataType, expected interface{}) {
		out, err := getDeploymentData(tr.db, deployment, &job, dataType)
		require.NoError(t, err)
		dbJobData := out.(DeploymentJobDataTestStruct)
		require.Equal(t, dsdtd.marshallJSON(t, expected), dsdtd.marshallJSON(t, dbJobData))
	}
	compareWithDB(td.jobName1, DeploymentJobInput, testInputData)
	compareWithDB(td.jobName1, DeploymentJobData, testJobData)
	// Check that the job with null data returns an error
	_, err = getDeploymentData(tr.db, deployment, &td.jobName2, DeploymentJobData)
	require.Equal(t, ErrDeploymentDataNoData, err)
}
