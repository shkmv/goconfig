package goconfig

import (
	"fmt"
	"sync"

	"github.com/shkmv/goconfig/internal"
	"github.com/shkmv/goconfig/sources"
)

// Config represents a configuration object.
type Config struct {
	sources []sources.Source
	mu      sync.RWMutex
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

// FromEtcd loads configuration from etcd.
func (c *Config) FromEtcd(endpoints []string, prefix string, user string, password string) *Config {
	etcdSource, err := sources.NewEtcdSource(endpoints, prefix, user, password)
	if err != nil {
		// TODO: handle error
		return c
	}
	c.sources = append(c.sources, etcdSource)
	return c
}

// Bind binds the configuration to a target struct.
func (c *Config) Bind(target any) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	merged := make(map[string]any)
	for _, src := range c.sources {
		data, err := src.Load()
		if err != nil {
			return fmt.Errorf("loading config from %T: %w", src, err)
		}
		merged = internal.Merge(merged, data)

		if watchedSrc, ok := src.(sources.WatchedSource); ok {
			err := watchedSrc.Watch(func(data map[string]any) {
				c.mu.Lock()
				defer c.mu.Unlock()
				newMerged := internal.Merge(merged, data)
				// TODO: log error
				_ = internal.Bind(newMerged, target)
			})
			if err != nil {
				return fmt.Errorf("watching config from %T: %w", src, err)
			}
		}
	}

	if err := internal.Bind(merged, target); err != nil {
		return fmt.Errorf("binding configuration to target: %w", err)
	}

	return nil
}
