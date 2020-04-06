package tools

// CloneMapString creates a new strings map by cloning the given map. It creates
// a shallow copy of the input map.
func CloneMapString(value map[string]string) map[string]string {
	cloned := make(map[string]string)
	for k, v := range value {
		cloned[k] = v
	}
	return cloned
}
