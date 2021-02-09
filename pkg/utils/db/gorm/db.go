package gorm

import (
	"errors"
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"log"
)

var (
	// ErrEmptyDatabaseName is returned when a database configuration contains an empty database name
	ErrEmptyDatabaseName = errors.New("db config contains empty database name")
)

// getDBConfigFromEnvVars reads environment variables to return a database connection configuration.
// The environment variables used are:
// * IGN_DB_ADDRESS Address of the DBMS.
// * IGN_DB_USERNAME Username to connect to the DBMS with.
// * IGN_DB_PASSWORD Password to connect to the DBMS with.
// * IGN_DB_NAME Name of the database to connect to.
// * IGN_DB_MAX_OPEN_CONNS - (Optional) You run the risk of getting a 'too many connections' error if this is not set.
func getDBConfigFromEnvVars() (*ign.DatabaseConfig, error) {
	// Get the db config
	var dbConfig ign.DatabaseConfig
	var err error
	dbConfig, err = ign.NewDatabaseConfigFromEnvVars()
	if err != nil {
		return nil, err
	}
	if len(dbConfig.Name) == 0 {
		return nil, ErrEmptyDatabaseName
	}

	return &dbConfig, nil
}

// GetDBFromEnvVars reads environment variables to return a Gorm database connection.
func GetDBFromEnvVars() (*gorm.DB, error) {
	// Get the db config
	dbConfig, err := getDBConfigFromEnvVars()
	if err != nil {
		return nil, err
	}

	// Connect to the db
	db, err := ign.InitDbWithCfg(dbConfig)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// GetTestDBFromEnvVars reads environment variables to return a Gorm database connection.
func GetTestDBFromEnvVars() (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		return nil, err
	}
	return db, nil
}

// MigrateModels migrates database models.
// If the model table already exists, it will be updated to reflect the model structure. The table will only have
// columns added or updated but not dropped.
func MigrateModels(tx *gorm.DB, models ...interface{}) error {
	if tx == nil {
		return errors.New("attempted to migrate with an invalid tx")
	}

	if err := tx.AutoMigrate(models...).Error; err != nil {
		log.Println("MigrateModels: Error while running AutoMigrate, error:", err)
		return err
	}

	return nil
}

// CleanAndMigrateModels drops existing target model tables and recreates them.
func CleanAndMigrateModels(tx *gorm.DB, models ...interface{}) error {
	if tx == nil {
		return errors.New("attempted to clean database with an invalid tx")
	}

	log.Printf("CleanAndMigrateModels: Dropping tables: %v\n", models)
	if err := tx.DropTableIfExists(models...).Error; err != nil {
		log.Println("CleanAndMigrateModels: Error while running DropTableIfExists, error:", err)
		return err
	}

	log.Printf("CleanAndMigrateModels: Migrating tables: %v\n", models)
	if err := MigrateModels(tx, models...); err != nil {
		log.Println("CleanAndMigrateModels: Error while running MigrateModels, error:", err)
		return err
	}

	return nil
}
