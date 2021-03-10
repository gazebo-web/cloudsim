package env

import (
	"github.com/caarlos0/env"
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

	// IngressNamespaceValue is the namespace where the gloo ingress lives.
	IngressNamespaceValue string `env:"CLOUDSIM_ORCHESTRATOR_INGRESS_NAMESPACE" envDefault:"default"`

	// IngressNameValue is the name of the Kubernetes Ingress used to route client requests from the Internet to
	// different internal services. This configuration is required to enable websocket connections to simulations.
	IngressNameValue string `env:"SUBT_ORCHESTRATOR_INGRESS_NAME"`

	// IngressHostValue contains the address of the host used to route incoming websocket connections.
	// It is used to select a specific rule to modify in an ingress.
	// The ingress resource referenced by the `IngressName` configuration must contain at least one rule with a host
	// value matching this configuration.
	IngressHostValue string
}

// IngressNamespace returns the ingress namespace.
func (o orchestratorEnvStore) IngressNamespace() string {
	return o.IngressNamespaceValue
}

// IngressName returns the ingress name.
func (o orchestratorEnvStore) IngressName() string {
	return o.IngressNameValue
}

// IngressHost returns the ingress host.
func (o orchestratorEnvStore) IngressHost() string {
	return o.IngressHostValue
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
  var orchEnvStore orchestratorEnvStore
	if err := env.Parse(&orchEnvStore); err != nil {
		panic(err)
	}

	return orchEnvStore
}
