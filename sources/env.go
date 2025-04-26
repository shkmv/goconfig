package sources

import (
	"os"
	"strings"
)

type EnvSource struct {
	prefix string
}

func NewEnvSource(prefix string) *EnvSource {
	return &EnvSource{
		prefix: prefix,
	}
}

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
		out[strings.ToLower(strings.ReplaceAll(cleanKey, "_", "."))] = value
	}
	return out, nil
}
