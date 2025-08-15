package sources

// setNestedValue sets a nested value in a map using the provided keys path.
func setNestedValue(data map[string]any, keys []string, value any) {
    if len(keys) == 0 {
        return
    }

    key := keys[0]
    if len(keys) == 1 {
        data[key] = value
        return
    }

    if _, ok := data[key]; !ok {
        data[key] = make(map[string]any)
    } else if _, ok := data[key].(map[string]any); !ok {
        data[key] = make(map[string]any)
    }

    setNestedValue(data[key].(map[string]any), keys[1:], value)
}

