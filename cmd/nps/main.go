package main

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/internal/nps/application"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"os"
)

// Configure your application.
//
// 1. Setup environment variables in a .env file
//     * IGN_DB_ADDRESS Address of the DBMS.
//     * IGN_DB_USERNAME Username to connect to the DBMS with.
//     * IGN_DB_PASSWORD Password to connect to the DBMS with.
//     * IGN_DB_NAME Name of the database to connect to.
//     * IGN_DB_MAX_OPEN_CONNS - (Optional) You run the risk of getting a
//                           'too many connections' error if this is not set.
func main() {
	// Create a new logger. This will be used to log messages.
	// The logger must be setup first.
	logger := ign.NewLoggerNoRollbar("NPS", ign.VerbosityDebug)

	// Create the application
	application, err := nps.NewApplication("1.0", logger)
	if err != nil {
		logger.Error("main: error:", err)
		os.Exit(1)
	}

	// Run the application
	application.Run()
}
