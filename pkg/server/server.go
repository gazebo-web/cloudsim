package server

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/auth0"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

var Server *ign.Server

type Config struct {
	Auth0 auth0.Auth0
	HTTPport string
	SSLport string
}

func New(config Config) (*ign.Server, error) {
	s, err := ign.Init(config.Auth0.PublicKey, "")
	if err != nil {
		return nil, err
	}
	Server = s

}
