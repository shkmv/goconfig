package sources

import (
	"os"
	"strings"
)

// EnvSource represents a configuration source that loads from environment variables.
type EnvSource struct {
	prefix string
}

// NewEnvSource creates a new EnvSource with the specified prefix.
// Environment variables starting with this prefix will be loaded into the configuration.
func NewEnvSource(prefix string) *EnvSource {
	return &EnvSource{
		prefix: prefix,
	}
}

// Load loads configuration values from environment variables.
func (e *EnvSource) Load() (map[string]any, error) {
	out := make(map[string]any)
	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		key := parts[0]
		if !strings.HasPrefix(key, e.prefix) {
			continue
		}
		value := parts[1]
		cleanKey := strings.TrimPrefix(key, e.prefix)
		dotKey := strings.ToLower(strings.ReplaceAll(cleanKey, "_", "."))

		setNestedValue(out, strings.Split(dotKey, "."), value)
	}
	return out, nil
}

func setNestedValue(data map[string]any, keys []string, value any) {
	if len(keys) == 0 {
		return
	}

	key := keys[0]
	if len(keys) == 1 {
		data[key] = value
		return
	}

	if _, ok := data[key]; !ok {
		data[key] = make(map[string]any)
	} else if _, ok := data[key].(map[string]any); !ok {
		data[key] = make(map[string]any)
	}

	setNestedValue(data[key].(map[string]any), keys[1:], value)
}
