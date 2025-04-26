package goconfig

import "github.com/shkmv/goconfig/sources"

type Option func(*Config)

func WithEnv(prefix string) Option {
	return func(c *Config) {
		c.sources = append(c.sources, sources.NewEnvSource(prefix))
	}
}
