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

// FromDotEnv loads configuration from a .env file.
func (c *Config) FromDotEnv(path string) *Config {
    c.sources = append(c.sources, sources.NewDotEnvSource(path))
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

	if err := internal.Bind(merged, target); err != nil {
		return fmt.Errorf("binding configuration to target: %w", err)
	}
	return nil
}
