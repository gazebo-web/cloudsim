package gorm

import (
	"errors"
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

// GetDBFromEnvVars reads environment variables to return a Gorm database connection.
// The environment variables used are:
// * IGN_DB_ADDRESS Address of the DBMS.
// * IGN_DB_USERNAME Username to connect to the DBMS with.
// * IGN_DB_PASSWORD Password to connect to the DBMS with.
// * IGN_DB_NAME Name of the database to connect to.
// * IGN_DB_MAX_OPEN_CONNS - (Optional) You run the risk of getting a 'too many connections' error if this is not set.
func GetDBFromEnvVars() (*gorm.DB, error) {
	// Get the db config
	var dbConfig ign.DatabaseConfig
	var err error
	dbConfig, err = ign.NewDatabaseConfigFromEnvVars()
	if err != nil {
		return nil, err
	}

	// Connect to the db
	db, err := ign.InitDbWithCfg(&dbConfig)
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

	tx.AutoMigrate(models...)

	return nil
}

// CleanAndMigrateModels drops existing target model tables and recreates them.
func CleanAndMigrateModels(tx *gorm.DB, models ...interface{}) error {
	if tx == nil {
		return errors.New("attempted to clean database with an invalid tx")
	}

	tx.DropTableIfExists(models...)

	return MigrateModels(tx, models...)
}
