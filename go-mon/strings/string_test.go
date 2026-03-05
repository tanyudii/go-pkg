package strings

import (
	"reflect"
	"testing"
)

func TestSplitString(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		sep      string
		expected []string
	}{
		{
			name:     "comma separated values",
			str:      "a,b,c",
			sep:      ",",
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "space separated values",
			str:      "hello world foo",
			sep:      " ",
			expected: []string{"hello", "world", "foo"},
		},
		{
			name:     "pipe separated values",
			str:      "one|two|three",
			sep:      "|",
			expected: []string{"one", "two", "three"},
		},
		{
			name:     "single value",
			str:      "single",
			sep:      ",",
			expected: []string{"single"},
		},
		{
			name:     "empty separator",
			str:      "abc",
			sep:      "",
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "separator not in string",
			str:      "noseparator",
			sep:      ",",
			expected: []string{"noseparator"},
		},
		{
			name:     "consecutive separators",
			str:      "a,,b",
			sep:      ",",
			expected: []string{"a", "", "b"},
		},
		{
			name:     "separator at start",
			str:      ",a,b",
			sep:      ",",
			expected: []string{"", "a", "b"},
		},
		{
			name:     "separator at end",
			str:      "a,b,",
			sep:      ",",
			expected: []string{"a", "b", ""},
		},
		{
			name:     "empty string",
			str:      "",
			sep:      ",",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SplitString(tt.str, tt.sep)

			if tt.expected == nil {
				if result != nil {
					t.Errorf("SplitString(%q, %q) = %v; want nil", tt.str, tt.sep, result)
				}
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("SplitString(%q, %q) = %v; want %v", tt.str, tt.sep, result, tt.expected)
			}
		})
	}
}

func TestSplitStringToMapBool(t *testing.T) {
	tests := []struct {
		name     string
		val      string
		sep      string
		expected map[string]bool
	}{
		{
			name:     "comma separated values",
			val:      "a,b,c",
			sep:      ",",
			expected: map[string]bool{"a": true, "b": true, "c": true},
		},
		{
			name:     "space separated values",
			val:      "foo bar baz",
			sep:      " ",
			expected: map[string]bool{"foo": true, "bar": true, "baz": true},
		},
		{
			name:     "single value",
			val:      "single",
			sep:      ",",
			expected: map[string]bool{"single": true},
		},
		{
			name:     "empty string",
			val:      "",
			sep:      ",",
			expected: nil,
		},
		{
			name:     "consecutive separators creates empty keys",
			val:      "a,,b",
			sep:      ",",
			expected: map[string]bool{"a": true, "": true, "b": true},
		},
		{
			name:     "duplicate values",
			val:      "a,b,a,c",
			sep:      ",",
			expected: map[string]bool{"a": true, "b": true, "c": true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SplitStringToMapBool(tt.val, tt.sep)

			if tt.expected == nil {
				if result != nil {
					t.Errorf("SplitStringToMapBool(%q, %q) = %v; want nil", tt.val, tt.sep, result)
				}
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("SplitStringToMapBool(%q, %q) = %v; want %v", tt.val, tt.sep, result, tt.expected)
			}
		})
	}
}

func TestInSeparatedString(t *testing.T) {
	tests := []struct {
		name     string
		val      string
		check    string
		sep      string
		expected bool
	}{
		{
			name:     "value exists",
			val:      "apple,banana,cherry",
			check:    "banana",
			sep:      ",",
			expected: true,
		},
		{
			name:     "value does not exist",
			val:      "apple,banana,cherry",
			check:    "orange",
			sep:      ",",
			expected: false,
		},
		{
			name:     "empty string to check",
			val:      "a,b,c",
			check:    "",
			sep:      ",",
			expected: false,
		},
		{
			name:     "check empty value when present",
			val:      "a,,c",
			check:    "",
			sep:      ",",
			expected: true,
		},
		{
			name:     "empty source string",
			val:      "",
			check:    "anything",
			sep:      ",",
			expected: false,
		},
		{
			name:     "space separated",
			val:      "foo bar baz",
			check:    "bar",
			sep:      " ",
			expected: true,
		},
		{
			name:     "single value matches",
			val:      "only",
			check:    "only",
			sep:      ",",
			expected: true,
		},
		{
			name:     "single value does not match",
			val:      "only",
			check:    "other",
			sep:      ",",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := InSeparatedString(tt.val, tt.check, tt.sep)
			if result != tt.expected {
				t.Errorf("InSeparatedString(%q, %q, %q) = %v; want %v", tt.val, tt.check, tt.sep, result, tt.expected)
			}
		})
	}
}

func TestAppendSeparatedString(t *testing.T) {
	tests := []struct {
		name     string
		val      string
		sep      string
		values   []string
		expected string
	}{
		{
			name:     "append to empty string",
			val:      "",
			sep:      ",",
			values:   []string{"a", "b"},
			expected: "a,b",
		},
		{
			name:     "append to existing string",
			val:      "x,y",
			sep:      ",",
			values:   []string{"a", "b"},
			expected: "x,y,a,b",
		},
		{
			name:     "append single value",
			val:      "existing",
			sep:      ",",
			values:   []string{"new"},
			expected: "existing,new",
		},
		{
			name:     "append with space separator",
			val:      "hello world",
			sep:      " ",
			values:   []string{"foo", "bar"},
			expected: "hello world foo bar",
		},
		{
			name:     "append with pipe separator",
			val:      "a|b",
			sep:      "|",
			values:   []string{"c"},
			expected: "a|b|c",
		},
		{
			name:     "append empty slice",
			val:      "original",
			sep:      ",",
			values:   []string{},
			expected: "original,",
		},
		{
			name:     "append empty strings",
			val:      "start",
			sep:      ",",
			values:   []string{"", ""},
			expected: "start,,",
		},
		{
			name:     "append multiple values to empty",
			val:      "",
			sep:      "-",
			values:   []string{"one", "two", "three"},
			expected: "one-two-three",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AppendSeparatedString(tt.val, tt.sep, tt.values...)
			if result != tt.expected {
				t.Errorf("AppendSeparatedString(%q, %q, %v) = %q; want %q", tt.val, tt.sep, tt.values, result, tt.expected)
			}
		})
	}
}
