package sort

import (
	"testing"
)

func TestSortConstants(t *testing.T) {
	if Asc != 1 {
		t.Errorf("Asc = %d; want 1", Asc)
	}

	if Desc != 2 {
		t.Errorf("Desc = %d; want 2", Desc)
	}
}

func TestMapAscending(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"asc", true},
		{"ascending", true},
		{"ASC", false}, // Map is case-sensitive
		{"ASCENDING", false},
		{"Asc", false},
		{"Ascending", false},
		{"desc", false},
		{"descending", false},
		{"foo", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := MapAscending[tt.input]
			if result != tt.expected {
				t.Errorf("MapAscending[%s] = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestMapDescendingRegression tests the bug fix where MapDescending
// was incorrectly using the same values as MapAscending ("asc"/"ascending")
// instead of the correct values ("desc"/"descending")
func TestMapDescendingRegression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"desc", true},
		{"descending", true},
		{"DESC", false}, // Map is case-sensitive
		{"DESCENDING", false},
		{"Desc", false},
		{"Descending", false},
		{"asc", false}, // This was the bug - these were incorrectly true
		{"ascending", false}, // This was the bug - these were incorrectly true
		{"foo", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := MapDescending[tt.input]
			if result != tt.expected {
				t.Errorf("MapDescending[%s] = %v; want %v (REGRESSION TEST)", tt.input, result, tt.expected)
			}
		})
	}
}

func TestNewSort(t *testing.T) {
	tests := []struct {
		name     string
		by       int32
		typ      int32
		expected *Sort
	}{
		{
			name:     "normal values",
			by:       5,
			typ:      Asc,
			expected: &Sort{By: 5, Type: Asc},
		},
		{
			name:     "negative by",
			by:       -5,
			typ:      Desc,
			expected: &Sort{By: 0, Type: Desc}, // Only By is validated, Type stays
		},
		{
			name:     "negative type",
			by:       5,
			typ:      -1,
			expected: &Sort{By: 5, Type: 0}, // Validate sets negative to 0
		},
		{
			name:     "both negative",
			by:       -5,
			typ:      -1,
			expected: &Sort{By: 0, Type: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewSort(tt.by, tt.typ)
			if result.By != tt.expected.By {
				t.Errorf("NewSort().By = %d; want %d", result.By, tt.expected.By)
			}
			if result.Type != tt.expected.Type {
				t.Errorf("NewSort().Type = %d; want %d", result.Type, tt.expected.Type)
			}
		})
	}
}

func TestNewSortFromString(t *testing.T) {
	tests := []struct {
		name         string
		by           int32
		typeString   string
		expectedType int32
	}{
		{
			name:         "ascending string",
			by:           5,
			typeString:   "asc",
			expectedType: 0, // String based, Type stays 0
		},
		{
			name:         "descending string",
			by:           5,
			typeString:   "desc",
			expectedType: 0,
		},
		{
			name:         "negative by",
			by:           -5,
			typeString:   "asc",
			expectedType: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewSortFromString(tt.by, tt.typeString)
			if result.By != tt.by && tt.by > 0 {
				t.Errorf("NewSortFromString().By = %d; want %d", result.By, tt.by)
			}
			if result.By < 0 {
				t.Errorf("NewSortFromString().By = %d; should be >= 0 after Validate", result.By)
			}
			if result.TypeString != tt.typeString {
				t.Errorf("NewSortFromString().TypeString = %s; want %s", result.TypeString, tt.typeString)
			}
		})
	}
}

func TestSort_IsAsc(t *testing.T) {
	tests := []struct {
		name     string
		sort     *Sort
		expected bool
	}{
		{
			name:     "type equals Asc",
			sort:     &Sort{Type: Asc},
			expected: true,
		},
		{
			name:     "typeString is asc",
			sort:     &Sort{TypeString: "asc"},
			expected: true,
		},
		{
			name:     "typeString is ascending",
			sort:     &Sort{TypeString: "ascending"},
			expected: true,
		},
		{
			name:     "type equals Desc",
			sort:     &Sort{Type: Desc},
			expected: false,
		},
		{
			name:     "typeString is desc",
			sort:     &Sort{TypeString: "desc"},
			expected: false,
		},
		{
			name:     "both match ascending",
			sort:     &Sort{Type: Asc, TypeString: "ascending"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.sort.IsAsc()
			if result != tt.expected {
				t.Errorf("Sort.IsAsc() = %v; want %v", result, tt.expected)
			}
		})
	}
}

func TestSort_IsDesc(t *testing.T) {
	tests := []struct {
		name     string
		sort     *Sort
		expected bool
	}{
		{
			name:     "type equals Desc",
			sort:     &Sort{Type: Desc},
			expected: true,
		},
		{
			name:     "typeString is desc",
			sort:     &Sort{TypeString: "desc"},
			expected: true,
		},
		{
			name:     "typeString is descending",
			sort:     &Sort{TypeString: "descending"},
			expected: true,
		},
		{
			name:     "type equals Asc",
			sort:     &Sort{Type: Asc},
			expected: false,
		},
		{
			name:     "typeString is asc",
			sort:     &Sort{TypeString: "asc"},
			expected: false,
		},
		{
			name:     "both match descending",
			sort:     &Sort{Type: Desc, TypeString: "descending"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.sort.IsDesc()
			if result != tt.expected {
				t.Errorf("Sort.IsDesc() = %v; want %v", result, tt.expected)
			}
		})
	}
}

func TestSort_Validate(t *testing.T) {
	tests := []struct {
		name          string
		sort          *Sort
		expectedBy    int32
		expectedType  int32
	}{
		{
			name:         "all positive",
			sort:         &Sort{By: 5, Type: 1},
			expectedBy:   5,
			expectedType: 1,
		},
		{
			name:         "negative By",
			sort:         &Sort{By: -5, Type: 1},
			expectedBy:   0,
			expectedType: 1,
		},
		{
			name:         "negative Type",
			sort:         &Sort{By: 5, Type: -1},
			expectedBy:   5,
			expectedType: 0,
		},
		{
			name:         "both negative",
			sort:         &Sort{By: -5, Type: -1},
			expectedBy:   0,
			expectedType: 0,
		},
		{
			name:         "zeros",
			sort:         &Sort{By: 0, Type: 0},
			expectedBy:   0,
			expectedType: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.sort.Validate()
			if tt.sort.By != tt.expectedBy {
				t.Errorf("Validate().By = %d; want %d", tt.sort.By, tt.expectedBy)
			}
			if tt.sort.Type != tt.expectedType {
				t.Errorf("Validate().Type = %d; want %d", tt.sort.Type, tt.expectedType)
			}
		})
	}
}
