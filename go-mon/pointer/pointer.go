package pointer

import "reflect"

func Val[T any](val T) *T {
	return &val
}

func Extract[T any](val *T) T {
	var result T
	if val == nil {
		return result
	}
	return *val
}

func EmptyNil[T any](val T) *T {
	if reflect.ValueOf(&val).Elem().IsZero() {
		return nil
	}
	return &val
}
