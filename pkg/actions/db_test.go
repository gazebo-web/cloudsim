package actions

import (
	"context"
	"github.com/gazebo-web/gz-go/v7"
	gormUtils "github.com/gazebo-web/gz-go/v7/database/gorm"
	"github.com/jinzhu/gorm"
	"testing"
)

// TestResource contains resources used for testing.
type TestResource struct {
	store  Store
	logger *gz.Logger
	db     *gorm.DB
}

type storeTestData struct {
}

// setupTest can be called at the start of a test in the actions package to get a set of common values used for testing.
func setupTest(t *testing.T) *TestResource {
	ctx := context.Background()
	logger := gz.LoggerFromContext(ctx)
	db, err := gormUtils.GetTestDBFromEnvVars()

	if err != nil {
		t.Fatalf("Could not connect to database: %s", err)
	}

	// Migrate the action models
	err = CleanAndMigrateDB(db)
	if err != nil {
		t.Fatalf("Could not migrate actions database models: %s", err)
	}

	// Create the test resource container
	testResources := TestResource{
		store:  NewStore(&storeTestData{}),
		logger: &logger,
		db:     db,
	}

	return &testResources
}
