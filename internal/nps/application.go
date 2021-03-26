package nps

// This file contains structures and functions that manage the application.
// This is the primary entry point for the nps application.

import (
	"context"
	"github.com/caarlos0/env"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/actions"
	useracc "gitlab.com/ignitionrobotics/web/cloudsim/pkg/users"
	ignGorm "gitlab.com/ignitionrobotics/web/cloudsim/pkg/utils/db/gorm"
	"gitlab.com/ignitionrobotics/web/fuelserver/permissions"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"gitlab.com/ignitionrobotics/web/ign-go/monitoring/prometheus"
)

const (
	// actionNameStartSimulation is the name used to register the start simulation action.
	actionNameStartSimulation = "start-simulation"
	// actionNameStopSimulation is the name used to register the stop simulation action.
	actionNameStopSimulation = "stop-simulation"

	// applicationName is the name of the current simulator's application.
	applicationName    = "nps"
	applicationVersion = "0.1.0"
)

// Application is an interface designed to manage this application.
type Application interface {
	// GetAPIRoutes returns the routes used by this application. This function
	// is implemented in api_routes.go
	GetAPIRoutes() ign.Routes

	// GetMonitoringRoutes returns the monitoring routes used by this
	// application. This function is implemented in monitoring_routes.go
	GetMonitoringRoutes() ign.Routes

	// Run will execute the application
	Run()
}

// application is the structure that holds application data
type application struct {
	// controller handles route requests
	controller Controller

	// db is a pointer to the gorm database interface.
	db *gorm.DB

	// server is a pointer to the Ignition Go web server.
	server *ign.Server
}

// NewApplication creates and returns a new application.
//
// `apiVersion` is the application version, such as "1.0". Routes are prefixed with this version string.
//
// `logger` is the Ignition Go logger to use for logging information.
//
// An application is returned on success, or nil on error.
func NewApplication(apiVersion string, logger ign.Logger) (Application, error) {

	logger.Debug("Creating ", applicationName, " version ", applicationVersion)

	// This will use the following environment variables to create a database
	// connection
	//     * IGN_DB_ADDRESS Address of the DBMS.
	//     * IGN_DB_USERNAME Username to connect to the DBMS with.
	//     * IGN_DB_PASSWORD Password to connect to the DBMS with.
	//     * IGN_DB_NAME Name of the database to connect to.
	//     * IGN_DB_MAX_OPEN_CONNS - (Optional) You run the risk of getting a
	//                           'too many connections' error if this is not set.
	logger.Debug("Initializing database connection")
	db, err := ignGorm.GetDBFromEnvVars()
	if err != nil {
		return nil, err
	}

	app := &application{
		controller: NewController(db, logger),
		db:         db,
	}

	// Update the database.
	app.db.AutoMigrate(
		&Simulation{},
		&RegisteredUsers{},
		// \todo: This is mandatory in order for jobs.LaunchInstances to work. Without this, the job will fail with "Error 1146: Table 'nps.action_deployments' doesn't exist". This seems like a hidden dependency. How about letting the job/action create the table if it's missing?
		&actions.Deployment{},
	)

	_ = actions.MigrateDB(db)

	// Create a server monitoring provider. Specifying a provider makes the
	// server automatically add middleware required to track metrics.
	monitoring := prometheus.NewPrometheusProvider("")

	app.server, err = ign.Init("", "", monitoring)

	if err := setupUsers(app, logger); err != nil {
		return nil, err
	}

	// Create a router
	logger.Debug("Initializing router")

	router := ign.NewRouter()
	apiPrefix := "/" + apiVersion

	// Create the API routes.
	apiRouter := router.PathPrefix(apiPrefix).Subrouter()
	app.server.ConfigureRouterWithRoutes(apiPrefix, apiRouter, app.GetAPIRoutes())

	// Create the Monitoring routes.
	monitorRouter := router.PathPrefix("/").Subrouter()
	app.server.ConfigureRouterWithRoutes("/", monitorRouter, app.GetMonitoringRoutes())

	// Set router.
	// Because a monitoring provider was set, this call will add monitoring
	// routes as well as setting the router
	app.server.SetRouter(router)

	return app, nil
}

// Run will execute the application
func (app *application) Run() {
	app.server.Run()
}

// SetupUsers connects the application to the user database and
func setupUsers(app *application, logger ign.Logger) error {

	// Read the user databas connection information from environment variables.
	type UserConfig struct {
		// See ign.DatabaseConfig for fields documentation
		UserName     string `env:"IGN_USER_DB_USERNAME" envDefault:":notset"`
		Password     string `env:"IGN_USER_DB_PASSWORD"`
		Address      string `env:"IGN_USER_DB_ADDRESS"`
		Name         string `env:"IGN_USER_DB_NAME" envDefault:"usersdb"`
		MaxOpenConns int    `env:"IGN_USER_DB_MAX_OPEN_CONNS" envDefault:"66"`
		EnableLog    bool   `env:"IGN_USER_DB_LOG" envDefault:"false"`
		SysAdmin     string `env:"IGN_SYSADMIN" envDefault:""`
	}

	userCfg := UserConfig{}
	// Also using env-to-struct approach to read configuration
	if err := env.Parse(&userCfg); err != nil {
		return errors.Wrap(err, "Error parsing environment into userDB UserConfig struct. %+v\n")
	}

	// Create the database config struct
	ignDbCfg := ign.DatabaseConfig{
		UserName:     userCfg.UserName,
		Password:     userCfg.Password,
		Address:      userCfg.Address,
		Name:         userCfg.Name,
		MaxOpenConns: userCfg.MaxOpenConns,
		EnableLog:    userCfg.EnableLog,
	}

	// Connect to the database.
	usersDb, err := ign.InitDbWithCfg(&ignDbCfg)
	if err != nil {
		return err
	}
	// Tell the server about the user database
	app.server.UsersDb = usersDb

	// Initialize permissions. This requires the `permissions/policy.conf` file.
	logger.Debug("Initializing user permissions")
	perm := &permissions.Permissions{}

	// \todo Error?: I have to pass in `userCfg.SysAdmin` here in order for
	// system administators to be loaded into the `casbin_rule` table in the
	// cloudsim_nps DB. Passing userCfg.SysAdmin to  useracc.NewService()
	// doesn't seem to do anything.
	err = perm.Init(app.db, userCfg.SysAdmin)
	if err != nil {
		return err
	}

	logCtx := ign.NewContextWithLogger(context.Background(), logger)
	userAccessorService, err := useracc.NewService(logCtx,
		perm, usersDb, userCfg.SysAdmin)
	if err != nil {
		return err
	}

	HTTPHandlerInstance, err = NewHTTPHandler(logCtx, userAccessorService)
	if err != nil {
		logger.Critical("Critical error trying to create the HTTPHandler", err)
		return err
	}

	return nil
}
