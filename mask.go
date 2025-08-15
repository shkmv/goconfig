package goconfig

import (
    "github.com/shkmv/goconfig/internal"
)

// MaskedMap returns a nested map representation of the target struct with
// fields tagged `secret:"true"` masked.
func MaskedMap(target any) (map[string]any, error) {
    return internal.Sanitize(target)
}

// MaskedJSON returns a JSON string representation of the target struct with
// fields tagged `secret:"true"` masked. Suitable for safe logging.
func MaskedJSON(target any) (string, error) {
    return internal.MaskedJSON(target)
}

