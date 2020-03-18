package logger

import (
	"github.com/caarlos0/env"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"strconv"
)

type config struct {
	LogVerbosity        string `env:"IGN_LOGGER_VERBOSITY"`
	RollbarLogVerbosity string `env:"IGN_LOGGER_ROLLBAR_VERBOSITY"`
}

func New() (ign.Logger, error) {
	cfg := config{}
	var err error

	if err = env.Parse(&cfg); err != nil {
		return nil, err
	}

	verbosity := ign.VerbosityWarning
	if cfg.LogVerbosity != "" {
		verbosity, err = strconv.Atoi(cfg.LogVerbosity)
		if err != nil {
			return nil, err
		}
	}

	rollbarVerbosity := ign.VerbosityWarning
	if cfg.RollbarLogVerbosity != "" {
		rollbarVerbosity, err = strconv.Atoi(cfg.RollbarLogVerbosity)
		if err != nil {
			return nil, err
		}
	}

	std := ign.ReadStdLogEnvVar()
	logger := ign.NewLoggerWithRollbarVerbosity("init", std, verbosity, rollbarVerbosity)
	return logger, nil
}