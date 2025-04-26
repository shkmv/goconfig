package goconfig

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigFromYAML(t *testing.T) {
	testFile := "testdata/config.yaml"

	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Fatalf("Test file %s does not exist", testFile)
	}

	cfg, err := Load[TestConfig](WithFile(testFile))
	if err != nil {
		t.Fatalf("Failed to load config from YAML: %v", err)
	}

	if cfg.DB.Host != "localhost" {
		t.Errorf("Expected DB.Host to be 'localhost', got '%s'", cfg.DB.Host)
	}

	if cfg.DB.Port != 5432 {
		t.Errorf("Expected DB.Port to be 5432, got %d", cfg.DB.Port)
	}

	if cfg.Port != 3000 {
		t.Errorf("Expected Port to be 3000, got %d", cfg.Port)
	}
}

func TestLoadComplexConfigFromYAML(t *testing.T) {
	testFile := "testdata/complex_config.yaml"

	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Fatalf("Test file %s does not exist", testFile)
	}

	type ComplexConfig struct {
		String string  `config:"string"`
		Int    int     `config:"int"`
		Float  float64 `config:"float"`
		Bool   bool    `config:"bool"`
		Nested struct {
			Value1 string `config:"value1"`
			Value2 int    `config:"value2"`
		} `config:"nested"`
		BoolValues struct {
			True1  bool `config:"true1"`
			True2  bool `config:"true2"`
			True3  bool `config:"true3"`
			False1 bool `config:"false1"`
			False2 bool `config:"false2"`
			False3 bool `config:"false3"`
		} `config:"bool_values"`
	}

	cfg, err := Load[ComplexConfig](WithFile(testFile))
	if err != nil {
		t.Fatalf("Failed to load complex config from YAML: %v", err)
	}

	if cfg.String != "test-string" {
		t.Errorf("Expected String to be 'test-string', got '%s'", cfg.String)
	}

	if cfg.Int != 42 {
		t.Errorf("Expected Int to be 42, got %d", cfg.Int)
	}

	if cfg.Float != 3.14159 {
		t.Errorf("Expected Float to be 3.14159, got %f", cfg.Float)
	}

	if cfg.Nested.Value1 != "nested-value" {
		t.Errorf("Expected Nested.Value1 to be 'nested-value', got '%s'", cfg.Nested.Value1)
	}

	if cfg.Nested.Value2 != 123 {
		t.Errorf("Expected Nested.Value2 to be 123, got %d", cfg.Nested.Value2)
	}

	if !cfg.BoolValues.True1 {
		t.Errorf("Expected Bool.True1 to be true, got %v", cfg.BoolValues.True1)
	}
	if !cfg.BoolValues.True2 {
		t.Errorf("Expected Bool.True2 to be true, got %v", cfg.BoolValues.True2)
	}
	if !cfg.BoolValues.True3 {
		t.Errorf("Expected Bool.True3 to be true, got %v", cfg.BoolValues.True3)
	}
	if cfg.BoolValues.False1 {
		t.Errorf("Expected Bool.False1 to be false, got %v", cfg.BoolValues.False1)
	}
	if cfg.BoolValues.False2 {
		t.Errorf("Expected Bool.False2 to be false, got %v", cfg.BoolValues.False2)
	}
	if cfg.BoolValues.False3 {
		t.Errorf("Expected Bool.False3 to be false, got %v", cfg.BoolValues.False3)
	}
}

func TestLoadConfigFromNonExistentYAML(t *testing.T) {
	_, err := Load[TestConfig](WithFile("/non/existent/file.yaml"))

	if err == nil {
		t.Error("Expected error when loading non-existent YAML file, but got nil")
	}
}

func TestLoadConfigFromInvalidYAML(t *testing.T) {
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

	_, err = Load[TestConfig](WithFile(tempFile))

	// Verify that an error is returned
	if err == nil {
		t.Error("Expected error when loading invalid YAML, but got nil")
	}
}

func TestLoadConfigFromMultipleSources(t *testing.T) {
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "partial.yaml")

	partialYAML := `
db:
  host: yaml-host
`
	err := os.WriteFile(tempFile, []byte(partialYAML), 0644)
	if err != nil {
		t.Fatalf("Failed to create test YAML file: %v", err)
	}

	os.Setenv("APP_DB_PORT", "5432")
	os.Setenv("APP_PORT", "3000")

	cfg, err := Load[TestConfig](
		WithFile(tempFile),
		WithEnv("APP_"),
	)
	if err != nil {
		t.Fatalf("Failed to load config from multiple sources: %v", err)
	}

	if cfg.DB.Host != "yaml-host" {
		t.Errorf("Expected DB.Host to be 'yaml-host', got '%s'", cfg.DB.Host)
	}

	if cfg.DB.Port != 5432 {
		t.Errorf("Expected DB.Port to be 5432, got %d", cfg.DB.Port)
	}

	if cfg.Port != 3000 {
		t.Errorf("Expected Port to be 3000, got %d", cfg.Port)
	}

	os.Unsetenv("APP_DB_PORT")
	os.Unsetenv("APP_PORT")
}

func TestLoadConfigWithOverridingValues(t *testing.T) {
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "base.yaml")

	baseYAML := `
db:
  host: yaml-host
  port: 5432
port: 3000
`
	err := os.WriteFile(tempFile, []byte(baseYAML), 0644)
	if err != nil {
		t.Fatalf("Failed to create test YAML file: %v", err)
	}

	os.Setenv("APP_DB_HOST", "env-host")
	os.Setenv("APP_PORT", "4000")

	cfg, err := Load[TestConfig](
		WithFile(tempFile),
		WithEnv("APP_"),
	)
	if err != nil {
		t.Fatalf("Failed to load config with overriding values: %v", err)
	}

	if cfg.DB.Host != "env-host" {
		t.Errorf("Expected DB.Host to be 'env-host', got '%s'", cfg.DB.Host)
	}

	if cfg.DB.Port != 5432 {
		t.Errorf("Expected DB.Port to be 5432, got %d", cfg.DB.Port)
	}

	if cfg.Port != 4000 {
		t.Errorf("Expected Port to be 4000, got %d", cfg.Port)
	}

	os.Unsetenv("APP_DB_HOST")
	os.Unsetenv("APP_PORT")
}
