// Package main Ignition Cloudsim Server RESET API
//
// This package provides a REST API to the Ignition CloudSim server.
//
// Schemes: https
// Host: cloudsim.ignitionrobotics.org
// BasePath: /1.0
// Version: 0.1.0
// License: Apache 2.0
// Contact: info@openrobotics.org
//
// swagger:meta
// go:generate swagger generate spec
package main

// Import this file's dependencies
import (
	"gitlab.com/ignitionrobotics/web/fuelserver/permissions"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"gitlab.com/ignitionrobotics/web/cloudsim/globals"
	igntran "gitlab.com/ignitionrobotics/web/cloudsim/ign-transport"
	sim "gitlab.com/ignitionrobotics/web/cloudsim/simulations"
	useracc "gitlab.com/ignitionrobotics/web/cloudsim/users"
	"context"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/caarlos0/env"
	"github.com/go-playground/form"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"
	"k8s.io/client-go/kubernetes"
	"log"
	"strconv"
	"strings"
)

// Impl note: we move this as a constant as it is used by tests.
const sysAdminForTest = "rootfortests"

var sysAdminIdentityForTest string

type appConfig struct {
	SysAdmin            string `env:"IGN_CLOUDSIM_SYSTEM_ADMIN"`
	Auth0RsaPublickey   string `env:"AUTH0_RSA256_PUBLIC_KEY"`
	SSLport             string `env:"IGN_CLOUDSIM_SSL_PORT" envDefault:":4431"`
	HTTPport            string `env:"IGN_CLOUDSIM_HTTP_PORT" envDefault:":8001"`
	LogVerbosity        string `env:"IGN_LOGGER_VERBOSITY"`
	RollbarLogVerbosity string `env:"IGN_LOGGER_ROLLBAR_VERBOSITY"`
	// Identity for the systemAdmin to be used during `go test`.
	SysAdminIdentityForTest string `env:"IGN_SYSTEM_ADMIN_IDENTITY_TEST"`
	ConnectToCloud          bool   `env:"IGN_CLOUDSIM_CONNECT_TO_CLOUD"`
	NodesManagerImpl        string `env:"IGN_CLOUDSIM_NODES_MGR_IMPL" envDefault:"ec2"`
	IgnTransportTopic       string `env:"IGN_TRANSPORT_TEST_TOPIC" envDefault:"/foo"`
	isGoTest                bool
	logger                  ign.Logger
	logCtx                  context.Context
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

// defined here to be able to Stop it during Server shutdown
var ecNm *sim.Ec2Client

/////////////////////////////////////////////////
/// Initialize this package
func init() {

	cfg := appConfig{}
	cfg.isGoTest = flag.Lookup("test.v") != nil

	// Using ENV approach to allow multiple layers of configuration.
	// See https://github.com/joho/godotenv
	// Load the original .env
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

	// override sys admin for tests
	if cfg.isGoTest {
		cfg.SysAdmin = sysAdminForTest
		sysAdminIdentityForTest = cfg.SysAdminIdentityForTest
	}

	logger := initLogging(cfg)
	cfg.logger = logger
	logCtx := ign.NewContextWithLogger(context.Background(), logger)
	cfg.logCtx = logCtx

	// Get the auth0 credentials.
	if cfg.Auth0RsaPublickey == "" {
		logger.Info("Missing AUTH0_RSA256_PUBLIC_KEY env variable. Authentication will not work.")
	}

	var err error

	globals.Server, err = ign.Init(cfg.Auth0RsaPublickey, "")

	if err != nil {
		// Log and shutdown the app , if there is an error during startup
		logger.Critical(err)
		log.Fatalf("Error while initializing app. %+v\n", err)
	}
	// Override ports
	globals.Server.HTTPPort = cfg.HTTPport
	globals.Server.SSLport = cfg.SSLport
	logger.Info(fmt.Sprintf("Using HTTP port [%s] and SSL port [%s]", globals.Server.HTTPPort, globals.Server.SSLport))

	// Create the main Router and set it to the server.
	// Note: here it is the place to define multiple APIs
	s := globals.Server
	mainRouter := ign.NewRouter()
	apiPrefix := "/" + globals.APIVersion
	r := mainRouter.PathPrefix(apiPrefix).Subrouter()
	s.ConfigureRouterWithRoutes(apiPrefix, r, sim.Routes)

	globals.Server.SetRouter(mainRouter)

	globals.DefaultEmailRecipients = sim.EnvVarToSlice("IGN_DEFAULT_EMAIL_RECIPIENT")
	globals.DefaultEmailSender, _ = ign.ReadEnvVar("IGN_DEFAULT_EMAIL_SENDER")
	globals.Validate = initValidator(cfg)
	globals.FormDecoder = form.NewDecoder()
	globals.Permissions = initPermissions(cfg)

	// TODO This should probably be stored in the service configuration
	globals.DisableSummaryEmails = false
	// Set the global configuration to true if the env var is set to true
	if value, err := ign.ReadEnvVar("IGN_DISABLE_SUMMARY_EMAILS"); err == nil && strings.ToLower(value) == "true" {
		globals.DisableSummaryEmails = true
	}
	logger.Info(fmt.Sprintf("Disable summary emails is set to %t.", globals.DisableSummaryEmails))
	// TODO This should probably be stored in the service configuration
	globals.DisableScoreGeneration = false
	// Set the global configuration to true if the env var is set to true
	if value, err := ign.ReadEnvVar("IGN_DISABLE_SCORE_GENERATION"); err == nil && strings.ToLower(value) == "true" {
		globals.DisableScoreGeneration = true
	}
	logger.Info(fmt.Sprintf("Disable score generation is set to %t.", globals.DisableScoreGeneration))

	// Initialize the users proxy
	userAccessor, err := initUserAccessor(logCtx, cfg)
	if err != nil {
		// Log and shutdown the app , if there is an error during startup
		logger.Critical("Critical error trying to create UserAccessor", err)
		log.Fatalf("%+v\n", err)
	}
	globals.UserAccessor = userAccessor

	logger.Info("[application.go] Started using database: " + globals.Server.DbConfig.Name)

	// Migrate database tables
	DBMigrate(logCtx, globals.Server.Db)
	DBAddDefaultData(logCtx, globals.Server.Db)
	// After loading initial data, apply custom indexes. Eg: fulltext indexes
	DBAddCustomIndexes(logCtx, globals.Server.Db)

	sim.HttpHandlerInstance, err = sim.NewHttpHandler(logCtx, globals.UserAccessor)
	if err != nil {
		logger.Critical("Critical error trying to create the HttpHandler", err)
		log.Fatalf("%+v\n", err)
	}

	// Create the default instance of the simulations service
	connectToCloudServices := !cfg.isGoTest && cfg.ConnectToCloud

	// Configure access to kubernetes cluster
	var kcli kubernetes.Interface
	kcli = sim.NewK8Factory(cfg.isGoTest, connectToCloudServices).NewK8(logCtx)
	globals.KClientset = kcli

	// Note: we were always creating the AWS session. And our logic relies on it.
	// TODO: we need to change this to avoid creating the AWS session if using minikube.
	cfg.awsSession = session.Must(session.NewSession())
	awsFactory := sim.NewAWSFactory(cfg.isGoTest)
	ec2Svc := awsFactory.NewEC2Svc(cfg.awsSession)
	globals.EC2Svc = ec2Svc
	s3Svc := awsFactory.NewS3Svc(cfg.awsSession)
	globals.S3Svc = s3Svc

	subT, err := sim.NewSubTApplication(logCtx, s3Svc)
	if err != nil {
		// Log and shutdown the app , if there is an error during startup
		logger.Critical("Critical error trying to create SubT Application", err)
		log.Fatalf("%+v\n", err)
	}

	// Create Node Manager instance
	var nm sim.NodeManager
	if cfg.NodesManagerImpl == "minikube" {
		nm, err = sim.NewLocalNodesClient(cfg.logCtx, kcli)
	} else if cfg.NodesManagerImpl == "ec2" {
		ecNm, err = sim.NewEC2Client(logCtx, kcli, ec2Svc)
		ecNm.RegisterPlatform(logCtx, subT)
		nm = ecNm
	} else {
		err = errors.New("Invalid .env value for NodeManager impl")
	}
	if cfg.isGoTest {
		// We switch the underlying implementation to nil to avoid unexpectedly creating
		// resources.
		sim.AssertMockedEC2(ec2Svc).SetImpl(nil)
	}

	if err != nil {
		cfg.logger.Critical("Critical error trying to create the nodeManager", err)
		log.Fatalf("%+v\n", err)
	}

	// Create the Simulations Service instance
	// First initialize the Jobs Pool factory, if needed
	var pFactory sim.PoolFactory
	if cfg.isGoTest {
		pFactory = sim.SynchronicPoolFactory
	}
	if sim.SimServImpl, err = sim.NewSimulationsService(logCtx, globals.Server.Db, nm, kcli, pFactory, userAccessor); err != nil {
		// Log and shutdown the app , if there is an error during startup
		logger.Critical("Critical error trying to create Simulations services", err)
		log.Fatalf("%+v\n", err)
	}
	sim.SimServImpl.RegisterApplication(logCtx, subT)
	sim.SimServImpl.Start(logCtx)

	// initialize ign-transport, for pub/sub
	initIgnTransport(cfg)
}

/////////////////////////////////////////////////
// initIgnTransport creates a new ignition transport node and subscribes to a test topic
func initIgnTransport(cfg appConfig) {
	// Create a new ignition transport node
	ignNode, _ := igntran.NewIgnTransportNode(nil)
	globals.IgnTransport = ignNode
	globals.IgnTransportTopic = cfg.IgnTransportTopic
}

/////////////////////////////////////////////////
// Run the router and server
func main() {
	// Launch the server
	globals.Server.Run()
	// Destroy Sim Service
	sim.SimServImpl.Stop(context.Background())
	if ecNm != nil {
		ecNm.Stop()
	}
	// Destroy ign_transport node
	globals.IgnTransport.Free()
}

func initValidator(cfg appConfig) *validator.Validate {
	validate := validator.New()
	InstallCustomValidators(validate)
	sim.InstallSimulationCustomValidators(validate)
	sim.InstallSubTCustomValidators(validate)
	return validate
}

func initLogging(cfg appConfig) ign.Logger {
	verbosity := ign.VerbosityWarning
	if cfg.LogVerbosity != "" {
		verbosity, _ = strconv.Atoi(cfg.LogVerbosity)
	}
	rollbarVerbosity := ign.VerbosityWarning
	if cfg.RollbarLogVerbosity != "" {
		rollbarVerbosity, _ = strconv.Atoi(cfg.RollbarLogVerbosity)
	}

	logStd := ign.ReadStdLogEnvVar()
	logger := ign.NewLoggerWithRollbarVerbosity("init", logStd, verbosity, rollbarVerbosity)
	return logger
}

// initializes casbin permissions
func initPermissions(cfg appConfig) *permissions.Permissions {

	// Dev note: we need to have a 'permissions/policy.conf' file in the project.
	// It will be used by the permissions package during initialization.

	if cfg.SysAdmin == "" {
		cfg.logger.Info("No IGN_CLOUDSIM_SYSTEM_ADMIN enivironment variable set. " +
			"No system administrator role will be created")
	}

	p := &permissions.Permissions{}
	p.Init(globals.Server.Db, cfg.SysAdmin)
	return p
}

// initUserAccessor initializes access to Users (from ign-fuel db)
func initUserAccessor(ctx context.Context, cfg appConfig) (useracc.UserAccessor, error) {

	// Dev notes:
	// Users live in the ign-fuel DB and not in the cloudsim DB. That's why we need
	// to read different environment variables to connect to that DB.
	// On the other hand, Permissions (casbin) will live on each specific application
	// (eg. ign-fuel, cloudsim). Each Casbin db will have the permissions specific to
	// its associated application.
	// For this case, Users live in the ign-fuel db, and cloudsim permissions in the
	// cloudsim db.

	type dbConfig struct {
		// See ign.DatabaseConfig for fields documentation
		UserName     string `env:"IGN_USER_DB_USERNAME" envDefault:":notset"`
		Password     string `env:"IGN_USER_DB_PASSWORD"`
		Address      string `env:"IGN_USER_DB_ADDRESS"`
		Name         string `env:"IGN_USER_DB_NAME" envDefault:"usersdb"`
		MaxOpenConns int    `env:"IGN_USER_DB_MAX_OPEN_CONNS" envDefault:"66"`
		EnableLog    bool   `env:"IGN_USER_DB_LOG" envDefault:"false"`
	}

	dbCfg := dbConfig{}
	// Also using env-to-struct approach to read configuration
	if err := env.Parse(&dbCfg); err != nil {
		return nil, errors.Wrap(err, "Error parsing environment into userDB dbConfig struct. %+v\n")
	}

	ignDbCfg := ign.DatabaseConfig{
		UserName:     dbCfg.UserName,
		Password:     dbCfg.Password,
		Address:      dbCfg.Address,
		Name:         dbCfg.Name,
		MaxOpenConns: dbCfg.MaxOpenConns,
		EnableLog:    dbCfg.EnableLog,
	}

	if cfg.isGoTest {
		ignDbCfg.Name = ignDbCfg.Name + "_test"
		// Parse verbose setting, and adjust logging accordingly
		if !flag.Parsed() {
			flag.Parse()
		}
		v := flag.Lookup("test.v")
		isTestVerbose := v.Value.String() == "true"
		if isTestVerbose {
			ignDbCfg.EnableLog = true
		}
	}

	usersDb, err := ign.InitDbWithCfg(&ignDbCfg)
	if err != nil {
		return nil, err
	}

	// \todo(anyone) This is a horrible hack. Here is the back story.
	// Private tokens are used to grant users access to secure routes. Route
	// access is determined in the ign-go package, typically through the use
	// of JWT. Private tokens however require access to the user database, where
	// the tokens are stored. So, web-cloudsim needs to tell ign-go about the
	// user database. And...that is how we got here.
	// In the future, this should be removed and replace with a proper User
	// service.
	globals.Server.UsersDb = usersDb

	ua, err := useracc.NewUserAccessor(ctx, globals.Permissions, usersDb, cfg.SysAdmin)
	if err != nil {
		return nil, err
	}

	if !cfg.isGoTest {
		ua.StartAutoLoadPolicy()
	}

	return ua, nil
}
