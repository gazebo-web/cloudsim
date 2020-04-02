package platform

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/auth0"
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
	// Are we using S3 for logs?
	S3LogsCopyEnabled bool `env:"AWS_GZ_LOGS_ENABLED" envDefault:"true"`
	PoolSizeLaunchSim    int `env:"SIMSVC_POOL_LAUNCH_SIM" envDefault:"10"`
	PoolSizeTerminateSim int `env:"SIMSVC_POOL_TERMINATE_SIM" envDefault:"10"`
	PoolSizeErrorHandler int `env:"SIMSVC_POOL_ERROR_HANDLER" envDefault:"20"`
	// Timeout in seconds to wait until a new Pod is ready. Default: 5 minutes.
	PodReadyTimeoutSeconds int `env:"SIMSVC_POD_READY_TIMEOUT_SECONDS" envDefault:"300"`
	// Timeout in seconds to wait until a new Node is ready. Default: 5 minutes.
	NodeReadyTimeoutSeconds int `env:"SIMSVC_NODE_READY_TIMEOUT_SECONDS" envDefault:"300"`
	// The number of live simulations a team can have running in parallel. Zero value means unlimited.
	MaxSimultaneousSimsPerOwner int `env:"SIMSVC_SIMULTANEOUS_SIMS_PER_TEAM" envDefault:"1"`
	// MaxDurationForSimulations is the maximum number of minutes a simulation can run in cloudsim.
	// This is a default value. Specific ApplicationTypes are expected to overwrite this.
	MaxDurationForSimulations int `env:"SIMSVC_SIM_MAX_DURATION_MINUTES" envDefault:"45"`
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