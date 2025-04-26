package sources

import (
	"os"

	"gopkg.in/yaml.v3"
)

// FileSource represents a source that loads configuration from a file.
type FileSource struct {
	path string
}

// NewFileSource creates a new FileSource instance.
func NewFileSource(path string) *FileSource {
	return &FileSource{path: path}
}

// Load loads the configuration from the file.
func (f *FileSource) Load() (map[string]any, error) {
	data, err := os.ReadFile(f.path)
	if err != nil {
		return nil, err
	}

	var out map[string]any
	if err := yaml.Unmarshal(data, &out); err != nil {
		return nil, err
	}

	// TODO: validate
	return out, nil
}
