package config

// Loader is used to load configs
type Loader interface {
	Load(cfgPath string, out interface{}) error
}
