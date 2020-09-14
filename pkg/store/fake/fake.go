package fake

import (
	"github.com/stretchr/testify/mock"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/cloud"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/store"
	"time"
)

// fake is a fake store.Store implementation.
type fake struct {
	machines     store.Machines
	orchestrator store.Orchestrator
	ignition     store.Ignition
}

// Machines mocks the Machine namespace.
func (f fake) Machines() store.Machines {
	return f.machines
}

// Orchestrator mocks the Orchestrator namespace.
func (f fake) Orchestrator() store.Orchestrator {
	return f.orchestrator
}

// Ignition mocks the Ignition namespace.
func (f fake) Ignition() store.Ignition {
	return f.ignition
}

// NewFakeStore initializes a new fake store implementation using fake configuration providers.
// This provider uses the mock library
func NewFakeStore(machines *Machines, orchestrator *Orchestrator, ignition *Ign) *fake {
	return &fake{
		machines:     machines,
		orchestrator: orchestrator,
		ignition:     ignition,
	}
}

// FakeIgn is a fake store.Ignition implementation.
type Ign struct {
	*mock.Mock
}

// IP mocks the IP method.
func (f *Ign) IP() string {
	args := f.Called()
	return args.String(0)
}

// GazeboServerLogsPath mocks the GazeboServerLogsPath method.
func (f *Ign) GazeboServerLogsPath() string {
	args := f.Called()
	return args.String(0)
}

// ROSLogsPath mocks the ROSLogsPath method.
func (f *Ign) ROSLogsPath() string {
	args := f.Called()
	return args.String(0)
}

// SidecarContainerLogsPath mocks the SidecarContainerLogsPath method.
func (f *Ign) SidecarContainerLogsPath() string {
	args := f.Called()
	return args.String(0)
}

// Verbosity mocks the Verbosity method.
func (f *Ign) Verbosity() string {
	args := f.Called()
	return args.String(0)
}

// NewFakeIgnition initializes a fake store.Ignition implementation.
func NewFakeIgnition() *Ign {
	return &Ign{
		Mock: new(mock.Mock),
	}
}

// Orchestrator is a fake store.Orchestrator implementation.
type Orchestrator struct {
	*mock.Mock
}

// Namespace mocks the Namespace method.
func (f Orchestrator) Namespace() string {
	args := f.Called()
	return args.String(0)
}

// TerminationGracePeriod mocks the TerminationGracePeriod method.
func (f Orchestrator) TerminationGracePeriod() time.Duration {
	args := f.Called()
	return args.Get(0).(time.Duration)
}

// Nameservers mocks the Nameservers method.
func (f Orchestrator) Nameservers() []string {
	args := f.Called()
	return args.Get(0).([]string)
}

// NewFakeOrchestrator initializes a new store.Orchestrator implementation.
func NewFakeOrchestrator() *Orchestrator {
	return &Orchestrator{
		Mock: new(mock.Mock),
	}
}

// Machines is a fake store.Machines implementation.
type Machines struct {
	*mock.Mock
}

// InstanceProfile mocks the InstanceProfile method.
func (f Machines) InstanceProfile() *string {
	args := f.Called()
	result := args.String(0)
	if len(result) == 0 {
		return nil
	}
	return &result
}

// KeyName mocks the KeyName method.
func (f Machines) KeyName() string {
	args := f.Called()
	return args.String(0)
}

// Type mocks the Type method.
func (f Machines) Type() string {
	args := f.Called()
	return args.String(0)
}

// FirewallRules mocks the FirewallRules method.
func (f Machines) FirewallRules() []string {
	args := f.Called()
	return args.Get(0).([]string)
}

// SubnetAndZone mocks the SubnetAndZone method.
func (f Machines) SubnetAndZone() (string, string) {
	args := f.Called()
	return args.String(0), args.String(1)
}

// Tags mocks the Tags method.
func (f Machines) Tags(simulation simulations.Simulation, nodeType string, nameSuffix string) []cloud.Tag {
	args := f.Called(simulation, nodeType, nameSuffix)
	return args.Get(0).([]cloud.Tag)
}

// InitScript mocks the InitScript method.
func (f Machines) InitScript() *string {
	args := f.Called()
	result := args.String(0)
	if len(result) == 0 {
		return nil
	}
	return &result
}

// BaseImage mocks the BaseImage method.
func (f Machines) BaseImage() string {
	args := f.Called()
	return args.String(0)
}

// Timeout mocks the Timeout method.
func (f Machines) Timeout() time.Duration {
	args := f.Called()
	return args.Get(0).(time.Duration)
}

// PollFrequency mocks the PollFrequency method.
func (f Machines) PollFrequency() time.Duration {
	args := f.Called()
	return args.Get(0).(time.Duration)
}

// Limit mocks the Limit method.
func (f Machines) Limit() int {
	args := f.Called()
	return args.Int(0)
}

// NamePrefix mocks the NamePrefix method.
func (f Machines) NamePrefix() string {
	args := f.Called()
	return args.String(0)
}

// NewFakeMachines initializes a new store.Machines implementation.
func NewFakeMachines() *Machines {
	return &Machines{
		Mock: new(mock.Mock),
	}
}
