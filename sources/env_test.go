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
