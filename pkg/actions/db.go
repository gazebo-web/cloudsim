package actions

import (
	"github.com/jinzhu/gorm"
	gormUtils "gitlab.com/ignitionrobotics/web/cloudsim/pkg/utils/db/gorm"
)

// MigrateDB migrates application models, indexes and keys to the database.
func migrateDB(tx *gorm.DB) error {
	return gormUtils.CleanAndMigrateModels(
		tx,
		&Deployment{},
		&deploymentData{},
		&DeploymentError{},
	)
}
