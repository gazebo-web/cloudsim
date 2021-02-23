package actions

import (
	"fmt"
	"github.com/imdario/mergo"
	"github.com/jinzhu/gorm"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
)

var (
	// NilJobDataType is used to indicate that a Job Input or Output does not receive or return data.
	NilJobDataType = reflect.TypeOf(struct{}{})
)

// JobFunc is the function signature used by job hooks and Execute function.
type JobFunc func(store Store, tx *gorm.DB, deployment *Deployment, value interface{}) (interface{}, error)

// JobErrorHandler is the job function type called when an error occurs in a job.
type JobErrorHandler func(store Store, tx *gorm.DB, deployment *Deployment, value interface{},
	err error) (interface{}, error)

// JobDataType is used to store Job input and output data types.
type JobDataType reflect.Type

// Job is the base atomic component in an Action. It contains the instructions necessary to perform and rollback a
// single operation. A Job is composed of three parts run in succession: PreHooks, Execute function and PostHooks.
// When ran (by calling Job.Run), a Job goes through a sequence of internal steps:
//
//   1. A number of optional PreHook functions are called to perform validation and preprocessing operations on the
//   	Execute function input data.
//   2. The Execute function is called to perform an operation with the input received from the PreHook functions.
//   3. A number of optional PostHook functions are called to perform validation and postprocessing operations on the
//   	Execute function output data before handing off the data to the next Job in the sequence.
//
// Hooks allow fitting a Job for application-specific scenarios. As an example, a Job that launches cloud instances
// needs to get the list of type of instances it needs to launch. Different applications can use hooks to specify this
// information and use the same execute logic to launch instances specific to their applications.
//
// Each Job should only perform one write or update operation, and it should only be performed in the Execute
// function. If more than one operation is needed, or a hook needs to write or update a persistent entry, consider
// creating additional jobs to handle the operation. Doing so will greatly simplify error and rollback handling.
// To create new jobs for an application, instance this struct. If an existing Job is to be fitted for an
// application, Job.Extend can be used to get an application-specific version of the Job.
//
// If a Job function (PreHook, Execute Funcion, PostHook) returns an error, the error is logged and the execution of
// the sequence of jobs stops and is rolled back. Job functions can handle errors by wrapping functions with error
// handlers by calling WrapErrorHandler. Job functions with error handlers allow logging errors and recovering from
// them. If a function or error-handled function returns an error, this will trigger a rollback of the entire job
// sequence.
//
// A Job may contain an optional RollbackHandler function. The RollbackHandler function is in charge of releasing any
// shared resources claimed (e.g. cloud instances, orchestration resources, etc.) and undoing any changes that may
// impact other operations. Rollback logic should always double check to understand the state of things, as the job
// may be rolled back before any resources are allocated. For example, the rollback logic for a job tasked with
// requesting cloud instances should check that there are machines allocated at all before proceeding to request the
// termination of said machines.
//
// The RollbackHandler function may require context from the Execute function to perform its operations (e.g. instance
// ids, pod ids, etc.). This shared context should be stored by the Execute function by calling deployment.SetJobData
// and creating a `deploymentData` type entry. The shared context can then be retrieved by the RollbackHandler by
// calling deployment.GetJobData.
//
// Jobs contain InputType and OutputType fields. These required fields contain the data types expected of the values
// received and returned by the Job, respectively. They are used to automate marhsalling and unmarshalling from
// persistent storage.
//
// The InputType and OutputType fields values should be obtained by calling GetJobDataType with a value of the type
// required by the Job. If a job does not require any input data, or does not return any input data, the special value
// `NilJobDataType` can be used to indicate such. Note that this does not mean that the job will not receive or return
// data; this information lets an actions.Service instance know how to recover data from the persistent records.
//
// As an illustration, consider the following scenario:
//
//          ┌────┐                 ┌────┐                 ┌────┐
// -string→ │ S1 │ -ExampleStruct→ | S2 | -ExampleStruct→ | S3 |
//          └────┘                 └────┘                 └────┘
//	        I: string              I: Nil                 I: ExampleStruct
//          O: ExampleStruct       O: Nil                 O: Nil
//
// S1 receives a string as input and returns an ExampleStruct instance.
// S2 does not require any data to run, but still receives and returns the ExampleStruct received from S1.
// S3 receives the ExampleStruct from S2 and performs some operations based on it.
//
// Since later jobs may require information from previous jobs, it is very important that all Jobs return a
// value. To enforce this, Actions will fail if no data is returned by a Job, even if the job has indicated that it
// returns no data. If a Job does not require or return any data to operate, it should pass along the data it
// receives from the previous job.
type Job struct {
	// Name contains the name of the Job. This value should be unique as it will be the Job identifier.
	Name string `validate:"required"`
	// PreHooks contains the hooks ran before the Execute function. Used to validate/preprocess job input data.
	PreHooks []JobFunc
	// Execute performs this job's action.
	// It should always return a value (not nil), not doing so will result in an error.
	Execute JobFunc `validate:"required"`
	// PostHooks contains the hooks ran after the Execute function. Used to validate/postprocess job output data.
	PostHooks []JobFunc
	// RollbackHandler contains the rollback function for this job.
	RollbackHandler JobErrorHandler
	// InputType contains the input type this job receives.
	// This should contain the return value of calling `GetJobDataType` with a value of the expected type.
	InputType JobDataType `validate:"required"`
	// OutputType contains the output type this job returns.
	// This should contain the return value of calling `GetJobDataType` with a value of the expected type.
	OutputType JobDataType `validate:"required"`
}

// GetJobDataType returns the data type for a job input or output.
// This is used to properly load job data from the persistent storage.
func GetJobDataType(value interface{}) JobDataType {
	if value == nil {
		return nil
	}

	// If the value is of type reflect.Type, return the value
	if _, ok := value.(reflect.Type); ok {
		return value.(JobDataType)
	}

	return reflect.TypeOf(value)
}

// GetJobDataTypeName returns the data type name for a job input or output.
func GetJobDataTypeName(value interface{}) string {
	dataType := GetJobDataType(value)

	if dataType == nil {
		return "nil"
	}

	return dataType.String()
}

// Extend customizes this job by modifying its hooks and error handlers.
// The extension cannot replace the Execute function. If you need to change the Execute function, you should create a
// new job instead.
func (j *Job) Extend(extension Job) *Job {
	// The extension name should match the
	// Ensure that the Execute function is not changed
	if extension.Execute != nil {
		panic(fmt.Sprintf("extend cannot replace %s execute, create a new job instead", j.Name))
	}
  fmt.Printf("ExtName[%s] JName[%s]\n", extension.Name, j.Name)

	// Make the extended job name the same as the job name
	j.Name = extension.Name

	// Create the extended job
	if err := mergo.Merge(&extension, *j); err != nil {
		panic(fmt.Sprintf("extend for %s failed to merge definitions: %s", j.Name, err.Error()))
	}

	return &extension
}

// registerTypes registers types used by this job in a registry.
func (j *Job) registerTypes(registry dataTypeRegistry) {
	registry.register(GetJobDataType(j.InputType))
	registry.register(GetJobDataType(j.OutputType))
}

// validate checks that all the required fields are in place
func (j *Job) validate() error {
	return validator.New().Struct(j)
}

// Run runs the job. It calls the job's pre-hooks, followed by its Execute method, and finally its post-hooks.
func (j *Job) Run(store Store, tx *gorm.DB, deployment *Deployment, value interface{}) (interface{}, error) {
	var err error
	// Ensure there is an Execute function
	if j.Execute == nil {
		panic(fmt.Sprintf("job %s does not have an execute function defined", j.Name))
	}

	// Job functions should never return `nil`, as other jobs may need the value received from previous jobs.
	// The only case where returning `nil` is acceptable is if the input was `nil` to begin with.
	inputValueIsNil := value == nil

	// Process pre-hooks
	if value, err = j.processHooks(store, tx, deployment, value, &j.PreHooks); err != nil {
		return nil, err
	}

	// Execute job
	if value, err = callJobFunc(j.Execute, store, tx, deployment, value); err != nil {
		return nil, err
	}

	// Process post-hooks
	if value, err = j.processHooks(store, tx, deployment, value, &j.PostHooks); err != nil {
		return nil, err
	}

	// Check that a value other than `nil` was returned.
	if !inputValueIsNil && value == nil {
		return nil, ErrJobNilOutput
	}

	return value, nil
}

// processHooks receives an input value and processes it using a sequence of hook functions.
func (j *Job) processHooks(store Store, tx *gorm.DB, deployment *Deployment, value interface{},
	hooks *[]JobFunc) (interface{}, error) {

	var err error
	for _, hook := range *hooks {
		// Process the values
		if value, err = callJobFunc(hook, store, tx, deployment, value); err != nil {
			return nil, err
		}
	}

	return value, nil
}

// callJobFunc calls a function of type JobFunc and checks that the output is valid.
func callJobFunc(jobFunc JobFunc, store Store, tx *gorm.DB, deployment *Deployment,
	value interface{}) (interface{}, error) {

	// Process the values
	var err error
	if value, err = jobFunc(store, tx, deployment, value); err != nil {
		return nil, err
	}

	return value, nil
}

// Jobs is a slice of Job
type Jobs []*Job

// jobsAreValid checks that all jobs are valid.
func (s *Jobs) jobsAreValid() error {
	for _, job := range *s {
		if err := job.validate(); err != nil {
			return err
		}
	}

	return nil
}

// registerTypes registers types used by this sequence of job in a registry.
func (s *Jobs) registerTypes(registry dataTypeRegistry) {
	for _, job := range *s {
		job.registerTypes(registry)
	}
}

// validate checks that the sequence of jobs is valid.
func (s *Jobs) validate() error {
	return s.jobsAreValid()
}
