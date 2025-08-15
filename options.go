package goconfig

import "github.com/shkmv/goconfig/sources"

// Option represents a configuration option.
type Option func(*Config)

// WithEnv adds an environment variable source to the configuration.
func WithEnv(prefix string) Option {
	return func(c *Config) {
		c.sources = append(c.sources, sources.NewEnvSource(prefix))
	}
}

// WithFile adds a file source to the configuration.
func WithFile(path string) Option {
    return func(c *Config) {
        c.sources = append(c.sources, sources.NewFileSource(path))
    }
}

// WithDotEnv adds a .env file source to the configuration.
func WithDotEnv(path string) Option {
    return func(c *Config) {
        c.sources = append(c.sources, sources.NewDotEnvSource(path))
    }
}
