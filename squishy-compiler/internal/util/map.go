package util

import "strings"

// Public Functions
func ConcatStringIndexedMapToIndexOnlyString[V any](m map[string]V, space string) string {
	if len(m) == 0 {
		return ""
	}

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	return strings.Join(keys, space) // TS is nice
}
