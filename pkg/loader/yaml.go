package loader

import (
	"github.com/pkg/errors"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/defaults"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// yamlLoader is a Loader implementation to parse YAML files.
type yamlLoader struct {
	logger ign.Logger
}

// log logs a message using the logger contained by the loader. If logger is `nil`, no logging is performed.
func (l *yamlLoader) log(interfaces ...interface{}) {
	if l.logger != nil {
		l.logger.Warning(interfaces...)
	}
}

// parseBytes parses a YAML file passed in as bytes.
// `out` must be a pointer to the resulting configuration object.
func (l *yamlLoader) parseBytes(data []byte, out interface{}) error {
	// Parse configuration file values
	err := yaml.Unmarshal(data, out)
	if err != nil {
		l.log("Failed to parse YAML file", err)
		return err
	}

	return nil
}

// Load reads and parses a YAML file from a path.
// `out` must be a pointer to the object the file will be parsed into.
func (l *yamlLoader) Load(path string, out interface{}) error {
	// Read the configuration file
	configFile, err := ioutil.ReadFile(path)
	if err != nil {
		l.log("Failed to read YAML file", err)
		return errors.Wrap(ErrLoadFailed, err.Error())
	}

	// Parse configuration
	if err := l.parseBytes(configFile, out); err != nil {
		return err
	}

	// Apply default values if available
	if err := defaults.SetDefaults(out); err != nil {
		return err
	}

	return nil
}

// NewYAMLLoader returns a YAML loader.
func NewYAMLLoader(logger ign.Logger) Loader {
	return &yamlLoader{
		logger: logger,
	}
}
