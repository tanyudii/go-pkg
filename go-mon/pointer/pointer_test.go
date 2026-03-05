package pointer

import (
	"testing"
)

func TestVal(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{"int", 42, 42},
		{"string", "hello", "hello"},
		{"float64", 3.14, 3.14},
		{"bool", true, true},
		{"zero int", 0, 0},
		{"zero float", 0.0, 0.0},
		{"empty string", "", ""},
		{"false", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.input.(type) {
			case int:
				result := Val(v)
				if result == nil {
					t.Fatal("Val() returned nil")
				}
				if *result != tt.expected.(int) {
					t.Errorf("Val(%d) = %d; want %d", v, *result, tt.expected.(int))
				}
			case string:
				result := Val(v)
				if result == nil {
					t.Fatal("Val() returned nil")
				}
				if *result != tt.expected.(string) {
					t.Errorf("Val(%s) = %s; want %s", v, *result, tt.expected.(string))
				}
			case float64:
				result := Val(v)
				if result == nil {
					t.Fatal("Val() returned nil")
				}
				if *result != tt.expected.(float64) {
					t.Errorf("Val(%f) = %f; want %f", v, *result, tt.expected.(float64))
				}
			case bool:
				result := Val(v)
				if result == nil {
					t.Fatal("Val() returned nil")
				}
				if *result != tt.expected.(bool) {
					t.Errorf("Val(%t) = %t; want %t", v, *result, tt.expected.(bool))
				}
			}
		})
	}
}

func TestExtract(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{"valid int pointer", func() *int { v := 42; return &v }(), 42},
		{"valid string pointer", func() *string { v := "hello"; return &v }(), "hello"},
		{"nil int pointer", (*int)(nil), 0},
		{"nil string pointer", (*string)(nil), ""},
		{"nil bool pointer", (*bool)(nil), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.input.(type) {
			case *int:
				result := Extract(v)
				if result != tt.expected.(int) {
					t.Errorf("Extract(%v) = %d; want %d", v, result, tt.expected.(int))
				}
			case *string:
				result := Extract(v)
				if result != tt.expected.(string) {
					t.Errorf("Extract(%v) = %s; want %s", v, result, tt.expected.(string))
				}
			case *bool:
				result := Extract(v)
				if result != tt.expected.(bool) {
					t.Errorf("Extract(%v) = %t; want %t", v, result, tt.expected.(bool))
				}
			}
		})
	}
}

func TestEmptyNil(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected bool
	}{
		{"non-zero int", 42, false},
		{"zero int", 0, true},
		{"non-empty string", "hello", false},
		{"empty string", "", true},
		{"true", true, false},
		{"false", false, true},
		{"non-zero float", 3.14, false},
		{"zero float", 0.0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.input.(type) {
			case int:
				result := EmptyNil(v)
				isNil := result == nil
				if isNil != tt.expected {
					t.Errorf("EmptyNil(%d) = %v; isNil=%v; want isNil=%v", v, result, isNil, tt.expected)
				}
			case string:
				result := EmptyNil(v)
				isNil := result == nil
				if isNil != tt.expected {
					t.Errorf("EmptyNil(%s) = %v; isNil=%v; want isNil=%v", v, result, isNil, tt.expected)
				}
			case bool:
				result := EmptyNil(v)
				isNil := result == nil
				if isNil != tt.expected {
					t.Errorf("EmptyNil(%t) = %v; isNil=%v; want isNil=%v", v, result, isNil, tt.expected)
				}
			case float64:
				result := EmptyNil(v)
				isNil := result == nil
				if isNil != tt.expected {
					t.Errorf("EmptyNil(%f) = %v; isNil=%v; want isNil=%v", v, result, isNil, tt.expected)
				}
			}
		})
	}
}
