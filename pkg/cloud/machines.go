package cloud

// CreateMachinesInput is the input for the Machines.Create operation.
// It will be used to create a certain number of machines.
type CreateMachinesInput struct {
	DryRun        bool
	KeyName       string
	MinCount      int64
	MaxCount      int64
	FirewallRules []string
	SubnetId      string
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

// Machines groups a set of methods to Create, Terminate and Count cloud machines.
type Machines interface {
	Create(input CreateMachinesInput) error
	Terminate(input TerminateMachinesInput) error
	Count(input CountMachinesInput) int
}
