package sources

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileSource_Load(t *testing.T) {
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "config.yaml")

	yamlContent := `
db:
  host: localhost
  port: 5432
port: 3000
`
	err := os.WriteFile(tempFile, []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test YAML file: %v", err)
	}

	fileSource := NewFileSource(tempFile)
	data, err := fileSource.Load()
	if err != nil {
		t.Fatalf("Failed to load configuration from file: %v", err)
	}

	if data == nil {
		t.Fatal("Expected non-nil data, got nil")
	}

	port, ok := data["port"]
	if !ok {
		t.Error("Expected 'port' key in configuration, but it was not found")
	} else if port != 3000 {
		t.Errorf("Expected port to be 3000, got %v", port)
	}

	db, ok := data["db"].(map[string]any)
	if !ok {
		t.Error("Expected 'db' key to be a map, but it was not found or has wrong type")
	} else {
		host, ok := db["host"]
		if !ok {
			t.Error("Expected 'db.host' key in configuration, but it was not found")
		} else if host != "localhost" {
			t.Errorf("Expected db.host to be 'localhost', got %v", host)
		}

		dbPort, ok := db["port"]
		if !ok {
			t.Error("Expected 'db.port' key in configuration, but it was not found")
		} else if dbPort != 5432 {
			t.Errorf("Expected db.port to be 5432, got %v", dbPort)
		}
	}
}

func TestFileSource_LoadNonExistentFile(t *testing.T) {
	fileSource := NewFileSource("/non/existent/file.yaml")
	_, err := fileSource.Load()

	if err == nil {
		t.Error("Expected error when loading non-existent file, but got nil")
	}
}

func TestFileSource_LoadInvalidYAML(t *testing.T) {
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "invalid.yaml")

	invalidYAML := `
this is not valid YAML
  - missing colon
    indentation is wrong
`
	err := os.WriteFile(tempFile, []byte(invalidYAML), 0644)
	if err != nil {
		t.Fatalf("Failed to create test YAML file: %v", err)
	}

	fileSource := NewFileSource(tempFile)
	_, err = fileSource.Load()

	if err == nil {
		t.Error("Expected error when loading invalid YAML, but got nil")
	}
}
