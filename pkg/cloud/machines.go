package cloud

// CreateMachinesInput is the input for the Machines.Create operation.
// It will be used to create a certain number of machines.
type CreateMachinesInput struct {
	DryRun        bool
	KeyName       string
	Type          string
	MinCount      int64
	MaxCount      int64
	FirewallRules []string
	SubnetID      string
	Tags          map[string]string
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
	Create(inputs []CreateMachinesInput) error
	// Terminate terminates a set of cloud machines that match a set of names.
	// The names are automatically created by the cloud provider.
	Terminate(input TerminateMachinesInput) error
	// Count returns the number of cloud machines that match a set of selectors.
	// The selectors should have been defined when creating the machines.
	Count(input CountMachinesInput) int
}
