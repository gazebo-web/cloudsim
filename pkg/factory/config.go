package factory

// Config is a factory configuration.
// It is passed to a factory to request a component type and provide it with configuration information.
type Config struct {
	// Type returns the object type to create using a factory.
	Type string `yaml:"type"`
	// Config contains configuration data required by the factory to create the request object.
	// Note that the factory creating the object has no way to know about the implementation details of the object
	// being created. As such, the NewFunc creating the object will need to marshal the data contained in the config.
	Config map[string]interface{} `yaml:"config"`
}
