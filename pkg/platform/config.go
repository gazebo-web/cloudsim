package platform

import (
	"context"
	"flag"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/auth0"
	"gitlab.com/ignitionrobotics/web/ign-go"
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
	isGoTest                bool
	logger                  ign.Logger
	logCtx                  context.Context
	Auth0					auth0.Auth0
	// From aws go documentation:
	// Sessions should be cached when possible, because creating a new Session
	// will load all configuration values from the environment, and config files
	// each time the Session is created. Sharing the Session value across all of
	// your service clients will ensure the configuration is loaded the fewest
	// number of times possible.
	awsSession *session.Session
	// Are we using S3 for logs?
	S3LogsCopyEnabled bool `env:"AWS_GZ_LOGS_ENABLED" envDefault:"true"`
}

func NewConfig() Config {
	cfg := Config{}
	cfg.isGoTest = flag.Lookup("test.v") != nil
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file. %+v\n", err)
	}
	if cfg.isGoTest {
		// allow the testing environment to define variables if not yet defined.
		godotenv.Load(".env.testing")
	}

	// Also using env-to-struct approach to read configuration
	if err := env.Parse(&cfg); err != nil {
		// This is a log.Fatal because ign.Logger is not setup yet
		log.Fatalf("Error parsing environment into appConfig struct. %+v\n", err)
	}
	return cfg
}