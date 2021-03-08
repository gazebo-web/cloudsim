package actions

import (
	"github.com/jinzhu/gorm"
	gormUtils "gitlab.com/ignitionrobotics/web/cloudsim/pkg/utils/db/gorm"
)

// CleanAndMigrateDB cleans and migrates action database models, indexes and keys.
func CleanAndMigrateDB(tx *gorm.DB) error {
	return gormUtils.CleanAndMigrateModels(
		tx,
		&Deployment{},
		&deploymentData{},
		&DeploymentError{},
	)
}

// MigrateDB migrates action database models, indexes and keys.
func MigrateDB(tx *gorm.DB) error {
	return gormUtils.MigrateModels(
		tx,
		&Deployment{},
		&deploymentData{},
		&DeploymentError{},
	)
}

// DropDB drops action database models, indexes and keys.
func DropDB(tx *gorm.DB) error {
	return gormUtils.DropModels(
		tx,
		&DeploymentError{},
		&deploymentData{},
		&Deployment{},
	)
}
