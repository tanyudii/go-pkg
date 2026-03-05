package slice

import (
	"cmp"
	"slices"
)

// ToMap converts a slice to a map using the provided key function.
// If duplicate keys exist, later values will overwrite earlier ones.
func ToMap[T any, K comparable](slice []T, keyFunc func(T) K) map[K]T {
	result := make(map[K]T, len(slice))
	for i := range slice {
		result[keyFunc(slice[i])] = slice[i]
	}
	return result
}

// ToMapBool converts a string slice to a map[string]bool with all values set to true.
// Useful for fast membership testing.
func ToMapBool(slice []string) map[string]bool {
	result := make(map[string]bool, len(slice))
	for i := range slice {
		result[slice[i]] = true
	}
	return result
}

// ToSlice transforms a slice of type T to a slice of type X using the provided function.
// Returns nil if the input slice is nil.
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

// Contains checks if a slice contains an element using the stdlib slices.Contains (Go 1.21+).
func Contains[T comparable](slice []T, elem T) bool {
	return slices.Contains(slice, elem)
}

// Index returns the index of the first occurrence of elem in slice, or -1 if not present.
// Uses stdlib slices.Index (Go 1.21+).
func Index[T comparable](slice []T, elem T) int {
	return slices.Index(slice, elem)
}

// Reverse reverses a slice in place using stdlib slices.Reverse (Go 1.21+).
func Reverse[T any](slice []T) {
	slices.Reverse(slice)
}

// Max returns the maximum element in a slice using stdlib slices.Max (Go 1.21+).
// Returns zero value if slice is empty.
func Max[T cmp.Ordered](slice []T) T {
	if len(slice) == 0 {
		var zero T
		return zero
	}
	return slices.Max(slice)
}

// Min returns the minimum element in a slice using stdlib slices.Min (Go 1.21+).
// Returns zero value if slice is empty.
func Min[T cmp.Ordered](slice []T) T {
	if len(slice) == 0 {
		var zero T
		return zero
	}
	return slices.Min(slice)
}

// Sort sorts a slice in ascending order using stdlib slices.Sort (Go 1.21+).
func Sort[T cmp.Ordered](slice []T) {
	slices.Sort(slice)
}

// SortFunc sorts a slice using the provided comparison function using stdlib slices.SortFunc (Go 1.21+).
// The comparison function should return:
//   - a negative number if a < b
//   - zero if a == b
//   - a positive number if a > b
func SortFunc[T any](slice []T, cmp func(a, b T) int) {
	slices.SortFunc(slice, cmp)
}

// SortBy sorts a slice using a boolean less-than function.
// This is a convenience wrapper that converts a boolean predicate to the int comparison needed by slices.SortFunc.
func SortBy[T any](slice []T, less func(a, b T) bool) {
	slices.SortFunc(slice, func(a, b T) int {
		if less(a, b) {
			return -1
		}
		if less(b, a) {
			return 1
		}
		return 0
	})
}

// Clone returns a shallow copy of a slice using stdlib slices.Clone (Go 1.21+).
func Clone[T any](slice []T) []T {
	return slices.Clone(slice)
}

// Equal checks if two slices are equal using stdlib slices.Equal (Go 1.21+).
func Equal[T comparable](a, b []T) bool {
	return slices.Equal(a, b)
}

// Chunk splits a slice into chunks of the specified size.
// The last chunk may be smaller than size if the slice length is not a multiple of size.
func Chunk[T any](slice []T, size int) [][]T {
	if size <= 0 {
		return nil
	}
	var chunks [][]T
	for i := 0; i < len(slice); i += size {
		end := i + size
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}

// Filter returns a new slice containing only the elements that satisfy the predicate.
func Filter[T any](slice []T, predicate func(T) bool) []T {
	result := make([]T, 0, len(slice))
	for _, v := range slice {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

// Unique removes duplicate elements from a slice while preserving order.
func Unique[T comparable](slice []T) []T {
	seen := make(map[T]struct{}, len(slice))
	result := make([]T, 0, len(slice))
	for _, v := range slice {
		if _, exists := seen[v]; !exists {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}

// Flatten flattens a 2D slice into a 1D slice.
func Flatten[T any](slice [][]T) []T {
	totalLen := 0
	for _, s := range slice {
		totalLen += len(s)
	}
	result := make([]T, 0, totalLen)
	for _, s := range slice {
		result = append(result, s...)
	}
	return result
}
