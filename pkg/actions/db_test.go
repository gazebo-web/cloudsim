package actions

import (
	"github.com/jinzhu/gorm"
	gormUtils "gitlab.com/ignitionrobotics/web/cloudsim/pkg/utils/db/gorm"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"golang.org/x/net/context"
	"testing"
)

// TestResource contains resources used for testing.
type TestResource struct {
	store  Store
	logger *ign.Logger
	db     *gorm.DB
}

type storeTestData struct {
	value int
}

// setupTest can be called at the start of a test in the actions package to get a set of common values used for testing.
func setupTest(t *testing.T) *TestResource {
	ctx := context.Background()
	logger := ign.LoggerFromContext(ctx)
	db, err := gormUtils.GetDBFromEnvVars()
	defer db.Close()

	if err != nil {
		t.Fatalf("Could not connect to database: %s", err)
	}

	// Migrate the action models
	err = MigrateDB(db)
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
