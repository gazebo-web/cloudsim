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
	// ErrUnknown is returned when an unknown error is triggered by a cloud provider.
	ErrUnknown = errors.New("unknown error")
	// ErrInsufficientMachines is returned when creating machines fails because there aren't enough machines.
	ErrInsufficientMachines = errors.New("insufficient machines")
	// ErrRequestsLimitExceeded is returned if the request limit has been reached.
	ErrRequestsLimitExceeded = errors.New("requests limit exceeded")
	// ErrMachineCreationFailed is returned when creating machines fails.
	ErrMachineCreationFailed = errors.New("machine creation failed")
	// ErrMissingMachineNames is returned when no machines ids are provided to be terminated.
	ErrMissingMachineNames = errors.New("missing machine names")
)

// Tag is a group of key-value pairs for a certain resource.
type Tag struct {
	Resource string
	Map      map[string]string
}

// CreateMachinesInput is the input for the Machines.Create operation.
// It will be used to create a certain number of machines.
type CreateMachinesInput struct {
	// ResourceName is a file naming convention used to identify a particular resource
	// In AWS: Amazon Resource Names (ARN).
	ResourceName string

	// KeyName is the SSH key-pair's name that will be used on the created machine.
	KeyName string

	// Type is the name of an instance type.
	// Instances types comprise varying combinations of CPU, memory, storage, and networking capacity
	// and give you the flexibility to choose the appropriate mix of resources for your applications.
	Type string

	// Image is the URL of the image that will be used to launch a machine.
	// In AWS: Amazon Machine Images (AMI).
	Image string

	// MinCount defines the minimum amount of machines that should be created.
	MinCount int64

	// MaxCount defines the maximum amount of machines that should be created.
	MaxCount int64

	// FirewallRules is a group of firewall configurations that will be applied to the machine.
	// In AWS: Security groups.
	FirewallRules []string

	// SubnetID is the ID of the subnet that defines a range of IP addresses.
	SubnetID string

	// Zone is a location inside a datacenter that is isolated from other zones.
	// In AWS: Availability zones.
	Zone string

	// Tags is a group of Tag that is being used to identify a machine.
	Tags []Tag

	// InitScript is the initialization script that will be executed when the machine gets created.
	InitScript string

	// Retries is the max amount of retries that will be executed when running in dry run mode.
	// Suggested value: 10.
	Retries int
}

// CreateMachinesOutput is the output for the Machines.Create operation.
// It will be used to display the machines that were created.
type CreateMachinesOutput struct {
	Instances []string
}

// Length returns the amount of instances that were initialized.
func (c *CreateMachinesOutput) Length() int {
	return len(c.Instances)
}

// ToTerminateMachinesInput converts the content of CreateMachinesOutput into TerminateMachinesInput.
func (c *CreateMachinesOutput) ToTerminateMachinesInput() *TerminateMachinesInput {
	return &TerminateMachinesInput{
		Names: c.Instances,
	}
}

// TerminateMachinesInput is the input for the Machines.Terminate operation.
// It will be used to terminate machines.
type TerminateMachinesInput struct {
	Names []string
}

// CountMachinesInput is the input for the Machines.Count operation.
// It will be used to count the number of machines that match a certain list of filters.
type CountMachinesInput struct {
	MaxResults int
	Filters    map[string][]string
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
