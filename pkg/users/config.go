package users

import (
	"github.com/caarlos0/env"
	"github.com/pkg/errors"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

// dbConfig
type dbConfig struct {
	UserName     string `env:"IGN_USER_DB_USERNAME" envDefault:":notset"`
	Password     string `env:"IGN_USER_DB_PASSWORD"`
	Address      string `env:"IGN_USER_DB_ADDRESS"`
	Name         string `env:"IGN_USER_DB_NAME" envDefault:"usersdb"`
	MaxOpenConns int    `env:"IGN_USER_DB_MAX_OPEN_CONNS" envDefault:"66"`
	EnableLog    bool   `env:"IGN_USER_DB_LOG" envDefault:"false"`
}

// newDbConfig
func newDbConfig() (*ign.DatabaseConfig, error) {
	cfg := dbConfig{}
	if err := env.Parse(&cfg); err != nil {
		return nil, errors.Wrap(err, "Error parsing environment into userDB dbConfig struct. %+v\n")
	}

	dbCfg := ign.DatabaseConfig{
		UserName:     cfg.UserName,
		Password:     cfg.Password,
		Address:      cfg.Address,
		Name:         cfg.Name,
		MaxOpenConns: cfg.MaxOpenConns,
		EnableLog:    cfg.EnableLog,
	}

	return &dbCfg, nil
}
