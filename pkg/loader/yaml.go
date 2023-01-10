package loader

import (
	"github.com/gazebo-web/cloudsim/v4/pkg/defaults"
	"github.com/gazebo-web/gz-go/v7"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
)

const extensionYAML = ".yaml"
const extensionYML = ".yml"

// yamlLoader is a Loader implementation to parse YAML files.
type yamlLoader struct {
	logger gz.Logger
}

// Filter returns the list of files that contain either .yaml or .yml.
func (l *yamlLoader) Filter(list []string) []string {
	result := make([]string, 0, len(list))
	for _, item := range list {
		if strings.Contains(item, extensionYAML) || strings.Contains(item, extensionYML) {
			result = append(result, item)
		}
	}
	return result
}

// TrimExt removes the .yaml and .yml extension from the given filename.
func (l *yamlLoader) TrimExt(filename string) string {
	return trimExts(filename, extensionYAML, extensionYML)
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
	configFile, err := os.ReadFile(path)
	if err != nil {
		l.log("Failed to read YAML file", err)
		return errors.Wrap(ErrLoadFailed, err.Error())
	}

	// Parse configuration
	if err := l.parseBytes(configFile, out); err != nil {
		return err
	}

	// Apply default values if available
	if err := defaults.SetValues(out); err != nil {
		return err
	}

	return nil
}

// NewYAMLLoader returns a YAML loader.
func NewYAMLLoader(logger gz.Logger) Loader {
	return &yamlLoader{
		logger: logger,
	}
}
