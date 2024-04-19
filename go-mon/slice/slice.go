package slice

func ToMap[T any, K comparable](slice []T, keyFunc func(T) K) map[K]T {
	result := make(map[K]T, len(slice))
	for i := range slice {
		result[keyFunc(slice[i])] = slice[i]
	}
	return result
}

func ToMapBool(slice []string) map[string]bool {
	result := make(map[string]bool, len(slice))
	for i := range slice {
		result[slice[i]] = true
	}
	return result
}

func ToSlice[T any, X any](d []T, fn func(d T) X) []X {
	if d == nil {
		return nil
	}
	dd := make([]X, len(d))
	for i, v := range d {
		dd[i] = fn(v)
	}
	return dd
}
