package auth0

import "github.com/caarlos0/env"

type Config struct {
	PublicKey   string `env:"AUTH0_RSA256_PUBLIC_KEY"`
}

func New() Config {
	auth := Config{}
	if err := env.Parse(&auth); err != nil {
		auth.PublicKey = ""
	}
	return auth
}