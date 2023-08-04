package slice

func FirstVal[T comparable](args ...T) T {
	var zero T
	for _, a := range args {
		if a != zero {
			return a
		}
	}
	return zero
}
