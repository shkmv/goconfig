package sources

import (
	"os"
	"testing"
)

func TestLoadSourceEnv(t *testing.T) {
	os.Setenv("TEST_KEY", "value")
	os.Setenv("TEST_KEY2", "value2")

	envSrc := NewEnvSource("TEST_")
	out, err := envSrc.Load()
	if err != nil {
		t.Fatal(err)
	}

	if out["key"] != "value" {
		t.Errorf("expected key to be value, got %s", out["key"])
	}

	if out["key2"] != "value2" {
		t.Errorf("expected key2 to be value2, got %s", out["key2"])
	}
}

func TestLoadSourceEnvWithNestedKeys(t *testing.T) {
	os.Unsetenv("TEST_KEY")
	os.Unsetenv("TEST_KEY2")

	os.Setenv("APP_DB_HOST", "localhost")
	os.Setenv("APP_DB_PORT", "5432")
	os.Setenv("APP_PORT", "3000")

	envSrc := NewEnvSource("APP_")
	out, err := envSrc.Load()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("out: %#v", out)

	dbMap, ok := out["db"].(map[string]any)
	if !ok {
		t.Fatalf("expected out[\"db\"] to be map[string]any, got %T", out["db"])
	}

	if dbMap["host"] != "localhost" {
		t.Errorf("expected db.host to be localhost, got %v", dbMap["host"])
	}

	if dbMap["port"] != "5432" {
		t.Errorf("expected db.port to be 5432, got %v", dbMap["port"])
	}

	if out["port"] != "3000" {
		t.Errorf("expected port to be 3000, got %v", out["port"])
	}
}
