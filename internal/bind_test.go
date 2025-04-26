package internal

import (
	"math"
	"reflect"
	"strings"
	"testing"
)

func TestBind(t *testing.T) {
	type Nested struct {
		Value string `config:"value"`
	}

	type Config struct {
		Name       string  `config:"server.name"`
		Port       int     `config:"server.port"`
		Enabled    bool    `config:"feature.enabled"`
		Nested     Nested  `config:"nested"`
		FloatVal   float64 `config:"floatValue"`
		Missing    string  `config:"missing.key"`
		NoTag      string
		unexported string `config:"unexported"` // Unexported field, should not be set
	}

	data := map[string]any{
		"server": map[string]any{
			"name": "TestServer",
			"port": 8080.0, // JSON/YAML typically decode numbers as float64
		},
		"feature": map[string]any{
			"enabled": true,
		},
		"nested": map[string]any{
			"value": "NestedValue",
		},
		"floatValue": 123.45,
		"unexported": "should not be set",
		"NoTag":      "should not be set",
	}

	t.Run("Successful Bind", func(t *testing.T) {
		target := Config{NoTag: "initial"}
		err := Bind(data, &target)
		if err != nil {
			t.Fatalf("Bind failed: %v", err)
		}

		expected := Config{
			Name:       "TestServer",
			Port:       8080,
			Enabled:    true,
			Nested:     Nested{Value: "NestedValue"},
			FloatVal:   123.45,
			Missing:    "",        // Should remain default value
			NoTag:      "initial", // Should remain initial value
			unexported: "",        // Should remain default value
		}

		if !reflect.DeepEqual(target, expected) {
			t.Errorf("Bind result mismatch.\nGot:  %#v\nWant: %#v", target, expected)
		}
	})

	t.Run("Nil Target", func(t *testing.T) {
		var target *Config = nil
		err := Bind(data, target)
		if err == nil {
			t.Errorf("Expected an error for nil target, got nil")
		} else {
			expectedErr := "target must be a non-nil pointer, got *internal.Config"
			if err.Error() != expectedErr {
				t.Errorf("Expected error '%s', got '%s'", expectedErr, err.Error())
			}
		}
	})

	t.Run("Non-Pointer Target", func(t *testing.T) {
		target := Config{}
		err := Bind(data, target)
		if err == nil {
			t.Errorf("Expected an error for non-pointer target, got nil")
		} else {
			expectedErr := "target must be a non-nil pointer, got internal.Config"
			if err.Error() != expectedErr {
				t.Errorf("Expected error '%s', got '%s'", expectedErr, err.Error())
			}
		}
	})

	t.Run("Empty Data", func(t *testing.T) {
		target := Config{Name: "Default"}
		emptyData := map[string]any{}
		err := Bind(emptyData, &target)
		if err != nil {
			t.Fatalf("Bind with empty data failed: %v", err)
		}
		expected := Config{Name: "Default"}
		if !reflect.DeepEqual(target, expected) {
			t.Errorf("Bind with empty data mismatch.\nGot:  %#v\nWant: %#v", target, expected)
		}
	})
}

func TestLookup(t *testing.T) {
	data := map[string]any{
		"a": 1,
		"b": map[string]any{
			"c": 2,
			"d": map[string]any{
				"e": 3,
			},
		},
		"f": "hello",
	}

	testCases := []struct {
		name     string
		keys     []string
		expected any
		found    bool
	}{
		{"Top Level Exists", []string{"a"}, 1, true},
		{"Nested Level 1 Exists", []string{"b", "c"}, 2, true},
		{"Nested Level 2 Exists", []string{"b", "d", "e"}, 3, true},
		{"Top Level Not Exists", []string{"z"}, nil, false},
		{"Nested Level 1 Not Exists", []string{"b", "z"}, nil, false},
		{"Nested Level 2 Not Exists", []string{"b", "d", "z"}, nil, false},
		{"Incorrect Path Type", []string{"f", "g"}, nil, false},
		{"Empty Keys", []string{}, nil, false},
		{"Nil Data", nil, nil, false}, // Test with nil map passed explicitly later
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			val, ok := lookup(data, tc.keys)
			if ok != tc.found {
				t.Errorf("lookup(%v) found mismatch: got %v, want %v", tc.keys, ok, tc.found)
			}
			if ok && !reflect.DeepEqual(val, tc.expected) {
				t.Errorf("lookup(%v) value mismatch: got %v (%T), want %v (%T)", tc.keys, val, val, tc.expected, tc.expected)
			}
		})
	}

	t.Run("Nil Data Input", func(t *testing.T) {
		val, ok := lookup(nil, []string{"a"})
		if ok {
			t.Errorf("lookup(nil, keys) found mismatch: got true, want false")
		}
		if val != nil {
			t.Errorf("lookup(nil, keys) value mismatch: got %v, want nil", val)
		}
	})
}

func createTestField(value any) reflect.Value {
	structType := reflect.StructOf([]reflect.StructField{
		{Name: "Field", Type: reflect.TypeOf(value)},
	})
	instance := reflect.New(structType).Elem()
	return instance.Field(0)
}

func TestAssing(t *testing.T) {
	testCases := []struct {
		name        string
		fieldType   reflect.Kind
		initialVal  any // Value to create the field with
		assignVal   any // Value to assign using assing
		expectedVal any // Expected value after assignment
		expectError bool
	}{
		{"Assign String", reflect.String, "", "test", "test", false},
		{"Assign Int from Float64", reflect.Int, 0, 123.0, int64(123), false}, // Common case for JSON numbers
		{"Assign Int64 from Float64", reflect.Int64, int64(0), 456.0, int64(456), false},
		{"Assign Float64", reflect.Float64, 0.0, 789.12, 789.12, false},
		{"Assign Float32", reflect.Float32, float32(0.0), 12.34, float64(12.34), false}, // Note: SetFloat takes float64
		{"Assign Bool True", reflect.Bool, false, true, true, false},
		{"Assign Bool False", reflect.Bool, true, false, false, false},
		{"Unsupported Type Slice", reflect.Slice, []int{}, []int{1}, nil, true}, // Example of an unsupported type
		{"Unsupported Type Struct", reflect.Struct, struct{}{}, struct{}{}, nil, true},
		{"Assign Int from String", reflect.Int, 0, "123", int64(123), false},
		{"Assign Int64 from String", reflect.Int64, int64(0), "456", int64(456), false},
		{"Assign Float64 from String", reflect.Float64, 0.0, "789.12", 789.12, false},
		{"Assign Float32 from String", reflect.Float32, float32(0.0), "12.34", float64(12.34), false},
		{"Assign Bool True from String", reflect.Bool, false, "true", true, false},
		{"Assign Bool False from String", reflect.Bool, true, "false", false, false},
		{"Assign Bool True from String (1)", reflect.Bool, false, "1", true, false},
		{"Assign Bool False from String (0)", reflect.Bool, true, "0", false, false},
		{"Assign Bool True from String (yes)", reflect.Bool, false, "yes", true, false},
		{"Assign Bool False from String (no)", reflect.Bool, true, "no", false, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a settable field value of the correct type
			fieldVal := createTestField(tc.initialVal)
			if fieldVal.Kind() != tc.fieldType {
				// Adjust if createTestField logic changes
				t.Skipf("Test setup issue: field kind mismatch for %s", tc.name)
			}

			err := assing(fieldVal, tc.assignVal)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected an error for type %s, but got nil", tc.fieldType)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for type %s: %v", tc.fieldType, err)
				}
				var finalVal any
				skipDeepEqual := false // Flag to skip DeepEqual for float types
				switch tc.fieldType {
				case reflect.String:
					finalVal = fieldVal.String()
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					finalVal = fieldVal.Int()
				case reflect.Float32, reflect.Float64:
					finalVal = fieldVal.Float()
					// Epsilon comparison for floats
					const epsilon = 1e-6
					gotFloat := finalVal.(float64)
					expectedFloat := tc.expectedVal.(float64)
					if math.Abs(gotFloat-expectedFloat) > epsilon {
						t.Errorf("Value mismatch for type %s: got %v (%T), want %v (%T)",
							tc.fieldType, finalVal, finalVal, tc.expectedVal, tc.expectedVal)
					}
					skipDeepEqual = true
				case reflect.Bool:
					finalVal = fieldVal.Bool()
				default:
					t.Fatalf("Unhandled field type in test check: %s", tc.fieldType)

				}

				// Use reflect.DeepEqual for non-float types
				if !skipDeepEqual && !reflect.DeepEqual(finalVal, tc.expectedVal) {
					t.Errorf("Value mismatch for type %s: got %v (%T), want %v (%T)",
						tc.fieldType, finalVal, finalVal, tc.expectedVal, tc.expectedVal)
				}
			}
		})
	}

	t.Run("Invalid String to Int Conversion Error", func(t *testing.T) {
		fieldVal := createTestField(0)
		err := assing(fieldVal, "not-an-int")
		if err == nil {
			t.Errorf("Expected assing to return an error for invalid string to int conversion, but it returned nil")
		} else {
			expectedErrorSubstring := "cannot convert string"
			if !strings.Contains(err.Error(), expectedErrorSubstring) {
				t.Errorf("Expected error message to contain '%s', got '%s'", expectedErrorSubstring, err.Error())
			}
			t.Logf("Received expected error from assing: %v", err)
		}
	})
}
