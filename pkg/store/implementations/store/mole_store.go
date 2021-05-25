package store

import (
	"github.com/caarlos0/env"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/defaults"
	storepkg "gitlab.com/ignitionrobotics/web/cloudsim/pkg/store"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/validate"
)

// moleStore is a store.Mole implementation.
type moleStore struct {
	// BridgePulsarAddressValue returns the address of the Pulsar service the Mole bridge should connect to.
	BridgePulsarAddressValue string `validate:"required"`

	// BridgePulsarPortValue returns the port on which the Pulsar service the mole bridge should connect to is running.
	BridgePulsarPortValue int `default:"6650"`

	// BridgePulsarHTTPPortValue returns the port on which the HTTP service the mole bridge should connect to is running.
	BridgePulsarHTTPPortValue int `default:"8080"`

	// BridgeTopicRegexValue returns the regex used by the Mole bridge to filter topics.
	BridgeTopicRegexValue string `default:"^subt/"`
}

// SetDefaults sets default values for the store.
func (m *moleStore) SetDefaults() error {
	return defaults.SetStructValues(m)
}

// LogsBucket returns the bucket to upload simulation logs to.
func (m *moleStore) BridgePulsarAddress() string {
	return m.BridgePulsarAddressValue
}

// LogsBucket returns the bucket to upload simulation logs to.
func (m *moleStore) BridgePulsarPort() int {
	return m.BridgePulsarPortValue
}

// LogsBucket returns the bucket to upload simulation logs to.
func (m *moleStore) BridgePulsarHTTPPort() int {
	return m.BridgePulsarHTTPPortValue
}

// BridgeTopicRegex returns the regex used by the Mole bridge to filter topics.
func (m *moleStore) BridgeTopicRegex() string {
	return m.BridgeTopicRegexValue
}

// newIgnitionStoreFromEnvVars initializes a new store.Mole implementation using environment variables.
func newMoleStoreFromEnvVars() (storepkg.Mole, error) {
	// Load store from env vars
	var m moleStore
	if err := env.Parse(&m); err != nil {
		return nil, err
	}
	// Set default values
	if err := defaults.SetValues(&m); err != nil {
		return nil, err
	}
	// Validate values
	if err := validate.Validate(m); err != nil {
		return nil, err
	}

	return &m, nil
}
