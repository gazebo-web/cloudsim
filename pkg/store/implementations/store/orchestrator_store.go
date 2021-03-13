package store

import (
	"github.com/caarlos0/env"
	storepkg "gitlab.com/ignitionrobotics/web/cloudsim/pkg/store"
	"time"
)

// orchestratorStore is a store.Orchestrator implementation using env vars.
type orchestratorStore struct {
	// TerminationGracePeriodSecondsValue is the amount of time in seconds that a simulation needs to terminate.
	TerminationGracePeriodSecondsValue int `default:"120" env:"CLOUDSIM_ORCHESTRATOR_TERMINATION_GRACE_SECONDS"`

	// NameserverValues is a comma separated list of nameservers that will be used to allow simulations
	// to access the internet to upload logs.
	NameserverValues []string `default:"[\"8.8.8.8\",\"1.1.1.1\"]" env:"CLOUDSIM_ORCHESTRATOR_NAMESERVERS" envSeparator:","`

	// NamespaceValue is the orchestrator namespace where simulations should be launched.
	NamespaceValue string `default:"default" env:"CLOUDSIM_ORCHESTRATOR_NAMESPACE"`

	// IngressNamespaceValue is the namespace where the gloo ingress lives.
	IngressNamespaceValue string `default:"default" env:"CLOUDSIM_ORCHESTRATOR_INGRESS_NAMESPACE"`

	// IngressNameValue is the name of the Kubernetes Ingress used to route client requests from the Internet to
	// different internal services. This configuration is required to enable websocket connections to simulations.
	IngressNameValue string `validate:"required" env:"SUBT_ORCHESTRATOR_INGRESS_NAME"`

	// IngressHostValue contains the address of the host used to route incoming websocket connections.
	// It is used to select a specific rule to modify in an ingress.
	// The ingress resource referenced by the `IngressName` configuration must contain at least one rule with a host
	// value matching this configuration.
	IngressHostValue string `validate:"required" env:"CLOUDSIM_ORCHESTRATOR_INGRESS_HOST"`
}

// IngressNamespace returns the ingress namespace.
func (o orchestratorStore) IngressNamespace() string {
	return o.IngressNamespaceValue
}

// IngressName returns the ingress name.
func (o orchestratorStore) IngressName() string {
	return o.IngressNameValue
}

// IngressHost returns the ingress host.
func (o orchestratorStore) IngressHost() string {
	return o.IngressHostValue
}

// TerminationGracePeriod duration that pods need to terminate gracefully.
func (o orchestratorStore) TerminationGracePeriod() time.Duration {
	return time.Duration(o.TerminationGracePeriodSecondsValue) * time.Second
}

// Nameservers returns a slice of the nameservers used to expose simulations to the internet.
func (o orchestratorStore) Nameservers() []string {
	return o.NameserverValues
}

// Namespace returns the base namespace that should be used for simulations.
func (o orchestratorStore) Namespace() string {
	return o.NamespaceValue
}

// newOrchestratorStoreFromEnvVars initializes a new store.Orchestrator implementation using environment variables to
// configure an orchestratorStore object.
func newOrchestratorStoreFromEnvVars() (storepkg.Orchestrator, error) {
	var o orchestratorStore
	if err := env.Parse(&o); err != nil {
		return nil, err
	}
	return &o, nil
}
