package internal

import (
	"reflect"
	"testing"
)

func TestMerge(t *testing.T) {
	dst := map[string]any{"a": 1, "b": map[string]any{"c": 2}}
	src := map[string]any{"b": map[string]any{"d": 3}, "e": 4}
	expected := map[string]any{"a": 1, "b": map[string]any{"c": 2, "d": 3}, "e": 4}

	result := Merge(dst, src)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}
