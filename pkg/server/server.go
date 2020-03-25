package server

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/auth0"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

type Config struct {
	Auth0 auth0.Config
	HTTPport string
	SSLport string
}

func New(config Config) (*ign.Server, error) {
	s, err := ign.Init(config.Auth0.PublicKey, "")
	if err != nil {
		return nil, err
	}
	return s, nil
}
