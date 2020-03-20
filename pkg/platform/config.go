package platform

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/auth0"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/email"
	"log"
)

type Config struct {
	SysAdmin            string `env:"IGN_CLOUDSIM_SYSTEM_ADMIN"`
	SSLport             string `env:"IGN_CLOUDSIM_SSL_PORT" envDefault:":4431"`
	HTTPport            string `env:"IGN_CLOUDSIM_HTTP_PORT" envDefault:":8001"`
	// Identity for the systemAdmin to be used during `go test`.
	SysAdminIdentityForTest string `env:"IGN_SYSTEM_ADMIN_IDENTITY_TEST"`
	ConnectToCloud          bool   `env:"IGN_CLOUDSIM_CONNECT_TO_CLOUD"`
	NodesManagerImpl        string `env:"IGN_CLOUDSIM_NODES_MGR_IMPL" envDefault:"ec2"`
	IgnTransportTopic       string `env:"IGN_TRANSPORT_TEST_TOPIC" envDefault:"/foo"`
	Auth0					auth0.Config
	Email 					email.Config
	// Are we using S3 for logs?
	S3LogsCopyEnabled bool `env:"AWS_GZ_LOGS_ENABLED" envDefault:"true"`
}

func NewConfig() Config {
	cfg := Config{}
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file. %+v\n", err)
	}

	// Also using env-to-struct approach to read configuration
	if err := env.Parse(&cfg); err != nil {
		// This is a log.Fatal because ign.Logger is not setup yet
		log.Fatalf("Error parsing environment into Platform config struct. %+v\n", err)
	}

	cfg.Auth0 = auth0.New()
	cfg.Email = email.New()
	return cfg
}

func NewTestConfig() Config {
	cfg := Config{}
	if err := godotenv.Load(".env.testing"); err != nil {
		log.Printf("Error loading .env.testing file. %+v\n", err)
	}

	// Also using env-to-struct approach to read configuration
	if err := env.Parse(&cfg); err != nil {
		// This is a log.Fatal because ign.Logger is not setup yet
		log.Fatalf("Error parsing environment into Platform config struct. %+v\n", err)
	}

	return cfg
}