package auth0

import "github.com/caarlos0/env"

type Auth0 struct {
	PublicKey   string `env:"AUTH0_RSA256_PUBLIC_KEY"`
}

func New() Auth0 {
	auth := Auth0{}
	if err := env.Parse(&auth); err != nil {
		auth.PublicKey = ""
	}
	return auth
}