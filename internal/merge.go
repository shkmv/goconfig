package internal

// Merge merges two maps into a new map.
func Merge(dst, src map[string]any) map[string]any {
	for k, v := range src {
		if vMap, ok := v.(map[string]any); ok {
			if dstMap, ok := dst[k].(map[string]any); ok {
				dst[k] = Merge(dstMap, vMap)
			} else {
				dst[k] = Merge(make(map[string]any), vMap)
			}
		} else {
			dst[k] = v
		}
	}
	return dst
}
