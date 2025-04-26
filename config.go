package goconfig

import (
	"fmt"

	"github.com/shkmv/goconfig/internal"
	"github.com/shkmv/goconfig/sources"
)

// Config represents a configuration object.
type Config struct {
	sources []sources.Source
}

func New() *Config {
	return &Config{}
}

// FromEnv loads configuration from environment variables.
func (c *Config) FromEnv(prefix string) *Config {
	c.sources = append(c.sources, sources.NewEnvSource(prefix))
	return c
}

// FromFile loads configuration from a file.
func (c *Config) FromFile(path string) *Config {
	c.sources = append(c.sources, sources.NewFileSource(path))
	return c
}

// Bind binds the configuration to a target struct.
func (c *Config) Bind(target any) error {
	merged := make(map[string]any)
	for _, src := range c.sources {
		data, err := src.Load()
		if err != nil {
			return fmt.Errorf("loading config from %T: %w", src, err)
		}
		merged = internal.Merge(merged, data)
	}

	return internal.Bind(merged, target)
}
