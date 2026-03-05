package pagination

import (
	"testing"
)

func TestPaginationConstants(t *testing.T) {
	if DefaultPage != 1 {
		t.Errorf("DefaultPage = %d; want 1", DefaultPage)
	}

	if DefaultLimit != 10 {
		t.Errorf("DefaultLimit = %d; want 10", DefaultLimit)
	}

	if MaximumLimit != 100 {
		t.Errorf("MaximumLimit = %d; want 100", MaximumLimit)
	}
}

func TestNewPagination(t *testing.T) {
	tests := []struct {
		name  string
		page  int32
		limit int32
	}{
		{"normal values", 1, 10},
		{"second page", 2, 20},
		{"large page", 100, 50},
		{"zero page", 0, 10},
		{"negative page", -1, 10},
		{"zero limit", 1, 0},
		{"negative limit", 1, -1},
		{"exceeds maximum", 1, 150},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPagination(tt.page, tt.limit)

			if p.Page <= 0 {
				t.Errorf("NewPagination(%d, %d).Page = %d; want > 0", tt.page, tt.limit, p.Page)
			}

			if p.Limit <= 0 || p.Limit > MaximumLimit {
				t.Errorf("NewPagination(%d, %d).Limit = %d; want (0, 100]", tt.page, tt.limit, p.Limit)
			}
		})
	}
}

func TestGetOffset(t *testing.T) {
	tests := []struct {
		name     string
		page     int32
		limit    int32
		expected int32
	}{
		{"first page", 1, 10, 0},
		{"second page", 2, 10, 10},
		{"third page", 3, 20, 40},
		{"large page", 10, 100, 900},
		{"zero page validated", 0, 10, 0}, // Will be validated to DefaultPage
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPagination(tt.page, tt.limit)
			offset := p.GetOffset()

			if offset != tt.expected {
				t.Errorf("GetOffset() = %d; want %d", offset, tt.expected)
			}
		})
	}
}

func TestSetPagination(t *testing.T) {
	tests := []struct {
		name          string
		page          int32
		limit         int32
		total         int64
		expectedTotal int64
		expectedMax   int32
	}{
		{
			name:          "exact division",
			page:          1,
			limit:         10,
			total:         100,
			expectedTotal: 100,
			expectedMax:   10,
		},
		{
			name:          "with remainder",
			page:          1,
			limit:         10,
			total:         105,
			expectedTotal: 105,
			expectedMax:   11,
		},
		{
			name:          "more than limit",
			page:          1,
			limit:         10,
			total:         25,
			expectedTotal: 25,
			expectedMax:   3,
		},
		{
			name:          "zero total",
			page:          1,
			limit:         10,
			total:         0,
			expectedTotal: 0,
			expectedMax:   0,
		},
		{
			name:          "less than limit",
			page:          1,
			limit:         100,
			total:         5,
			expectedTotal: 5,
			expectedMax:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPagination(tt.page, tt.limit)
			p.SetPagination(tt.total)

			if p.Total != tt.expectedTotal {
				t.Errorf("Total = %d; want %d", p.Total, tt.expectedTotal)
			}

			if p.TotalPage != tt.expectedMax {
				t.Errorf("TotalPage = %d; want %d", p.TotalPage, tt.expectedMax)
			}
		})
	}
}

func TestHasNextPage(t *testing.T) {
	tests := []struct {
		name     string
		page     int32
		limit    int32
		total    int64
		expected bool
	}{
		{"has next", 1, 10, 100, true},
		{"last page", 10, 10, 100, false},
		{"zero total", 1, 10, 0, false},
		{"middle page", 5, 10, 100, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPagination(tt.page, tt.limit)
			p.SetPagination(tt.total)

			result := p.HasNextPage()
			if result != tt.expected {
				t.Errorf("HasNextPage() = %v; want %v", result, tt.expected)
			}
		})
	}
}

func TestIncr(t *testing.T) {
	p := NewPagination(1, 10)

	if p.Page != 1 {
		t.Errorf("Initial Page = %d; want 1", p.Page)
	}

	p.Incr()

	if p.Page != 2 {
		t.Errorf("After Incr() Page = %d; want 2", p.Page)
	}

	p.Incr()

	if p.Page != 3 {
		t.Errorf("After second Incr() Page = %d; want 3", p.Page)
	}
}

func TestNewPaginationWithMax(t *testing.T) {
	p := NewPaginationWithMax(5)

	if p.Page != 5 {
		t.Errorf("Page = %d; want 5", p.Page)
	}

	if p.Limit != MaximumLimit {
		t.Errorf("Limit = %d; want %d", p.Limit, MaximumLimit)
	}
}

func TestSetPaginationWithoutTotal(t *testing.T) {
	p := NewPagination(1, 10)

	p.SetPagination() // No total provided

	if p.Total != 0 {
		t.Errorf("Total = %d; want 0", p.Total)
	}

	if p.TotalPage != 0 {
		t.Errorf("TotalPage = %d; want 0", p.TotalPage)
	}
}
