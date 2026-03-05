package actor

import (
	"context"
	"testing"

	gotex "pkg.tanyudii.me/go-pkg/go-tex"
)

func TestEntityHaveActor_SetCreator(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		wantCreatedBy *string
		wantUpdatedBy *string
	}{
		{
			name:          "nil context - no changes",
			ctx:           nil,
			wantCreatedBy: nil,
			wantUpdatedBy: nil,
		},
		{
			name:          "context without gotex - no changes",
			ctx:           context.Background(),
			wantCreatedBy: nil,
			wantUpdatedBy: nil,
		},
		{
			name: "context with gotex - sets both fields",
			ctx: gotex.NewContext(context.Background(), &gotex.Gotex{
				UserID: "user123",
			}),
			wantCreatedBy: func() *string { s := "user123"; return &s }(),
			wantUpdatedBy: func() *string { s := "user123"; return &s }(),
		},
		{
			name: "context with different user ID",
			ctx: gotex.NewContext(context.Background(), &gotex.Gotex{
				UserID: "admin456",
			}),
			wantCreatedBy: func() *string { s := "admin456"; return &s }(),
			wantUpdatedBy: func() *string { s := "admin456"; return &s }(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EntityHaveActor{}
			err := e.SetCreator(tt.ctx)

			if err != nil {
				t.Errorf("SetCreator() error = %v, want nil", err)
			}

			if !compareStringPtr(e.CreatedBy, tt.wantCreatedBy) {
				t.Errorf("SetCreator() CreatedBy = %v, want %v", e.CreatedBy, tt.wantCreatedBy)
			}

			if !compareStringPtr(e.UpdatedBy, tt.wantUpdatedBy) {
				t.Errorf("SetCreator() UpdatedBy = %v, want %v", e.UpdatedBy, tt.wantUpdatedBy)
			}
		})
	}
}

func TestEntityHaveActor_SetEditor(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		initialState  *EntityHaveActor
		wantUpdatedBy *string
	}{
		{
			name:          "nil context - no changes",
			ctx:           nil,
			initialState:  &EntityHaveActor{CreatedBy: func() *string { s := "creator"; return &s }()},
			wantUpdatedBy: nil,
		},
		{
			name:          "context without gotex - no changes",
			ctx:           context.Background(),
			initialState:  &EntityHaveActor{CreatedBy: func() *string { s := "creator"; return &s }()},
			wantUpdatedBy: nil,
		},
		{
			name: "context with gotex - sets updated by",
			ctx: gotex.NewContext(context.Background(), &gotex.Gotex{
				UserID: "editor123",
			}),
			initialState:  &EntityHaveActor{CreatedBy: func() *string { s := "creator"; return &s }()},
			wantUpdatedBy: func() *string { s := "editor123"; return &s }(),
		},
		{
			name: "overwrites existing updated by",
			ctx: gotex.NewContext(context.Background(), &gotex.Gotex{
				UserID: "newEditor456",
			}),
			initialState: &EntityHaveActor{
				CreatedBy:  func() *string { s := "creator"; return &s }(),
				UpdatedBy:  func() *string { s := "oldEditor"; return &s }(),
			},
			wantUpdatedBy: func() *string { s := "newEditor456"; return &s }(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := tt.initialState
			err := e.SetEditor(tt.ctx)

			if err != nil {
				t.Errorf("SetEditor() error = %v, want nil", err)
			}

			if !compareStringPtr(e.UpdatedBy, tt.wantUpdatedBy) {
				t.Errorf("SetEditor() UpdatedBy = %v, want %v", e.UpdatedBy, tt.wantUpdatedBy)
			}

			// Ensure CreatedBy is not changed
			if e.CreatedBy != tt.initialState.CreatedBy {
				t.Errorf("SetEditor() CreatedBy should not change, got %v", e.CreatedBy)
			}
		})
	}
}

func TestEntityHaveActor_SetDestroyer(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		initialState  *EntityHaveActor
		wantDeletedBy *string
	}{
		{
			name:          "nil context - no changes",
			ctx:           nil,
			initialState:  &EntityHaveActor{CreatedBy: func() *string { s := "creator"; return &s }()},
			wantDeletedBy: nil,
		},
		{
			name:          "context without gotex - no changes",
			ctx:           context.Background(),
			initialState:  &EntityHaveActor{CreatedBy: func() *string { s := "creator"; return &s }()},
			wantDeletedBy: nil,
		},
		{
			name: "context with gotex - sets deleted by",
			ctx: gotex.NewContext(context.Background(), &gotex.Gotex{
				UserID: "deleter123",
			}),
			initialState:  &EntityHaveActor{CreatedBy: func() *string { s := "creator"; return &s }()},
			wantDeletedBy: func() *string { s := "deleter123"; return &s }(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := tt.initialState
			err := e.SetDestroyer(tt.ctx)

			if err != nil {
				t.Errorf("SetDestroyer() error = %v, want nil", err)
			}

			if !compareStringPtr(e.DeletedBy, tt.wantDeletedBy) {
				t.Errorf("SetDestroyer() DeletedBy = %v, want %v", e.DeletedBy, tt.wantDeletedBy)
			}

			// Ensure CreatedBy and UpdatedBy are not changed
			if e.CreatedBy != tt.initialState.CreatedBy {
				t.Errorf("SetDestroyer() CreatedBy should not change, got %v", e.CreatedBy)
			}
			if e.UpdatedBy != tt.initialState.UpdatedBy {
				t.Errorf("SetDestroyer() UpdatedBy should not change, got %v", e.UpdatedBy)
			}
		})
	}
}

func TestEntityHaveCEditor_SetCreator(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		wantCreatedBy *string
		wantUpdatedBy *string
	}{
		{
			name:          "nil context - no changes",
			ctx:           nil,
			wantCreatedBy: nil,
			wantUpdatedBy: nil,
		},
		{
			name:          "context without gotex - no changes",
			ctx:           context.Background(),
			wantCreatedBy: nil,
			wantUpdatedBy: nil,
		},
		{
			name: "context with gotex - sets both fields",
			ctx: gotex.NewContext(context.Background(), &gotex.Gotex{
				UserID: "user789",
			}),
			wantCreatedBy: func() *string { s := "user789"; return &s }(),
			wantUpdatedBy: func() *string { s := "user789"; return &s }(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EntityHaveCEditor{}
			err := e.SetCreator(tt.ctx)

			if err != nil {
				t.Errorf("SetCreator() error = %v, want nil", err)
			}

			if !compareStringPtr(e.CreatedBy, tt.wantCreatedBy) {
				t.Errorf("SetCreator() CreatedBy = %v, want %v", e.CreatedBy, tt.wantCreatedBy)
			}

			if !compareStringPtr(e.UpdatedBy, tt.wantUpdatedBy) {
				t.Errorf("SetCreator() UpdatedBy = %v, want %v", e.UpdatedBy, tt.wantUpdatedBy)
			}
		})
	}
}

func TestEntityHaveCEditor_SetEditor(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		initialState  *EntityHaveCEditor
		wantUpdatedBy *string
	}{
		{
			name:          "nil context - no changes",
			ctx:           nil,
			initialState:  &EntityHaveCEditor{CreatedBy: func() *string { s := "creator"; return &s }()},
			wantUpdatedBy: nil,
		},
		{
			name:          "context without gotex - no changes",
			ctx:           context.Background(),
			initialState:  &EntityHaveCEditor{CreatedBy: func() *string { s := "creator"; return &s }()},
			wantUpdatedBy: nil,
		},
		{
			name: "context with gotex - sets updated by",
			ctx: gotex.NewContext(context.Background(), &gotex.Gotex{
				UserID: "editor789",
			}),
			initialState:  &EntityHaveCEditor{CreatedBy: func() *string { s := "creator"; return &s }()},
			wantUpdatedBy: func() *string { s := "editor789"; return &s }(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := tt.initialState
			err := e.SetEditor(tt.ctx)

			if err != nil {
				t.Errorf("SetEditor() error = %v, want nil", err)
			}

			if !compareStringPtr(e.UpdatedBy, tt.wantUpdatedBy) {
				t.Errorf("SetEditor() UpdatedBy = %v, want %v", e.UpdatedBy, tt.wantUpdatedBy)
			}

			// Ensure CreatedBy is not changed
			if e.CreatedBy != tt.initialState.CreatedBy {
				t.Errorf("SetEditor() CreatedBy should not change, got %v", e.CreatedBy)
			}
		})
	}
}

func TestEntityHaveCreator_SetCreator(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		wantCreatedBy *string
	}{
		{
			name:          "nil context - no changes",
			ctx:           nil,
			wantCreatedBy: nil,
		},
		{
			name:          "context without gotex - no changes",
			ctx:           context.Background(),
			wantCreatedBy: nil,
		},
		{
			name: "context with gotex - sets created by",
			ctx: gotex.NewContext(context.Background(), &gotex.Gotex{
				UserID: "creator999",
			}),
			wantCreatedBy: func() *string { s := "creator999"; return &s }(),
		},
		{
			name: "overwrites existing created by",
			ctx: gotex.NewContext(context.Background(), &gotex.Gotex{
				UserID: "newCreator888",
			}),
			wantCreatedBy: func() *string { s := "newCreator888"; return &s }(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EntityHaveCreator{}
			err := e.SetCreator(tt.ctx)

			if err != nil {
				t.Errorf("SetCreator() error = %v, want nil", err)
			}

			if !compareStringPtr(e.CreatedBy, tt.wantCreatedBy) {
				t.Errorf("SetCreator() CreatedBy = %v, want %v", e.CreatedBy, tt.wantCreatedBy)
			}
		})
	}
}

func TestEntityActorInterface(t *testing.T) {
	// Test that EntityHaveActor implements EntityActor interface
	var _ EntityActor = &EntityHaveActor{}

	e := &EntityHaveActor{}
	ctx := gotex.NewContext(context.Background(), &gotex.Gotex{UserID: "testuser"})

	// Test all interface methods
	if err := e.SetCreator(ctx); err != nil {
		t.Errorf("EntityActor.SetCreator() error = %v", err)
	}
	if e.CreatedBy == nil || *e.CreatedBy != "testuser" {
		t.Error("EntityActor.SetCreator() failed to set CreatedBy")
	}

	if err := e.SetEditor(ctx); err != nil {
		t.Errorf("EntityActor.SetEditor() error = %v", err)
	}
	if e.UpdatedBy == nil || *e.UpdatedBy != "testuser" {
		t.Error("EntityActor.SetEditor() failed to set UpdatedBy")
	}

	if err := e.SetDestroyer(ctx); err != nil {
		t.Errorf("EntityActor.SetDestroyer() error = %v", err)
	}
	if e.DeletedBy == nil || *e.DeletedBy != "testuser" {
		t.Error("EntityActor.SetDestroyer() failed to set DeletedBy")
	}
}

func TestEntityCEditorInterface(t *testing.T) {
	// Test that EntityHaveCEditor implements EntityCEditor interface
	var _ EntityCEditor = &EntityHaveCEditor{}

	e := &EntityHaveCEditor{}
	ctx := gotex.NewContext(context.Background(), &gotex.Gotex{UserID: "testuser2"})

	// Test all interface methods
	if err := e.SetCreator(ctx); err != nil {
		t.Errorf("EntityCEditor.SetCreator() error = %v", err)
	}
	if e.CreatedBy == nil || *e.CreatedBy != "testuser2" {
		t.Error("EntityCEditor.SetCreator() failed to set CreatedBy")
	}

	if err := e.SetEditor(ctx); err != nil {
		t.Errorf("EntityCEditor.SetEditor() error = %v", err)
	}
	if e.UpdatedBy == nil || *e.UpdatedBy != "testuser2" {
		t.Error("EntityCEditor.SetEditor() failed to set UpdatedBy")
	}
}

func TestEntityCreatorInterface(t *testing.T) {
	// Test that EntityHaveCreator implements EntityCreator interface
	var _ EntityCreator = &EntityHaveCreator{}

	e := &EntityHaveCreator{}
	ctx := gotex.NewContext(context.Background(), &gotex.Gotex{UserID: "testuser3"})

	// Test interface method
	if err := e.SetCreator(ctx); err != nil {
		t.Errorf("EntityCreator.SetCreator() error = %v", err)
	}
	if e.CreatedBy == nil || *e.CreatedBy != "testuser3" {
		t.Error("EntityCreator.SetCreator() failed to set CreatedBy")
	}
}

// Helper function to compare string pointers
func compareStringPtr(a, b *string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}
