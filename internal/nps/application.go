package nps

// This file contains structures and functions that manage the application.
// This is the primary entry point for the nps application.

import (
	"github.com/jinzhu/gorm"
	ignGorm "gitlab.com/ignitionrobotics/web/cloudsim/pkg/utils/db/gorm"
	"gitlab.com/ignitionrobotics/web/fuelserver/permissions"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"gitlab.com/ignitionrobotics/web/ign-go/monitoring/prometheus"
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

	logger.Debug("Creating the application")

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

	// Initialize permissions. This requires the `permissions/policy.conf` file.
	logger.Debug("Initializing user permissions")
	perm := &permissions.Permissions{}
	err = perm.Init(db, "sysadmin")
	if err != nil {
		return nil, err
	}

	app := &application{
		controller: NewController(db, logger),
		db:         db,
	}

	// Create a server monitoring provider. Specifying a provider makes the
	// server automatically add middleware required to track metrics.
	monitoring := prometheus.NewPrometheusProvider("")

	app.server, err = ign.Init("", "", monitoring)

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
