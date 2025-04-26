package internal

import (
	"fmt"
	"reflect"
	"strings"
)

// Bind recursively binds data from a map to the fields of a target struct
// based on `config` tags.
func Bind(data map[string]any, target any) error {
	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("target must be a non-nil pointer, got %T", target)
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("target pointer must point to a struct, got %s", v.Kind())
	}
	t := v.Type()

	for i := range t.NumField() {
		field := t.Field(i)
		tag := field.Tag.Get("config")
		if tag == "" || !field.IsExported() {
			continue
		}

		keys := strings.Split(tag, ".")
		val, ok := lookup(data, keys)
		if !ok {
			continue
		}

		fieldVal := v.Field(i)
		if !fieldVal.CanSet() {
			continue
		}

		if fieldVal.Kind() == reflect.Struct {
			if subData, ok := val.(map[string]any); ok {
				if !fieldVal.CanAddr() {
					return fmt.Errorf("cannot get address of field %s to bind nested struct", field.Name)
				}
				if err := Bind(subData, fieldVal.Addr().Interface()); err != nil {
					return fmt.Errorf("error binding nested struct field %s: %w", field.Name, err)
				}
			} else if val != nil {
				return fmt.Errorf("type mismatch for field %s: expected map[string]any for nested struct, got %T", field.Name, val)
			}
		} else if fieldVal.Kind() == reflect.Ptr && fieldVal.Elem().Kind() == reflect.Struct {
			if subData, ok := val.(map[string]any); ok {
				if fieldVal.IsNil() {
					newStruct := reflect.New(fieldVal.Type().Elem())
					fieldVal.Set(newStruct)
				}
				if err := Bind(subData, fieldVal.Interface()); err != nil {
					return fmt.Errorf("error binding nested pointer field %s: %w", field.Name, err)
				}
			} else if val != nil {
				return fmt.Errorf("type mismatch for pointer field %s: expected map[string]any for nested struct pointer, got %T", field.Name, val)
			}
		} else {
			if err := assing(fieldVal, val); err != nil {
				return fmt.Errorf("error assigning value to field %s: %w", field.Name, err)
			}
		}
	}

	return nil
}

func lookup(data map[string]any, keys []string) (any, bool) {
	if len(keys) == 0 || data == nil {
		return nil, false
	}
	key := keys[0]
	val, ok := data[key]
	if !ok {
		return nil, false
	}

	if len(keys) == 1 {
		return val, true
	}

	if subdata, ok := val.(map[string]any); ok {
		return lookup(subdata, keys[1:])
	}

	return nil, false
}

// assing assigns a value to a reflect.Value field, handling basic type conversions.
func assing(fieldVal reflect.Value, val any) error {
	if val == nil {
		// Cannot assign nil directly unless it's a pointer, slice, map, etc.
		switch fieldVal.Kind() {
		case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
			if fieldVal.IsNil() {
				return nil
			}
			return fmt.Errorf("cannot assign nil to non-nil field %s", fieldVal.Type())
		default:
			return fmt.Errorf("cannot assign nil to field %s of kind %s", fieldVal.Type(), fieldVal.Kind())
		}
		// TODO: Implement support for setting nil for nilable types if this feature will be needed.
		// zeroValue := reflect.Zero(fieldVal.Type())
		// fieldVal.Set(zeroValue)
		// return nil
	}

	if strVal, ok := val.(string); ok {
		switch fieldVal.Kind() {
		case reflect.String:
			fieldVal.SetString(strVal)
			return nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			var intVal int64
			if _, err := fmt.Sscanf(strVal, "%d", &intVal); err != nil {
				return fmt.Errorf("cannot convert string '%s' to int: %w", strVal, err)
			}
			if fieldVal.OverflowInt(intVal) {
				return fmt.Errorf("integer overflow assigning %d to %s", intVal, fieldVal.Type())
			}
			fieldVal.SetInt(intVal)
			return nil
		case reflect.Float32, reflect.Float64:
			var floatVal float64
			if _, err := fmt.Sscanf(strVal, "%f", &floatVal); err != nil {
				return fmt.Errorf("cannot convert string '%s' to float: %w", strVal, err)
			}
			if fieldVal.Kind() == reflect.Float32 && fieldVal.OverflowFloat(floatVal) {
				return fmt.Errorf("float32 overflow assigning %f", floatVal)
			}
			fieldVal.SetFloat(floatVal)
			return nil
		case reflect.Bool:
			var boolVal bool
			if strVal == "true" || strVal == "1" || strVal == "yes" || strVal == "y" {
				boolVal = true
			} else if strVal == "false" || strVal == "0" || strVal == "no" || strVal == "n" {
				boolVal = false
			} else {
				return fmt.Errorf("cannot convert string '%s' to bool", strVal)
			}
			fieldVal.SetBool(boolVal)
			return nil
		}
	}

	valValue := reflect.ValueOf(val)

	if fieldVal.Kind() != valValue.Kind() {
		// Allow assigning float64 (common from JSON/YAML) to int types
		if (fieldVal.Kind() == reflect.Int || fieldVal.Kind() == reflect.Int8 || fieldVal.Kind() == reflect.Int16 || fieldVal.Kind() == reflect.Int32 || fieldVal.Kind() == reflect.Int64) && valValue.Kind() == reflect.Float64 {
		} else if (fieldVal.Kind() == reflect.Float32 || fieldVal.Kind() == reflect.Float64) && (valValue.Kind() == reflect.Float64 || valValue.Kind() == reflect.Float32) {
		} else if fieldVal.Kind() == reflect.String && valValue.Kind() == reflect.String {
		} else if fieldVal.Kind() == reflect.Bool && valValue.Kind() == reflect.Bool {
		} else {
			return fmt.Errorf("type mismatch: cannot assign %T to field of type %s", val, fieldVal.Type())
		}
	}

	switch fieldVal.Kind() {
	case reflect.String:
		if s, ok := val.(string); ok {
			fieldVal.SetString(s)
		} else {
			return fmt.Errorf("value for string field is not a string: %T", val)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// Handle potential float64 input from map[string]any (JSON)
		if f64, ok := val.(float64); ok {
			if fieldVal.OverflowInt(int64(f64)) {
				return fmt.Errorf("integer overflow assigning %f to %s", f64, fieldVal.Type())
			}
			fieldVal.SetInt(int64(f64))
		} else if i64, ok := val.(int64); ok {
			if fieldVal.OverflowInt(i64) {
				return fmt.Errorf("integer overflow assigning %d to %s", i64, fieldVal.Type())
			}
			fieldVal.SetInt(i64)
		} else if i, ok := val.(int); ok {
			if fieldVal.OverflowInt(int64(i)) {
				return fmt.Errorf("integer overflow assigning %d to %s", i, fieldVal.Type())
			}
			fieldVal.SetInt(int64(i))
		} else {
			return fmt.Errorf("value for integer field is not a number: %T", val)
		}
	case reflect.Float32, reflect.Float64:
		var f64 float64
		var ok bool
		if f64, ok = val.(float64); !ok {
			// Try int conversion if not float64
			var i int
			if i, ok = val.(int); ok {
				f64 = float64(i)
			} else {
				var i64 int64
				if i64, ok = val.(int64); ok {
					f64 = float64(i64)
				} else {
					return fmt.Errorf("value for float field is not a number: %T", val)
				}
			}
		}
		if fieldVal.Kind() == reflect.Float32 && fieldVal.OverflowFloat(f64) {
			return fmt.Errorf("float32 overflow assigning %f", f64)
		}
		fieldVal.SetFloat(f64)

	case reflect.Bool:
		if b, ok := val.(bool); ok {
			fieldVal.SetBool(b)
		} else if s, ok := val.(string); ok {
			sLower := strings.ToLower(s)
			if sLower == "true" || sLower == "1" || sLower == "yes" || sLower == "y" {
				fieldVal.SetBool(true)
			} else if sLower == "false" || sLower == "0" || sLower == "no" || sLower == "n" {
				fieldVal.SetBool(false)
			} else {
				return fmt.Errorf("cannot convert string '%s' to bool", s)
			}
		} else {
			return fmt.Errorf("value for bool field is not a bool or string: %T", val)
		}
	default:
		return fmt.Errorf("unsupported type %s in assing function", fieldVal.Type())
	}
	return nil
}
