package actions

import (
	"errors"
	"github.com/jinzhu/gorm"
)

var (
	// ErrExecuteInputRestoreNoAction is returned when an ExecuteInput restore is attempted without providing an action
	ErrExecuteInputRestoreNoAction = errors.New("cannot restore an execute input with no action")
	// ErrExecuteInputRestoreNoDeployment is returned when an ExecuteInput restore is attempted without providing a
	// deployment
	ErrExecuteInputRestoreNoDeployment = errors.New("cannot restore an execute input with no deployment")
)

// ExecuteInputer is an interface that action execution inputs must implement.
// This interface allows having a common set of data shared by all action-specific input types.
type ExecuteInputer interface {
	// getExecuteInput returns a generic execute input for the action to be executed.
	getExecuteInput() *ExecuteInput
	// getDeployment returns the execute input's deployment.
	getDeployment() *Deployment
	// initialize initializes the input.
	initialize(tx *gorm.DB, action *Action) error
	// isNew indicates whether this input is new or was restored from a previous execution
	isNew() bool
}

// ExecuteInput contains the set of fields that all action execution inputs must contain.
// An action's custom execution input can compose a pointer to this struct (i.e. *ExecuteInput) to automatically
// implement the required set of fields.
type ExecuteInput struct {
	// ApplicationName is the name of the application of the action to execute.
	ApplicationName *string
	// ActionName is the name of the action to execute.
	ActionName string
	// Deployment contains an optional action deployment.
	// If this value is not `nil`, then the action will be restarted.
	Deployment *Deployment
	// index contains the current job index.
	index int
}

// newExecuteInput creates a new ExecuteInput for a specific deployment.
func newExecuteInput(tx *gorm.DB, action *Action) (*ExecuteInput, error) {
	input := &ExecuteInput{}

	// Initialize the input
	if err := input.initialize(tx, action); err != nil {
		return nil, err
	}

	return input, nil
}

// getExecuteInput is the default implementation that makes all structs that compose *ExecuteInput automatically
// implement ExecuteInputer.
func (ei *ExecuteInput) getExecuteInput() *ExecuteInput {
	return ei
}

// getDeployment returns the execute input's deployment.
func (ei *ExecuteInput) getDeployment() *Deployment {
	return ei.getExecuteInput().Deployment
}

// isNew indicates whether this input is new or was restored from a previous execution
func (ei *ExecuteInput) isNew() bool {
	return ei.index == 0
}

// initialize initializes this input.
// If the input contains a deployment, this method restores the input to the state of the deployment.
// If not, a new deployment is created for the input.
func (ei *ExecuteInput) initialize(tx *gorm.DB, action *Action) error {
	// If the input is for an existing deployment, restore the state and return
	if err := ei.restore(action); err == nil {
		return nil
	}

	// If the input does not contain a previous deployment, create a new one
	if ei.Deployment == nil {
		deployment, err := newDeployment(tx, action)
		if err != nil {
			return err
		}
		ei.Deployment = deployment
	}

	return nil
}

// restore updates the state of this ExecutionInput to allow the service to continue executing an action after an
// interruption (e.g. a server restart).
func (ei *ExecuteInput) restore(action *Action) error {
	// Ensure that this input contains an action
	if action == nil {
		return ErrExecuteInputRestoreNoAction
	}

	// Ensure that this input contains an existing deployment
	if ei.Deployment == nil {
		return ErrExecuteInputRestoreNoDeployment
	}

	if err := ei.restoreIndex(action); err != nil {
		return err
	}

	return nil
}

// restoreIndex restores the job index of this execute input based on the state of the deployment.
func (ei *ExecuteInput) restoreIndex(action *Action) (err error) {
	ei.index, err = action.getJobIndex(&ei.Deployment.CurrentJob)

	return
}
