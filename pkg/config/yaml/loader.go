package yaml

import (
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/config"
	"gitlab.com/ignitionrobotics/web/ign-go"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type loader struct {
	logger ign.Logger
}

// log logs a message using the logger contained by the loader. If logger is `nil`, no logging is performed.
func (l *loader) log(interfaces ...interface{}) {
	if l.logger != nil {
		l.logger.Warning(interfaces...)
	}
}

// LoadConfiguration parses a YAML file passed in as bytes.
// `out` must be a pointer to the resulting configuration object.
func (l *loader) parseBytes(cfg []byte, out interface{}) error {
	// Parse configuration file values
	err := yaml.Unmarshal(cfg, out)
	if err != nil {
		l.log("Failed to parse YAML file", err)
		return err
	}

	return nil
}

// Load reads and parses a YAML file from a path.
// `out` must be a pointer to the resulting configuration object.
func (l *loader) Load(cfgPath string, out interface{}) error {
	// Read the configuration file
	configFile, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		l.log("Failed to read YAML file", err)
		return err
	}

	// Parse configuration
	return l.parseBytes(configFile, out)
}

// NewLoader returns a YAML loader.
func NewLoader(logger ign.Logger) config.Loader {
	return &loader{
		logger: logger,
	}
}
