package db

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

func NewConfig() ign.DatabaseConfig {
	config, err := ign.NewDatabaseConfigFromEnvVars()
	if err != nil {
		panic(err)
	}
	return config
}

func NewTestConfig() ign.DatabaseConfig {
	config := NewConfig()
	config.Name = config.Name + "_test"
	return config
}

func Must(db *gorm.DB, err error) *gorm.DB {
	if err != nil {
		panic(err)
	}
	return db
}

func NewDB(config ign.DatabaseConfig) (*gorm.DB, error) {
	db, err := ign.InitDbWithCfg(&config)
	if err != nil {
		return nil, err
	}
	return db, nil
}
