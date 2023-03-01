package actions

import (
	gormUtils "github.com/gazebo-web/gz-go/v7/database/gorm"
	"github.com/jinzhu/gorm"
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
