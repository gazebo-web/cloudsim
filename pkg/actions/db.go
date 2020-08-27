package actions

import (
	"errors"
	"github.com/jinzhu/gorm"
)

// migrateModels migrates database models.
func migrateModels(tx *gorm.DB, clean bool, models ...interface{}) error {
	if tx != nil {
		// Optionally clean db
		if clean {
			tx.DropTableIfExists(
				models...,
			)
		}
		// Create models
		tx.AutoMigrate(
			models...,
		)
	} else {
		return errors.New("attempted to migrate with an invalid tx")
	}

	return nil
}

// MigrateDB migrates application models, indexes and keys to the database.
func MigrateDB(tx *gorm.DB, clean bool) error {
	return migrateModels(
		tx,
		clean,
		&Deployment{},
		&deploymentData{},
		&DeploymentError{},
	)
}
