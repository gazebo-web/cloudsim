package machines

import "github.com/pkg/errors"

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
	// ErrExternalServiceError is returned when the external service returns an internal error out of this component's
	// control.
	ErrExternalServiceError = errors.New("external service error")
	// ErrMachineCreationFailed is returned when creating machines fails.
	ErrMachineCreationFailed = errors.New("machine creation failed")
	// ErrMissingMachineNames is returned when no machines ids are provided to be terminated.
	ErrMissingMachineNames = errors.New("missing machine names")
	// ErrMissingMachineFilters is returned when no tags ids are provided to be terminated.
	ErrMissingMachineFilters = errors.New("missing machine filters")
	// ErrInvalidTerminateRequest is used to return an error when validating a termination machine request fails.
	ErrInvalidTerminateRequest = errors.New("invalid terminate machines request")
	// ErrInvalidClusterID is returned when an invalid cluster id is passed when creating machines.
	ErrInvalidClusterID = errors.New("invalid cluster id")
	// ErrRetryable is used to wrap errors returned by cloud providers to signal that the operations are retryable.
	ErrRetryable = errors.New("retryable error")
)

// WrapRetryableError wraps an error with the ErrRetryable error.
// This is typically done to signal that an API call is retryable.
func WrapRetryableError(err error) error {
	return errors.Wrap(ErrRetryable, err.Error())
}

// ErrorIsRetryable checks that an error is wrapped with the ErrRetryable error.
func ErrorIsRetryable(err error) bool {
	return errors.Is(err, ErrRetryable)
}

// Tag is a group of key-value pairs for a certain resource.
type Tag struct {
	Resource string
	Map      map[string]string
}

// NewTags returns a new tag slice with a single tag.
// resourceType is the Resource of the tag in the slice.
func NewTags(resourceType string) []Tag {
	return []Tag{
		{
			Resource: resourceType,
			Map: map[string]string{},
		},
	}
}

// CreateMachinesInput is the input for the Machines.Create operation.
// It will be used to create a certain number of machines.
type CreateMachinesInput struct {
	// InstanceProfile is used to identify a particular resource.
	// In AWS: Used to assign an AWS IAM profile to EC2 instances EC2 so that they can join the EKS cluster.
	InstanceProfile *string

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
	// If not provided, the machines component will provide a subnet automatically.
	SubnetID *string

	// Zone is a location inside a datacenter that is isolated from other zones.
	// If not provided, the machines components will select a zone automatically.
	// In AWS: Availability zones.
	Zone *string

	// Tags is a group of Tag that is being used to identify a machine.
	Tags []Tag

	// InitScript is the initialization script that will be executed when the machine gets created.
	InitScript *string

	// Retries is the max amount of retries that will be executed when running in dry run mode.
	// Suggested value: 10.
	Retries int

	// Labels is the map of labels that will be assigned to the node when it joins the cluster.
	// In AWS, it will be the labels that are assigned to the node in order to join the EKS cluster.
	Labels map[string]string

	// ClusterID identifies the cluster that the nodes should join.
	// In AWS: It's the cluster name.
	ClusterID string
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
func (c *CreateMachinesOutput) ToTerminateMachinesInput() TerminateMachinesInput {
	return TerminateMachinesInput{
		Instances: c.Instances,
	}
}

// ToWaitMachinesOKInput converts the content of CreateMachinesOutput into WaitMachinesOKInput.
func (c *CreateMachinesOutput) ToWaitMachinesOKInput() WaitMachinesOKInput {
	return WaitMachinesOKInput{
		Instances: c.Instances,
	}
}

// TerminateMachinesInput is the input for the Machines.Terminate operation.
// It will be used to terminate machines.
type TerminateMachinesInput struct {
	// Instances has the list of machine ids.
	Instances []string
	// Filters has a list of filters to identify machines that should be deleted.
	Filters map[string][]string
}

// ValidateInstances validates that the instances ids given in the TerminateMachinesInput are valid.
func (in *TerminateMachinesInput) ValidateInstances() error {
	if in.Instances == nil || len(in.Instances) == 0 {
		return ErrMissingMachineNames
	}
	return nil
}

// ValidateFilters validates that the filters given in the TerminateMachinesInput are valid.
func (in *TerminateMachinesInput) ValidateFilters() error {
	if in.Filters == nil || len(in.Filters) == 0 {
		return ErrMissingMachineFilters
	}
	return nil
}

// Validate validates that the TerminateMachinesInput request is valid.
func (in *TerminateMachinesInput) Validate() error {
	if in.ValidateInstances() != nil && in.ValidateFilters() != nil {
		return ErrInvalidTerminateRequest
	}
	return nil
}

// CountMachinesInput is the input for the Machines.Count operation.
// It will be used to count the number of machines that match a certain list of filters.
type CountMachinesInput struct {
	Filters map[string][]string
}

// WaitMachinesOKInput represents a set of machines that need to be waited until they are ready.
type WaitMachinesOKInput struct {
	// Instances has the list of machine ids.
	Instances []string
}

// ListMachinesInput is the input used to list the machines running in a certain cloud provider
type ListMachinesInput struct {
	// Filters are used to filter instances by a certain config.
	//
	// In AWS:
	//    * availability-zone - The Availability Zone of the instance.
	//
	//    * event.code - The code for the scheduled event (instance-reboot | system-reboot
	//    | system-maintenance | instance-retirement | instance-stop).
	//
	//    * event.description - A description of the event.
	//
	//    * event.instance-event-id - The ID of the event whose date and time you
	//    are modifying.
	//
	//    * event.not-after - The latest end time for the scheduled event (for example,
	//    2014-09-15T17:15:20.000Z).
	//
	//    * event.not-before - The earliest start time for the scheduled event (for
	//    example, 2014-09-15T17:15:20.000Z).
	//
	//    * event.not-before-deadline - The deadline for starting the event (for
	//    example, 2014-09-15T17:15:20.000Z).
	//
	//    * instance-state-code - The code for the instance state, as a 16-bit unsigned
	//    integer. The high byte is used for internal purposes and should be ignored.
	//    The low byte is set based on the state represented. The valid values are
	//    0 (pending), 16 (running), 32 (shutting-down), 48 (terminated), 64 (stopping),
	//    and 80 (stopped).
	//
	//    * instance-state-name - The state of the instance (pending | running |
	//    shutting-down | terminated | stopping | stopped).
	//
	//    * instance-status.reachability - Filters on instance status where the
	//    name is reachability (passed | failed | initializing | insufficient-data).
	//
	//    * instance-status.status - The status of the instance (ok | impaired |
	//    initializing | insufficient-data | not-applicable).
	//
	//    * system-status.reachability - Filters on system status where the name
	//    is reachability (passed | failed | initializing | insufficient-data).
	//
	//    * system-status.status - The system status of the instance (ok | impaired
	//    | initializing | insufficient-data | not-applicable).
	Filters map[string][]string
}

// ListMachinesItem represents a single instance listed by the output of Machines.List.
type ListMachinesItem struct {
	// InstanceID is the unique identifier for a single instance.
	InstanceID string
	// State is the state of the instance.
	//
	// In AWS, the state will be any the following values:
	//  - pending
	//  - running
	//  - shutting-down
	//  - stopping
	//  - stopped
	//  - terminated
	State string
}

// ListMachinesOutput is the output value returned by Machines.List. It includes a list of ListMachinesItem.
type ListMachinesOutput struct {
	// Instances represents the actual list of instances returned by Machines.List.
	// Each item has information like the InstanceID and the State.
	Instances []ListMachinesItem
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

	// WaitOK is used to wait for the given machines input to be OK.
	WaitOK(input []WaitMachinesOKInput) error

	// List returns a list of machines based on the given input.
	List(input ListMachinesInput) (*ListMachinesOutput, error)
}
