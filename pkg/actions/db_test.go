package actions

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/require"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"testing"
)

// TestResource contains resources used for testing.
type TestResource struct {
	ctx    *context.Context
	logger *ign.Logger
	db     *gorm.DB
}

func setupTest(t *testing.T) *TestResource {
	ctx := context.Background()
	logger := ign.LoggerFromContext(ctx)

	// Get the db config
	var dbConfig ign.DatabaseConfig
	var err error
	dbConfig, err = ign.NewDatabaseConfigFromEnvVars()
	require.NoError(t, err, "Could not read database config from env vars")

	// Connect to the db
	db, err := ign.InitDbWithCfg(&dbConfig)
	require.NoError(t, err, "Could not connect to the db.")

	// Migrate database tables
	err = MigrateDB(db, true)
	require.NoError(t, err, "Could not migrate actions to the db.")

	// Create the test resource container
	testResources := TestResource{
		ctx:    &ctx,
		logger: &logger,
		db:     db,
	}

	return &testResources
}
