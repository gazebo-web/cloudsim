package cloud

import "errors"

var (
	// ErrMissingKeyName is returned when the key name is missing.
	ErrMissingKeyName = errors.New("missing key name")
	// ErrInvalidMachinesCount is returned when the Min and Max count validation fails.
	ErrInvalidMachinesCount = errors.New("invalid machines count")
	// ErrInvalidSubnetID is returned when the subnet ID provided is invalid.
	ErrInvalidSubnetID = errors.New("invalid subnet")
	// ErrDryRunFailed is returned when a dry run operation fails.
	ErrDryRunFailed = errors.New("dry run failed")
	// ErrUnknown is returned if an invalid errors is returned from AWS.
	ErrUnknown = errors.New("unknown error")
	// ErrInsufficientMachines is returned when creating machines fails because there aren't enough machines.
	ErrInsufficientMachines = errors.New("insufficient machines")
	// ErrRequestsLimitExceeded is returned if the request limit has been reached.
	ErrRequestsLimitExceeded = errors.New("requests limit exceeded")
	// ErrMachineCreationFailed is returned when creating machines fails.
	ErrMachineCreationFailed = errors.New("machine creation failed")
)

// CreateMachinesInput is the input for the Machines.Create operation.
// It will be used to create a certain number of machines.
type CreateMachinesInput struct {
	ResourceName  string
	DryRun        bool
	KeyName       string
	Type          string
	MinCount      int64
	MaxCount      int64
	FirewallRules []string
	SubnetID      string
	Zone          string
	Tags          map[string]map[string]string
	InitScript    string
	Retries       int
}

// CreateMachinesOutput is the output for the Machines.Create operation.
// It will be used to display the machines that were created.
type CreateMachinesOutput struct {
	Instances   []string
	length      int
	isLengthSet bool
}

// Length returns the amount of instances that were initialized.
func (c *CreateMachinesOutput) Length() int {
	if !c.isLengthSet {
		c.length = len(c.Instances)
		c.isLengthSet = true
	}
	return c.length
}

// ToTerminateMachinesInput converts the content of CreateMachinesOutput into TerminateMachinesInput.
func (c *CreateMachinesOutput) ToTerminateMachinesInput(dryRun bool) *TerminateMachinesInput {
	return &TerminateMachinesInput{
		Names:  c.Instances,
		DryRun: dryRun,
	}
}

// TerminateMachinesInput is the input for the Machines.Terminate operation.
// It will be used to terminate machines.
type TerminateMachinesInput struct {
	Names  []string
	DryRun bool
}

// CountMachinesInput is the input for the Machines.Count operation.
// It will be used to count the number of machines that match a certain list of tags.
type CountMachinesInput struct {
	MaxResults int
	Tags       map[string][]string
}

// Machines requests physical instances from a cloud provider on which to deploy applications
type Machines interface {
	// Create creates a set of cloud machines with a certain configuration.
	Create(input []CreateMachinesInput) ([]CreateMachinesOutput, error)
	// Terminate terminates a set of cloud machines that match a set of names.
	// The names are automatically created by the cloud provider.
	Terminate(input TerminateMachinesInput) error
	// Count returns the number of cloud machines that match a set of selectors.
	// The selectors should have been defined when creating the machines.
	Count(input CountMachinesInput) int
}
