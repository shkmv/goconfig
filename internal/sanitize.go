package internal

import (
    "encoding/json"
    "fmt"
    "reflect"
    "strings"
)

// Sanitize traverses the target struct using `config` tags and returns a nested
// map representation with fields marked `secret:"true"` masked out.
// The target can be a struct or a pointer to struct.
func Sanitize(target any) (map[string]any, error) {
    v := reflect.ValueOf(target)
    if v.Kind() == reflect.Ptr {
        if v.IsNil() {
            return nil, fmt.Errorf("target must not be nil")
        }
        v = v.Elem()
    }
    if v.Kind() != reflect.Struct {
        return nil, fmt.Errorf("target must be a struct or pointer to struct, got %T", target)
    }

    out := make(map[string]any)
    if err := sanitizeStruct(v, nil, out); err != nil {
        return nil, err
    }
    return out, nil
}

// MaskedJSON returns a JSON string representation of the sanitized struct.
func MaskedJSON(target any) (string, error) {
    m, err := Sanitize(target)
    if err != nil {
        return "", err
    }
    b, err := json.Marshal(m)
    if err != nil {
        return "", fmt.Errorf("marshal masked json: %w", err)
    }
    return string(b), nil
}

func sanitizeStruct(v reflect.Value, prefix []string, out map[string]any) error {
    t := v.Type()
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        if !field.IsExported() {
            continue
        }
        tag := field.Tag.Get("config")
        if tag == "" {
            continue
        }
        // path from tag and prefix
        path := append([]string{}, prefix...)
        path = append(path, strings.Split(tag, ".")...)

        // handle nested
        fv := v.Field(i)
        switch fv.Kind() {
        case reflect.Struct:
            if err := sanitizeStruct(fv, path, out); err != nil {
                return err
            }
            continue
        case reflect.Ptr:
            if fv.IsNil() {
                // nothing to add
                continue
            }
            if fv.Elem().Kind() == reflect.Struct {
                if err := sanitizeStruct(fv.Elem(), path, out); err != nil {
                    return err
                }
                continue
            }
        }

        // leaf value
        masked := field.Tag.Get("secret")
        if isTruthy(masked) {
            setNested(out, path, "***")
            continue
        }

        // Use the field's current value for non-secret entries
        setNested(out, path, fv.Interface())
    }
    return nil
}

func setNested(dst map[string]any, keys []string, val any) {
    if len(keys) == 0 {
        return
    }
    if len(keys) == 1 {
        dst[keys[0]] = val
        return
    }
    key := keys[0]
    child, ok := dst[key].(map[string]any)
    if !ok {
        child = make(map[string]any)
        dst[key] = child
    }
    setNested(child, keys[1:], val)
}

