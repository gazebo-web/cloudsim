package auth0

import "github.com/caarlos0/env"

// Config represents a set of options to configure Auth0.
type Config struct {
	PublicKey string `env:"AUTH0_RSA256_PUBLIC_KEY"`
}

// New returns a new Auth0 Config.
func New() Config {
	auth := Config{}
	if err := env.Parse(&auth); err != nil {
		auth.PublicKey = ""
	}
	return auth
}
