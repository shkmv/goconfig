package goconfig

// Load generic configuration into a struct.
// example:
//
//	cfg := Load[Config](WithEnv("APP"))
func Load[T any](opts ...Option) (T, error) {
	var zero T

	cfg := New()
	for _, opt := range opts {
		opt(cfg)
	}

	var target T
	if err := cfg.Bind(&target); err != nil {
		return zero, err
	}

	return target, nil
}
