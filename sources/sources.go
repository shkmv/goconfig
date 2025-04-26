package sources

type Source interface {
	Load() (map[string]any, error)
}
