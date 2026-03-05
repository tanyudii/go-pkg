package slice

import (
	"testing"
)

func TestToMap(t *testing.T) {
	type Person struct {
		ID   int
		Name string
	}

	people := []Person{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
		{ID: 3, Name: "Charlie"},
	}

	result := ToMap(people, func(p Person) int {
		return p.ID
	})

	if len(result) != 3 {
		t.Errorf("Expected map length 3, got %d", len(result))
	}

	if result[1].Name != "Alice" {
		t.Errorf("Expected Alice, got %s", result[1].Name)
	}

	if result[2].Name != "Bob" {
		t.Errorf("Expected Bob, got %s", result[2].Name)
	}
}

func TestToMapBool(t *testing.T) {
	input := []string{"apple", "banana", "cherry"}
	result := ToMapBool(input)

	if len(result) != 3 {
		t.Errorf("Expected map length 3, got %d", len(result))
	}

	if !result["apple"] {
		t.Errorf("Expected apple to be true")
	}

	if !result["banana"] {
		t.Errorf("Expected banana to be true")
	}

	if result["orange"] {
		t.Errorf("Expected orange to be false (not in map)")
	}
}

func TestToSlice(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	result := ToSlice(input, func(i int) string {
		return string(rune('a' + i - 1))
	})

	expected := []string{"a", "b", "c", "d", "e"}
	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
	}

	for i := range expected {
		if result[i] != expected[i] {
			t.Errorf("Index %d: expected %s, got %s", i, expected[i], result[i])
		}
	}
}

func TestToSliceNil(t *testing.T) {
	var input []int
	result := ToSlice(input, func(i int) int {
		return i * 2
	})

	if result != nil {
		t.Errorf("Expected nil for nil input, got %v", result)
	}
}

func TestContains(t *testing.T) {
	t.Run("int slice", func(t *testing.T) {
		tests := []struct {
			name     string
			slice    []int
			elem     int
			expected bool
		}{
			{"element exists", []int{1, 2, 3, 4, 5}, 3, true},
			{"element does not exist", []int{1, 2, 3, 4, 5}, 6, false},
			{"empty slice", []int{}, 1, false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := Contains(tt.slice, tt.elem)
				if result != tt.expected {
					t.Errorf("Contains(%v, %v) = %v; want %v", tt.slice, tt.elem, result, tt.expected)
				}
			})
		}
	})

	t.Run("string slice", func(t *testing.T) {
		if !Contains([]string{"a", "b", "c"}, "b") {
			t.Error("Expected to find 'b' in string slice")
		}
		if Contains([]string{"a", "b", "c"}, "d") {
			t.Error("Expected not to find 'd' in string slice")
		}
	})
}

func TestIndex(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		elem     int
		expected int
	}{
		{"element exists", []int{1, 2, 3, 4, 5}, 3, 2},
		{"element does not exist", []int{1, 2, 3, 4, 5}, 6, -1},
		{"empty slice", []int{}, 1, -1},
		{"first element", []int{1, 2, 3}, 1, 0},
		{"last element", []int{1, 2, 3}, 3, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Index(tt.slice, tt.elem)
			if result != tt.expected {
				t.Errorf("Index(%v, %v) = %d; want %d", tt.slice, tt.elem, result, tt.expected)
			}
		})
	}
}

func TestReverse(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{"normal slice", []int{1, 2, 3, 4, 5}, []int{5, 4, 3, 2, 1}},
		{"empty slice", []int{}, []int{}},
		{"single element", []int{42}, []int{42}},
		{"two elements", []int{1, 2}, []int{2, 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make a copy to avoid modifying original
			input := make([]int, len(tt.input))
			copy(input, tt.input)

			Reverse(input)

			if len(input) != len(tt.expected) {
				t.Errorf("Reverse() length = %d; want %d", len(input), len(tt.expected))
			}

			for i := range input {
				if input[i] != tt.expected[i] {
					t.Errorf("Reverse()[%d] = %d; want %d", i, input[i], tt.expected[i])
				}
			}
		})
	}
}

func TestMax(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		expected int
	}{
		{"normal slice", []int{1, 5, 3, 9, 2}, 9},
		{"empty slice", []int{}, 0},
		{"single element", []int{42}, 42},
		{"negative numbers", []int{-5, -2, -10}, -2},
		{"all same", []int{5, 5, 5}, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Max(tt.slice)
			if result != tt.expected {
				t.Errorf("Max(%v) = %d; want %d", tt.slice, result, tt.expected)
			}
		})
	}
}

func TestMin(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		expected int
	}{
		{"normal slice", []int{1, 5, 3, 9, 2}, 1},
		{"empty slice", []int{}, 0},
		{"single element", []int{42}, 42},
		{"negative numbers", []int{-5, -2, -10}, -10},
		{"all same", []int{5, 5, 5}, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Min(tt.slice)
			if result != tt.expected {
				t.Errorf("Min(%v) = %d; want %d", tt.slice, result, tt.expected)
			}
		})
	}
}

func TestSort(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{"unsorted", []int{3, 1, 4, 1, 5, 9, 2, 6}, []int{1, 1, 2, 3, 4, 5, 6, 9}},
		{"already sorted", []int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
		{"reverse sorted", []int{5, 4, 3, 2, 1}, []int{1, 2, 3, 4, 5}},
		{"empty", []int{}, []int{}},
		{"single", []int{42}, []int{42}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := make([]int, len(tt.input))
			copy(input, tt.input)

			Sort(input)

			if len(input) != len(tt.expected) {
				t.Errorf("Sort() length = %d; want %d", len(input), len(tt.expected))
			}

			for i := range input {
				if input[i] != tt.expected[i] {
					t.Errorf("Sort()[%d] = %d; want %d", i, input[i], tt.expected[i])
				}
			}
		})
	}
}

func TestSortFunc(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	people := []Person{
		{"Alice", 30},
		{"Bob", 25},
		{"Charlie", 35},
	}

	SortFunc(people, func(a, b Person) int {
		if a.Age < b.Age {
			return -1
		}
		if a.Age > b.Age {
			return 1
		}
		return 0
	})

	expected := []string{"Bob", "Alice", "Charlie"}
	for i, p := range people {
		if p.Name != expected[i] {
			t.Errorf("Index %d: expected %s, got %s", i, expected[i], p.Name)
		}
	}
}

func TestSortBy(t *testing.T) {
	input := []int{3, 1, 4, 1, 5, 9, 2, 6}

	SortBy(input, func(a, b int) bool {
		return a > b // Sort descending
	})

	expected := []int{9, 6, 5, 4, 3, 2, 1, 1}
	for i := range input {
		if input[i] != expected[i] {
			t.Errorf("SortBy()[%d] = %d; want %d", i, input[i], expected[i])
		}
	}
}

func TestClone(t *testing.T) {
	original := []int{1, 2, 3, 4, 5}
	cloned := Clone(original)

	if len(cloned) != len(original) {
		t.Errorf("Clone() length = %d; want %d", len(cloned), len(original))
	}

	for i := range original {
		if cloned[i] != original[i] {
			t.Errorf("Clone()[%d] = %d; want %d", i, cloned[i], original[i])
		}
	}

	// Modify clone and check original is unchanged
	cloned[0] = 999
	if original[0] == 999 {
		t.Error("Clone() did not create a separate copy")
	}
}

func TestEqual(t *testing.T) {
	tests := []struct {
		name     string
		a, b     []int
		expected bool
	}{
		{"equal slices", []int{1, 2, 3}, []int{1, 2, 3}, true},
		{"different length", []int{1, 2}, []int{1, 2, 3}, false},
		{"different values", []int{1, 2, 3}, []int{1, 2, 4}, false},
		{"both empty", []int{}, []int{}, true},
		{"both nil", nil, nil, true},
		{"one nil one empty", nil, []int{}, true}, // slices.Equal treats nil and [] as equal
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Equal(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Equal(%v, %v) = %v; want %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestChunk(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		size     int
		expected [][]int
	}{
		{
			name:     "even chunks",
			slice:    []int{1, 2, 3, 4, 5, 6},
			size:     2,
			expected: [][]int{{1, 2}, {3, 4}, {5, 6}},
		},
		{
			name:     "uneven chunks",
			slice:    []int{1, 2, 3, 4, 5},
			size:     2,
			expected: [][]int{{1, 2}, {3, 4}, {5}},
		},
		{
			name:     "chunk size larger than slice",
			slice:    []int{1, 2, 3},
			size:     5,
			expected: [][]int{{1, 2, 3}},
		},
		{
			name:     "empty slice",
			slice:    []int{},
			size:     2,
			expected: nil,
		},
		{
			name:     "invalid size",
			slice:    []int{1, 2, 3},
			size:     0,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Chunk(tt.slice, tt.size)

			if len(result) != len(tt.expected) {
				t.Fatalf("Chunk() length = %d; want %d", len(result), len(tt.expected))
			}

			for i := range result {
				if !Equal(result[i], tt.expected[i]) {
					t.Errorf("Chunk()[%d] = %v; want %v", i, result[i], tt.expected[i])
				}
			}
		})
	}
}

func TestFilter(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		predicate func(int) bool
		expected []int
	}{
		{
			name:  "even numbers",
			slice: []int{1, 2, 3, 4, 5, 6},
			predicate: func(n int) bool {
				return n%2 == 0
			},
			expected: []int{2, 4, 6},
		},
		{
			name:  "greater than 3",
			slice: []int{1, 2, 3, 4, 5},
			predicate: func(n int) bool {
				return n > 3
			},
			expected: []int{4, 5},
		},
		{
			name:     "none match",
			slice:    []int{1, 2, 3},
			predicate: func(n int) bool { return n > 10 },
			expected: []int{},
		},
		{
			name:     "all match",
			slice:    []int{1, 2, 3},
			predicate: func(n int) bool { return n > 0 },
			expected: []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Filter(tt.slice, tt.predicate)

			if !Equal(result, tt.expected) {
				t.Errorf("Filter() = %v; want %v", result, tt.expected)
			}
		})
	}
}

func TestUnique(t *testing.T) {
	t.Run("int slice", func(t *testing.T) {
		tests := []struct {
			name     string
			input    []int
			expected []int
		}{
			{"with duplicates", []int{1, 2, 2, 3, 3, 3, 4}, []int{1, 2, 3, 4}},
			{"no duplicates", []int{1, 2, 3, 4}, []int{1, 2, 3, 4}},
			{"empty", []int{}, []int{}},
			{"all same", []int{5, 5, 5, 5}, []int{5}},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := Unique(tt.input)
				if !Equal(result, tt.expected) {
					t.Errorf("Unique() = %v; want %v", result, tt.expected)
				}
			})
		}
	})

	t.Run("string slice", func(t *testing.T) {
		result := Unique([]string{"a", "b", "a", "c"})
		expected := []string{"a", "b", "c"}
		if !Equal(result, expected) {
			t.Errorf("Unique() = %v; want %v", result, expected)
		}
	})
}

func TestFlatten(t *testing.T) {
	tests := []struct {
		name     string
		input    [][]int
		expected []int
	}{
		{
			name:     "normal 2D slice",
			input:    [][]int{{1, 2}, {3, 4, 5}, {6}},
			expected: []int{1, 2, 3, 4, 5, 6},
		},
		{
			name:     "empty outer slice",
			input:    [][]int{},
			expected: []int{},
		},
		{
			name:     "empty inner slices",
			input:    [][]int{{}, {}, {}},
			expected: []int{},
		},
		{
			name:     "mixed empty and non-empty",
			input:    [][]int{{1, 2}, {}, {3, 4}, {}},
			expected: []int{1, 2, 3, 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Flatten(tt.input)

			if !Equal(result, tt.expected) {
				t.Errorf("Flatten() = %v; want %v", result, tt.expected)
			}
		})
	}
}
