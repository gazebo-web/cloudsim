package fake

import (
	"github.com/stretchr/testify/mock"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/simulations"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/store"
	"time"
)

// Fake is a fake store.Store implementation.
type Fake struct {
	machines     store.Machines
	orchestrator store.Orchestrator
	ignition     store.Ignition
	mole         store.Mole
}

// Machines mocks the Machine store.
func (f *Fake) Machines() store.Machines {
	return f.machines
}

// Orchestrator mocks the Orchestrator store.
func (f *Fake) Orchestrator() store.Orchestrator {
	return f.orchestrator
}

// Ignition mocks the Ignition store.
func (f *Fake) Ignition() store.Ignition {
	return f.ignition
}

// Mole mocks the Mole store.
func (f *Fake) Mole() store.Mole {
	return f.mole
}

// NewFakeStore initializes a new fake store implementation using fake configuration providers.
// This provider uses the mock library
func NewFakeStore(machines *Machines, orchestrator *Orchestrator, ignition *Ign, mole *Mole) *Fake {
	return &Fake{
		machines:     machines,
		orchestrator: orchestrator,
		ignition:     ignition,
		mole:         mole,
	}
}

// NewDefaultFakeStore initializes a new fake store implementation using default fake configuration providers.
// This provider uses the mock library
func NewDefaultFakeStore() *Fake {
	return &Fake{
		machines:     NewFakeMachines(),
		orchestrator: NewFakeOrchestrator(),
		ignition:     NewFakeIgnition(),
		mole:         NewFakeMole(),
	}
}

// Ign is a fake store.Ignition implementation.
type Ign struct {
	*mock.Mock
}

var _ store.Ignition = (*Ign)(nil)

// DefaultRecipients mocks the DefaultRecipients method.
func (f *Ign) DefaultRecipients() []string {
	args := f.Called()
	return args.Get(0).([]string)
}

// DefaultSender mocks the DefaultSender method.
func (f *Ign) DefaultSender() string {
	args := f.Called()
	return args.String(0)
}

// LogsBucket mocks the LogsBucket method.
func (f *Ign) LogsBucket() string {
	args := f.Called()
	return args.String(0)
}

// GetWebsocketPath mocks the GetWebsocketPath method.
func (f *Ign) GetWebsocketPath(groupID simulations.GroupID) string {
	args := f.Called(groupID)
	return args.String(0)
}

// AccessKeyLabel mocks the AccessKeyLabel method.
func (f *Ign) AccessKeyLabel() string {
	args := f.Called()
	return args.String(0)
}

// SecretAccessKeyLabel mocks the SecretAccessKeyLabel method.
func (f *Ign) SecretAccessKeyLabel() string {
	args := f.Called()
	return args.String(0)
}

// LogsCopyEnabled mocks the LogsCopyEnabled method.
func (f *Ign) LogsCopyEnabled() bool {
	args := f.Called()
	return args.Bool(0)
}

// Region mocks the Region method.
func (f *Ign) Region() string {
	args := f.Called()
	return args.String(0)
}

// SecretsName mocks the SecretsName method.
func (f *Ign) SecretsName() string {
	args := f.Called()
	return args.String(0)
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

// Mole is a fake store.Mole implementation.
type Mole struct {
	*mock.Mock
}

var _ store.Mole = (*Mole)(nil)

// BridgePulsarAddress returns the address of the Pulsar service the Mole bridge should connect to.
func (m *Mole) BridgePulsarAddress() string {
	args := m.Called()
	return args.String(0)
}

// BridgePulsarPort returns the port on which the Pulsar service the mole bridge should connect to is running.
func (m *Mole) BridgePulsarPort() int {
	args := m.Called()
	return args.Int(0)
}

// BridgePulsarHTTPPort returns the port on which the HTTP service the mole bridge should connect to is running.
func (m *Mole) BridgePulsarHTTPPort() int {
	args := m.Called()
	return args.Int(0)
}

// BridgeTopicRegex returns the regex used by the Mole bridge to filter topics.
func (m *Mole) BridgeTopicRegex() string {
	args := m.Called()
	return args.String(0)
}

// NewFakeMole initializes a fake store.Mole implementation.
func NewFakeMole() *Mole {
	return &Mole{
		Mock: new(mock.Mock),
	}
}

// Orchestrator is a fake store.Orchestrator implementation.
type Orchestrator struct {
	*mock.Mock
}

// Timeout mocks the Timeout method.
func (f Orchestrator) Timeout() time.Duration {
	args := f.Called()
	return args.Get(0).(time.Duration)
}

// PollFrequency mocks the PollFrequency method.
func (f Orchestrator) PollFrequency() time.Duration {
	args := f.Called()
	return args.Get(0).(time.Duration)
}

// IngressNamespace mocks the IngressNamespace method.
func (f Orchestrator) IngressNamespace() string {
	args := f.Called()
	return args.String(0)
}

// IngressName mocks the IngressName method.
func (f Orchestrator) IngressName() string {
	args := f.Called()
	return args.String(0)
}

// IngressHost mocks the IngressHost method.
func (f Orchestrator) IngressHost() string {
	args := f.Called()
	return args.String(0)
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

// ClusterName mocks the ClusterName method.
func (f Machines) ClusterName() string {
	args := f.Called()
	return args.String(0)
}

// NewFakeMachines initializes a new store.Machines implementation.
func NewFakeMachines() *Machines {
	return &Machines{
		Mock: new(mock.Mock),
	}
}
