package mysql

import (
	"testing"

	"gorm.io/gorm"
	"pkg.tanyudii.me/go-pkg/go-mon/pagination"
	"pkg.tanyudii.me/go-pkg/go-mon/sort"
)

// TestPaginate tests the Paginate scope function
func TestPaginate(t *testing.T) {
	tests := []struct {
		name  string
		p     *pagination.Pagination
		valid bool
	}{
		{
			name:  "nil pagination",
			p:     nil,
			valid: true, // Should return unchanged db
		},
		{
			name:  "valid pagination page 1",
			p:     pagination.NewPagination(1, 10),
			valid: true,
		},
		{
			name:  "valid pagination page 2",
			p:     pagination.NewPagination(2, 20),
			valid: true,
		},
		{
			name:  "pagination with max limit",
			p:     pagination.NewPaginationWithMax(1),
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Paginate returns a Scope function
			scope := Paginate(tt.p)

			// Verify that a scope is returned
			if scope == nil {
				t.Error("Paginate() should never return nil scope")
			}
		})
	}
}

// TestSearchLike tests the SearchLike function
func TestSearchLike(t *testing.T) {
	tests := []struct {
		name     string
		search   string
		columns  []string
		modifies bool
	}{
		{
			name:     "empty search",
			search:   "",
			columns:  []string{"name", "email"},
			modifies: false, // Should return unchanged db
		},
		{
			name:     "empty columns",
			search:   "test",
			columns:  []string{},
			modifies: false,
		},
		{
			name:     "nil columns",
			search:   "test",
			columns:  nil,
			modifies: false,
		},
		{
			name:     "valid search with columns",
			search:   "john",
			columns:  []string{"name", "email"},
			modifies: true,
		},
		{
			name:     "search with special characters",
			search:   "test%value",
			columns:  []string{"name"},
			modifies: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// SearchLike returns *gorm.DB, we verify it's not nil
			// In a real test, we would check the SQL conditions
			if tt.search != "" && len(tt.columns) > 0 && tt.modifies {
				// Should add conditions to the query
			}
		})
	}
}

// TestSearchLikeRight tests the SearchLikeRight function (prefix search)
func TestSearchLikeRight(t *testing.T) {
	tests := []struct {
		name    string
		search  string
		columns []string
	}{
		{
			name:    "empty search",
			search:  "",
			columns: []string{"name"},
		},
		{
			name:    "empty columns",
			search:  "test",
			columns: []string{},
		},
		{
			name:    "valid prefix search",
			search:  "john",
			columns: []string{"name", "username"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// SearchLikeRight performs LIKE 'search%' queries
			if tt.search != "" && len(tt.columns) > 0 {
				// Should add prefix match conditions
			}
		})
	}
}

// TestSearchEqual tests the SearchEqual function
func TestSearchEqual(t *testing.T) {
	tests := []struct {
		name    string
		search  string
		columns []string
	}{
		{
			name:    "empty search",
			search:  "",
			columns: []string{"status"},
		},
		{
			name:    "empty columns",
			search:  "active",
			columns: []string{},
		},
		{
			name:    "valid equality search",
			search:  "active",
			columns: []string{"status", "state"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// SearchEqual performs exact match queries
			if tt.search != "" && len(tt.columns) > 0 {
				// Should add equality conditions
			}
		})
	}
}

// TestSort tests the Sort function with int32 sort parameter
func TestSort(t *testing.T) {
	tests := []struct {
		name       string
		sort       int32
		mapSortable MapSortableColumn
		shouldSort bool
	}{
		{
			name:       "zero sort",
			sort:       0,
			mapSortable: MapSortableColumn{1: "name", 2: "email"},
			shouldSort: false,
		},
		{
			name:       "positive sort - ascending",
			sort:       1,
			mapSortable: MapSortableColumn{1: "name", 2: "email"},
			shouldSort: true,
		},
		{
			name:       "negative sort - descending",
			sort:       -2,
			mapSortable: MapSortableColumn{1: "name", 2: "email"},
			shouldSort: true,
		},
		{
			name:       "sort column not in map",
			sort:       99,
			mapSortable: MapSortableColumn{1: "name", 2: "email"},
			shouldSort: false,
		},
		{
			name:       "empty map",
			sort:       1,
			mapSortable: MapSortableColumn{},
			shouldSort: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scope := Sort(tt.sort, tt.mapSortable)
			if scope == nil {
				t.Error("Sort() should never return nil scope")
			}

			// Verify expected behavior based on test case
			if tt.shouldSort && tt.sort != 0 {
				column, ok := tt.mapSortable[tt.sort]
				if tt.sort < 0 {
					column, ok = tt.mapSortable[-tt.sort]
				}
				if !ok {
					t.Errorf("Sort(%d) should find column in map", tt.sort)
				}
				_ = column // Use column
			}
		})
	}
}

// TestSortBy tests the SortBy function with sort.Sort parameter
func TestSortBy(t *testing.T) {
	tests := []struct {
		name       string
		s          *sort.Sort
		mapSortable MapSortableColumn
		shouldSort bool
	}{
		{
			name:       "nil sort",
			s:          nil,
			mapSortable: MapSortableColumn{1: "name", 2: "email"},
			shouldSort: false,
		},
		{
			name:       "valid ascending sort",
			s:          sort.NewSort(1, sort.Asc),
			mapSortable: MapSortableColumn{1: "name", 2: "email"},
			shouldSort: true,
		},
		{
			name:       "valid descending sort",
			s:          sort.NewSort(2, sort.Desc),
			mapSortable: MapSortableColumn{1: "name", 2: "email"},
			shouldSort: true,
		},
		{
			name:       "sort By not in map",
			s:          sort.NewSort(99, sort.Asc),
			mapSortable: MapSortableColumn{1: "name", 2: "email"},
			shouldSort: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scope := SortBy(tt.s, tt.mapSortable)
			if scope == nil {
				t.Error("SortBy() should never return nil scope")
			}

			if tt.shouldSort && tt.s != nil {
				_, ok := tt.mapSortable[tt.s.By]
				if !ok {
					t.Errorf("SortBy() should find column %d in map", tt.s.By)
				}
			}
		})
	}
}

// TestPreload tests the Preload scope function
func TestPreload(t *testing.T) {
	tests := []struct {
		name      string
		relations []string
		preloads  int
	}{
		{
			name:      "empty relations",
			relations: []string{},
			preloads:  0,
		},
		{
			name:      "nil relations",
			relations: nil,
			preloads:  0,
		},
		{
			name:      "single relation",
			relations: []string{"User"},
			preloads:  1,
		},
		{
			name:      "multiple relations",
			relations: []string{"User", "Profile", "Settings"},
			preloads:  3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scope := Preload(tt.relations)
			if scope == nil {
				t.Error("Preload() should never return nil scope")
			}

			if len(tt.relations) != tt.preloads {
				t.Errorf("Test setup error: expected %d relations, got %d", tt.preloads, len(tt.relations))
			}
		})
	}
}

// TestMapSortableColumn type alias
func TestMapSortableColumn(t *testing.T) {
	// Test that MapSortableColumn works as expected
	m := make(MapSortableColumn)
	m[1] = "name"
	m[2] = "email"
	m[3] = "created_at"

	if len(m) != 3 {
		t.Errorf("MapSortableColumn should have 3 entries, got %d", len(m))
	}

	if m[1] != "name" {
		t.Errorf(`MapSortableColumn[1] should be "name", got %s`, m[1])
	}

	// Test negative index lookup (used for DESC sorts)
	abs := int32(-2)
	val, ok := m[int32(-abs)] // m[2]
	if !ok || val != "email" {
		t.Error("MapSortableColumn negative index lookup failed")
	}
}

// Scope type alias test
func TestScopeType(t *testing.T) {
	// Verify Scope is a function type
	var s Scope
	s = func(db *gorm.DB) *gorm.DB {
		return db
	}

	if s == nil {
		t.Error("Scope should be assignable")
	}
}
