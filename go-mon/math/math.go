package math

import "cmp"

// Signed is a constraint that permits any signed integer or float type.
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64
}

// Abs returns the absolute value of a signed numeric type.
// Uses Signed constraint to ensure only signed types can be used.
func Abs[T Signed](val T) T {
	if val < 0 {
		return -val
	}
	return val
}

// Max returns the maximum of two ordered values using cmp.Compare (Go 1.21+).
func Max[T cmp.Ordered](a, b T) T {
	if cmp.Compare(a, b) < 0 {
		return b
	}
	return a
}

// Min returns the minimum of two ordered values using cmp.Compare (Go 1.21+).
func Min[T cmp.Ordered](a, b T) T {
	if cmp.Compare(a, b) < 0 {
		return a
	}
	return b
}

// AbsInt32 returns the absolute value of an int32.
// Deprecated: Use Abs instead.
func AbsInt32(val int32) int32 {
	if val < 0 {
		return -val
	}
	return val
}

// AbsInt64 returns the absolute value of an int64.
// Deprecated: Use Abs instead.
func AbsInt64(val int64) int64 {
	if val < 0 {
		return -val
	}
	return val
}
