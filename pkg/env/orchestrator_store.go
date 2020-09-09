package env

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/store"
	"time"
)

// orchestratorEnvStore is a store.Orchestrator implementation using env vars.
type orchestratorEnvStore struct {
	// TerminationGracePeriodSecondsValue is the amount of time in seconds that a simulation needs to terminate.
	TerminationGracePeriodSecondsValue int `env:"CLOUDSIM_ORCHESTRATOR_TERMINATION_GRACE_SECONDS" envDefault:"120"`

	// NameserverValues is a comma separated list of nameservers that will be used to allow simulations
	// to access the internet to upload logs.
	NameserverValues []string `env:"CLOUDSIM_ORCHESTRATOR_NAMESERVERS" envDefault:"8.8.8.8,1.1.1.1" envSeparator:","`

	// NamespaceValue is the orchestrator namespace where simulations should be launched.
	NamespaceValue string `env:"CLOUDSIM_ORCHESTRATOR_NAMESPACE" envDefault:"default"`
}

// TerminationGracePeriod duration that pods need to terminate gracefully.
func (o orchestratorEnvStore) TerminationGracePeriod() time.Duration {
	return time.Duration(o.TerminationGracePeriodSecondsValue) * time.Second
}

// Nameservers returns a slice of the nameservers used to expose simulations to the internet.
func (o orchestratorEnvStore) Nameservers() []string {
	return o.NameserverValues
}

// Namespace returns the base namespace that should be used for simulations.
func (o orchestratorEnvStore) Namespace() string {
	return o.NamespaceValue
}

// newOrchestratorStore initializes a new store.Orchestrator implementation using orchestratorEnvStore.
func newOrchestratorStore() store.Orchestrator {
	return &orchestratorEnvStore{}
}
