package goconfig

import (
    "os"
    "path/filepath"
    "testing"
    "strings"

    "github.com/shkmv/goconfig/sources"
)

type TestConfig struct {
	DB struct {
		Host string `config:"host"`
		Port int    `config:"port"`
	} `config:"db"`
	Port int `config:"port"`
}

func TestLoadConfigFromEnv(t *testing.T) {
	os.Setenv("APP_DB_HOST", "localhost")
	os.Setenv("APP_DB_PORT", "5432")
	os.Setenv("APP_PORT", "3000")

	cfg, err := Load[TestConfig](WithEnv("APP_"))
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
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

func TestLoadConfigWithInvalidTypes(t *testing.T) {
	os.Setenv("APP_PORT", "not_a_number")

	_, err := Load[TestConfig](WithEnv("APP_"))

	if err == nil {
		t.Error("Expected error when loading invalid type, but got nil")
	}
}

func TestLoadConfigWithEmptyValues(t *testing.T) {
	os.Unsetenv("APP_DB_HOST")
	os.Unsetenv("APP_DB_PORT")
	os.Unsetenv("APP_PORT")

	cfg, err := Load[TestConfig](WithEnv("APP_"))
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg.DB.Host != "" {
		t.Errorf("Expected DB.Host to be empty, got '%s'", cfg.DB.Host)
	}

	if cfg.DB.Port != 0 {
		t.Errorf("Expected DB.Port to be 0, got %d", cfg.DB.Port)
	}

	if cfg.Port != 0 {
		t.Errorf("Expected Port to be 0, got %d", cfg.Port)
	}
}

func TestLoadConfigWithDifferentTypes(t *testing.T) {
	type ComplexConfig struct {
		String string  `config:"string"`
		Int    int     `config:"int"`
		Float  float64 `config:"float"`
		Bool   bool    `config:"bool"`
		Nested struct {
			Value1 string `config:"value1"`
			Value2 int    `config:"value2"`
		} `config:"nested"`
	}

	os.Setenv("APP_STRING", "test-string")
	os.Setenv("APP_INT", "42")
	os.Setenv("APP_FLOAT", "3.14159")
	os.Setenv("APP_BOOL", "true")
	os.Setenv("APP_NESTED_VALUE1", "nested-value")
	os.Setenv("APP_NESTED_VALUE2", "123")

	cfg, err := Load[ComplexConfig](WithEnv("APP_"))
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
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

	if !cfg.Bool {
		t.Errorf("Expected Bool to be true, got %v", cfg.Bool)
	}

	if cfg.Nested.Value1 != "nested-value" {
		t.Errorf("Expected Nested.Value1 to be 'nested-value', got '%s'", cfg.Nested.Value1)
	}

	if cfg.Nested.Value2 != 123 {
		t.Errorf("Expected Nested.Value2 to be 123, got %d", cfg.Nested.Value2)
	}

	os.Unsetenv("APP_STRING")
	os.Unsetenv("APP_INT")
	os.Unsetenv("APP_FLOAT")
	os.Unsetenv("APP_BOOL")
	os.Unsetenv("APP_NESTED_VALUE1")
	os.Unsetenv("APP_NESTED_VALUE2")
}

func TestLoadConfigWithBooleanValues(t *testing.T) {
	type BoolConfig struct {
		Bool struct {
			True1  bool `config:"true1"`
			True2  bool `config:"true2"`
			True3  bool `config:"true3"`
			False1 bool `config:"false1"`
			False2 bool `config:"false2"`
			False3 bool `config:"false3"`
		} `config:"bool"`
	}

	os.Setenv("APP_BOOL_TRUE1", "true")
	os.Setenv("APP_BOOL_TRUE2", "1")
	os.Setenv("APP_BOOL_TRUE3", "yes")
	os.Setenv("APP_BOOL_FALSE1", "false")
	os.Setenv("APP_BOOL_FALSE2", "0")
	os.Setenv("APP_BOOL_FALSE3", "no")

	t.Logf("APP_BOOL_TRUE1: %s", os.Getenv("APP_BOOL_TRUE1"))
	t.Logf("APP_BOOL_TRUE2: %s", os.Getenv("APP_BOOL_TRUE2"))
	t.Logf("APP_BOOL_TRUE3: %s", os.Getenv("APP_BOOL_TRUE3"))

	envSrc := sources.NewEnvSource("APP_")
	envData, err := envSrc.Load()
	if err != nil {
		t.Fatalf("Failed to load from env source: %v", err)
	}
	t.Logf("Env data: %#v", envData)

	cfg, err := Load[BoolConfig](WithEnv("APP_"))
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	t.Logf("Bool.True1: %v", cfg.Bool.True1)
	t.Logf("Bool.True2: %v", cfg.Bool.True2)
	t.Logf("Bool.True3: %v", cfg.Bool.True3)
	t.Logf("Bool.False1: %v", cfg.Bool.False1)
	t.Logf("Bool.False2: %v", cfg.Bool.False2)
	t.Logf("Bool.False3: %v", cfg.Bool.False3)

	if !cfg.Bool.True1 {
		t.Errorf("Expected Bool.True1 to be true, got %v", cfg.Bool.True1)
	}
	if !cfg.Bool.True2 {
		t.Errorf("Expected Bool.True2 to be true, got %v", cfg.Bool.True2)
	}
	if !cfg.Bool.True3 {
		t.Errorf("Expected Bool.True3 to be true, got %v", cfg.Bool.True3)
	}
	if cfg.Bool.False1 {
		t.Errorf("Expected Bool.False1 to be false, got %v", cfg.Bool.False1)
	}
	if cfg.Bool.False2 {
		t.Errorf("Expected Bool.False2 to be false, got %v", cfg.Bool.False2)
	}
	if cfg.Bool.False3 {
		t.Errorf("Expected Bool.False3 to be false, got %v", cfg.Bool.False3)
	}

	os.Unsetenv("APP_BOOL_TRUE1")
	os.Unsetenv("APP_BOOL_TRUE2")
	os.Unsetenv("APP_BOOL_TRUE3")
	os.Unsetenv("APP_BOOL_FALSE1")
	os.Unsetenv("APP_BOOL_FALSE2")
	os.Unsetenv("APP_BOOL_FALSE3")
}

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

func TestLoadConfigFromDotEnv(t *testing.T) {
    type EnvCfg struct {
        DB struct {
            Host string `config:"host"`
            Port int    `config:"port"`
        } `config:"db"`
        Port int `config:"port"`
    }

    tempDir := t.TempDir()
    envPath := filepath.Join(tempDir, ".env")
    content := `
# Comment
DB_HOST=dotenv-host
DB_PORT=6543 # inline comment
PORT="8080"
export UNUSED=1
`
    if err := os.WriteFile(envPath, []byte(content), 0644); err != nil {
        t.Fatalf("Failed to write dotenv file: %v", err)
    }

    cfg, err := Load[EnvCfg](WithDotEnv(envPath))
    if err != nil {
        t.Fatalf("Failed to load config from .env: %v", err)
    }

    if cfg.DB.Host != "dotenv-host" {
        t.Errorf("Expected DB.Host to be 'dotenv-host', got '%s'", cfg.DB.Host)
    }
    if cfg.DB.Port != 6543 {
        t.Errorf("Expected DB.Port to be 6543, got %d", cfg.DB.Port)
    }
    if cfg.Port != 8080 {
        t.Errorf("Expected Port to be 8080, got %d", cfg.Port)
    }
}

func TestLoadConfigFromDotEnvAndEnvMerge(t *testing.T) {
    type EnvCfg struct {
        DB struct {
            Host string `config:"host"`
            Port int    `config:"port"`
        } `config:"db"`
        Port int `config:"port"`
    }

    tempDir := t.TempDir()
    envPath := filepath.Join(tempDir, ".env")
    content := `
DB_HOST=dotenv-host
DB_PORT=6543
PORT=8080
`
    if err := os.WriteFile(envPath, []byte(content), 0644); err != nil {
        t.Fatalf("Failed to write dotenv file: %v", err)
    }

    os.Setenv("APP_DB_HOST", "env-host")
    os.Setenv("APP_PORT", "9090")

    cfg, err := Load[EnvCfg](WithDotEnv(envPath), WithEnv("APP_"))
    if err != nil {
        t.Fatalf("Failed to load config from .env + env: %v", err)
    }

    // Env should override dotenv where overlapping
    if cfg.DB.Host != "env-host" {
        t.Errorf("Expected DB.Host to be 'env-host', got '%s'", cfg.DB.Host)
    }
    if cfg.DB.Port != 6543 {
        t.Errorf("Expected DB.Port to be 6543, got %d", cfg.DB.Port)
    }
    if cfg.Port != 9090 {
        t.Errorf("Expected Port to be 9090, got %d", cfg.Port)
    }

    os.Unsetenv("APP_DB_HOST")
    os.Unsetenv("APP_PORT")
}

func TestMaskedJSONRedactsSecrets(t *testing.T) {
    type S struct {
        DB struct {
            Host string `config:"host"`
            Pass string `config:"pass" secret:"true"`
        } `config:"db"`
        Token string `config:"token" secret:"true"`
        Port  int    `config:"port"`
    }
    cfg := S{}
    cfg.DB.Host = "localhost"
    cfg.DB.Pass = "super-secret"
    cfg.Token = "tkn"
    cfg.Port = 8080

    masked, err := MaskedJSON(cfg)
    if err != nil {
        t.Fatalf("MaskedJSON failed: %v", err)
    }
    if strings.Contains(masked, "super-secret") || strings.Contains(masked, "tkn") {
        t.Errorf("MaskedJSON leaked secret values: %s", masked)
    }
    if !strings.Contains(masked, "\"***\"") {
        t.Errorf("MaskedJSON does not contain masked placeholders: %s", masked)
    }
    if !strings.Contains(masked, "localhost") || !strings.Contains(masked, "8080") {
        t.Errorf("MaskedJSON missing non-secret fields: %s", masked)
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

func TestRequiredFieldsMissing(t *testing.T) {
    type ReqConfig struct {
        DB struct {
            Host string `config:"host" required:"true"`
        } `config:"db"`
        Port int `config:"port" required:"true"`
    }

    // Ensure env is clean
    os.Unsetenv("APP_DB_HOST")
    os.Unsetenv("APP_PORT")

    _, err := Load[ReqConfig](WithEnv("APP_"))
    if err == nil {
        t.Error("Expected error for missing required fields, got nil")
    }
}

func TestRequiredFieldsPresent(t *testing.T) {
    type ReqConfig struct {
        DB struct {
            Host string `config:"host" required:"true"`
        } `config:"db"`
        Port int `config:"port" required:"true"`
    }

    os.Setenv("APP_DB_HOST", "localhost")
    os.Setenv("APP_PORT", "3000")

    cfg, err := Load[ReqConfig](WithEnv("APP_"))
    if err != nil {
        t.Fatalf("Did not expect error when required fields are present, got: %v", err)
    }

    if cfg.DB.Host != "localhost" || cfg.Port != 3000 {
        t.Errorf("Unexpected values: DB.Host=%s Port=%d", cfg.DB.Host, cfg.Port)
    }

    os.Unsetenv("APP_DB_HOST")
    os.Unsetenv("APP_PORT")
}
