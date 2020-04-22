package server

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/auth0"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

// Config is a configuration to create a new Ignition Server.
type Config struct {
	Auth0    auth0.Config
	HTTPport string
	SSLport  string
}

// New initializes a new Ignition Server with the given configuration.
func New(config Config) (*ign.Server, error) {
	s, err := ign.Init(config.Auth0.PublicKey, "")
	if err != nil {
		return nil, err
	}
	return s, nil
}
