package sources

// Source is an interface for configuration sources.
// Implementations of this interface can load configuration data from different sources
// such as environment variables, files, etc.
type Source interface {
	// Load loads and returns configuration data as a map.
	// It returns an error if the loading process fails.
	Load() (map[string]any, error)
}
