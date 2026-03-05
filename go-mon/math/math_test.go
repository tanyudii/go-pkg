package math

import (
	"testing"
)

func TestAbs(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{"positive int", int(5), int(5)},
		{"negative int", int(-5), int(5)},
		{"zero int", int(0), int(0)},
		{"positive int8", int8(10), int8(10)},
		{"negative int8", int8(-10), int8(10)},
		{"zero int8", int8(0), int8(0)},
		{"positive int16", int16(100), int16(100)},
		{"negative int16", int16(-100), int16(100)},
		{"positive int32", int32(1000), int32(1000)},
		{"negative int32", int32(-1000), int32(1000)},
		{"positive int64", int64(10000), int64(10000)},
		{"negative int64", int64(-10000), int64(10000)},
		{"positive float32", float32(1.5), float32(1.5)},
		{"negative float32", float32(-1.5), float32(1.5)},
		{"positive float64", float64(2.7), float64(2.7)},
		{"negative float64", float64(-2.7), float64(2.7)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.input.(type) {
			case int:
				result := Abs(v)
				if result != tt.expected.(int) {
					t.Errorf("Abs(%d) = %d; want %d", v, result, tt.expected.(int))
				}
			case int8:
				result := Abs(v)
				if result != tt.expected.(int8) {
					t.Errorf("Abs(%d) = %d; want %d", v, result, tt.expected.(int8))
				}
			case int16:
				result := Abs(v)
				if result != tt.expected.(int16) {
					t.Errorf("Abs(%d) = %d; want %d", v, result, tt.expected.(int16))
				}
			case int32:
				result := Abs(v)
				if result != tt.expected.(int32) {
					t.Errorf("Abs(%d) = %d; want %d", v, result, tt.expected.(int32))
				}
			case int64:
				result := Abs(v)
				if result != tt.expected.(int64) {
					t.Errorf("Abs(%d) = %d; want %d", v, result, tt.expected.(int64))
				}
			case float32:
				result := Abs(v)
				if result != tt.expected.(float32) {
					t.Errorf("Abs(%f) = %f; want %f", v, result, tt.expected.(float32))
				}
			case float64:
				result := Abs(v)
				if result != tt.expected.(float64) {
					t.Errorf("Abs(%f) = %f; want %f", v, result, tt.expected.(float64))
				}
			}
		})
	}
}

func TestAbsInt32(t *testing.T) {
	tests := []struct {
		name  string
		input int32
		want  int32
	}{
		{"positive", 42, 42},
		{"negative", -42, 42},
		{"zero", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AbsInt32(tt.input); got != tt.want {
				t.Errorf("AbsInt32(%d) = %d; want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestAbsInt64(t *testing.T) {
	tests := []struct {
		name  string
		input int64
		want  int64
	}{
		{"positive", 420, 420},
		{"negative", -420, 420},
		{"zero", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AbsInt64(tt.input); got != tt.want {
				t.Errorf("AbsInt64(%d) = %d; want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestMax(t *testing.T) {
	tests := []struct {
		name     string
		a, b     interface{}
		expected interface{}
	}{
		{"int a > b", 5, 3, 5},
		{"int a < b", 3, 5, 5},
		{"int equal", 4, 4, 4},
		{"float32 a > b", float32(5.5), float32(3.3), float32(5.5)},
		{"float32 a < b", float32(3.3), float32(5.5), float32(5.5)},
		{"float64 a > b", 7.7, 4.4, 7.7},
		{"float64 a < b", 4.4, 7.7, 7.7},
		{"string a > b", "zzz", "aaa", "zzz"},
		{"string a < b", "aaa", "zzz", "zzz"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.a.(type) {
			case int:
				result := Max(v, tt.b.(int))
				if result != tt.expected.(int) {
					t.Errorf("Max(%d, %d) = %d; want %d", v, tt.b.(int), result, tt.expected.(int))
				}
			case float32:
				result := Max(v, tt.b.(float32))
				if result != tt.expected.(float32) {
					t.Errorf("Max(%f, %f) = %f; want %f", v, tt.b.(float32), result, tt.expected.(float32))
				}
			case float64:
				result := Max(v, tt.b.(float64))
				if result != tt.expected.(float64) {
					t.Errorf("Max(%f, %f) = %f; want %f", v, tt.b.(float64), result, tt.expected.(float64))
				}
			case string:
				result := Max(v, tt.b.(string))
				if result != tt.expected.(string) {
					t.Errorf("Max(%s, %s) = %s; want %s", v, tt.b.(string), result, tt.expected.(string))
				}
			}
		})
	}
}

func TestMin(t *testing.T) {
	tests := []struct {
		name     string
		a, b     interface{}
		expected interface{}
	}{
		{"int a > b", 5, 3, 3},
		{"int a < b", 3, 5, 3},
		{"int equal", 4, 4, 4},
		{"float32 a > b", float32(5.5), float32(3.3), float32(3.3)},
		{"float32 a < b", float32(3.3), float32(5.5), float32(3.3)},
		{"float64 a > b", 7.7, 4.4, 4.4},
		{"float64 a < b", 4.4, 7.7, 4.4},
		{"string a > b", "zzz", "aaa", "aaa"},
		{"string a < b", "aaa", "zzz", "aaa"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.a.(type) {
			case int:
				result := Min(v, tt.b.(int))
				if result != tt.expected.(int) {
					t.Errorf("Min(%d, %d) = %d; want %d", v, tt.b.(int), result, tt.expected.(int))
				}
			case float32:
				result := Min(v, tt.b.(float32))
				if result != tt.expected.(float32) {
					t.Errorf("Min(%f, %f) = %f; want %f", v, tt.b.(float32), result, tt.expected.(float32))
				}
			case float64:
				result := Min(v, tt.b.(float64))
				if result != tt.expected.(float64) {
					t.Errorf("Min(%f, %f) = %f; want %f", v, tt.b.(float64), result, tt.expected.(float64))
				}
			case string:
				result := Min(v, tt.b.(string))
				if result != tt.expected.(string) {
					t.Errorf("Min(%s, %s) = %s; want %s", v, tt.b.(string), result, tt.expected.(string))
				}
			}
		})
	}
}
