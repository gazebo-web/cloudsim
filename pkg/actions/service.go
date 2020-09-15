package actions

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
)

var (
	// ErrGenerateNoActionName is raised when an application action name generation does not contain an action name.
	ErrGenerateNoActionName = errors.New("tried to generate an application action name with no action name")
	// ErrActionIsNil is raised when a nil action is registered.
	ErrActionIsNil = errors.New("action is nil")
	// ErrActionExists is raised when an action is registered but it already exists.
	ErrActionExists = errors.New("action already exists")
	// ErrActionNotFound is raised when a requested action is not registered in the service.
	ErrActionNotFound = errors.New("action not found")
	// ErrJobNilOutput is raised when a job returns nil
	ErrJobNilOutput = errors.New("job cannot return nil, pass through the input value instead")
	// ErrExecutionStopped is raised when the execution is forcibly stopped by a user command.
	ErrExecutionStopped = errors.New("action execution was stopped by a user command")
)

// Servicer is the interface for action services.
type Servicer interface {
	// RegisterAction registers an action for a specific application.
	RegisterAction(applicationName *string, actionName string, action *Action) error
	// Execute executes an action.
	Execute(store Store, tx *gorm.DB, executeInput ExecuteInputer, jobInput interface{}) error
}

// service provides operations to register and execute actions.
type service struct {
	actions map[string]*Action
}

// NewService returns a pointer to an action Servicer implementation.
func NewService() Servicer {
	service := &service{}
	service.actions = make(map[string]*Action, 0)

	return service
}

// generateApplicationActionName generates the name for an application-specific action.
// actionName cannot be an empty string.
func generateApplicationActionName(applicationName *string, actionName string) (string, error) {
	if actionName == "" {
		return "", ErrGenerateNoActionName
	}

	appName := ""
	if applicationName != nil {
		appName = *applicationName
	}

	return fmt.Sprintf("%s%s", appName, actionName), nil
}

// getAction gets an action based on the application and action names.
// actionName cannot be an empty string.
func (s *service) getAction(applicationName *string, actionName string) (*Action, error) {
	// Get the application-specific action name
	applicationActionName, err := generateApplicationActionName(applicationName, actionName)
	if err != nil {
		return nil, err
	}

	// Get the action
	var action *Action
	var exists bool
	if action, exists = s.actions[applicationActionName]; !exists {
		return nil, ErrActionNotFound
	}

	return action, nil
}

// RegisterAction registers an action for a specific application.
// Actions for an application have the application name prefixed to the action.
// actionName cannot be an empty string.
// action cannot be `nil`.
func (s *service) RegisterAction(applicationName *string, actionName string, action *Action) error {
	// Make sure the action is not nil
	if action == nil {
		return ErrActionIsNil
	}

	// Get the application-specific action name
	applicationActionName, err := generateApplicationActionName(applicationName, actionName)
	if err != nil {
		return err
	}

	// Make sure the action was not registered before
	if _, exists := s.actions[applicationActionName]; exists {
		return ErrActionExists
	}

	// Register the action
	s.actions[applicationActionName] = action

	return nil
}

// Execute executes an action by running each job in the action's job sequence.
// Executing an action includes running an action from scratch, restarting an action (e.g. due to a server restart) and
// handling errors that may come up while running actions.
func (s *service) Execute(store Store, tx *gorm.DB, executeInput ExecuteInputer, jobInput interface{}) (err error) {
	// input contains generic execution information necessary such as the action, the deployment and current job index
	input := executeInput.getExecuteInput()

	// Get the action
	action, err := s.getAction(input.ApplicationName, input.ActionName)
	if err != nil {
		return err
	}

	// Initialize the execution input.
	// This step ensures that the ExecuteInput contains a valid Deployment:
	//   * If a previous deployment for the input exists then it is used.
	//   * If no deployment is found then a new deployment is created and assigned to the input.
	err = executeInput.initialize(tx, action)
	if err != nil {
		return err
	}
	deployment := executeInput.getDeployment()

	// Change the deployment status to Finished after returning
	defer func() {
		// Recover from panics
		if r := recover(); r != nil {
			panicErr := fmt.Errorf("panic: %s", r)
			if err == nil {
				err = panicErr
			}
			if errJobError := deployment.addJobError(tx, nil, panicErr); errJobError != nil {
				err = errJobError
			}
		}

		// Only override the returned error if the Status field fails to update
		if errSetStatus := deployment.setFinishedStatus(tx); errSetStatus != nil {
			err = errSetStatus
		}
	}()

	// Process the sequence of jobs
	if deployment.isRunning() {
		err = s.processJobs(store, tx, action, executeInput, jobInput)
	}

	// Rollback if the deployment has been marked for rollback
	if deployment.isRollingBack() {
		return s.rollback(store, tx, action, executeInput, err)
	}

	return nil
}

// processJobs processes jobs in an action. It is in charge of running each job, validating the output,
// and chaining it with the following job. The jobs processed depend on the `executeInput` state.
// If the `executeInput`'s deployment is not new, this method will only process the current job onwards.
// `jobInput` is also automatically loaded from persistent storage (and overwritten) if the `executeInput`'s
// deployment is not new.
func (s *service) processJobs(store Store, tx *gorm.DB, action *Action, executeInput ExecuteInputer,
	jobInput interface{}) (err error) {
	// input contains generic execution information necessary such as the action, the deployment and current job index
	input := executeInput.getExecuteInput()
	deployment := executeInput.getDeployment()

	// Mark the deployment for rollback if an error was returned
	defer func() {
		if err != nil {
			if errSetStatus := deployment.setRollbackStatus(tx, err); errSetStatus != nil {
				err = errSetStatus
			}
		}
	}()

	// If the executeInput is not new, get the jobInput from persistent storage
	if !executeInput.isNew() {
		if jobInput, err = deployment.GetJobData(tx, &deployment.CurrentJob, deploymentJobInput); err != nil {
			return err
		}
	}

	// Process jobs
	for ; input.index < len(action.Jobs); input.index++ {

		job := action.Jobs[input.index]

		// Update the deployment job
		if err := deployment.setJob(tx, job.Name, jobInput); err != nil {
			return err
		}

		// Run the job
		jobInput, err = job.Run(store, tx, deployment, jobInput)
		// If an error was found, add it to the deployment and return
		if err != nil {
			if err := deployment.addJobError(tx, nil, err); err != nil {
				return err
			}
			return err
		}
	}

	return nil
}

// rollback rolls back an execution, releasing any resources taken (e.g. cloud instances, orchestration resources,
// etc.) and undoing any changes that may affect other executions.
// All error handlers for the current and previous jobs will be executed, to allow them to reset resources.
func (s *service) rollback(store Store, tx *gorm.DB, action *Action, executeInput ExecuteInputer, err error) error {
	// input contains generic execution information necessary such as the action, the deployment and current job index
	input := executeInput.getExecuteInput()
	deployment := executeInput.getDeployment()

	// If err is nil then this rollback is being resumed and the rollback error has to be restored
	if err == nil {
		err = deployment.getRollbackError()
	}

	for ; input.index >= 0; input.index-- {

		job := action.Jobs[input.index]

		// Update the current deployment job
		if err := deployment.setJob(tx, job.Name, nil); err != nil {
			return err
		}

		// Run rollback logic for the current job if defined
		if job.RollbackHandler != nil {
			_, handlerErr := job.RollbackHandler(store, tx, deployment, nil, err)

			// If an error was found, add it to the deployment and return
			if handlerErr != nil {
				rollbackErr := fmt.Errorf("rollback: %s", handlerErr.Error())

				if err := deployment.addJobError(tx, nil, rollbackErr); err != nil {
					return err
				}

				return err
			}
		}
	}

	return err
}
