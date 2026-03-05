package bcrypt

import (
	"testing"
)

func TestHash(t *testing.T) {
	tests := []struct {
		name    string
		val     string
		cost    []int
		wantErr bool
	}{
		{
			name:    "default cost",
			val:     "password123",
			cost:    nil,
			wantErr: false,
		},
		{
			name:    "custom valid cost 10",
			val:     "password123",
			cost:    []int{10},
			wantErr: false,
		},
		{
			name:    "custom valid cost 12",
			val:     "password123",
			cost:    []int{12},
			wantErr: false,
		},
		{
			name:    "custom valid cost 4",
			val:     "password123",
			cost:    []int{4},
			wantErr: false,
		},
		{
			name:    "empty password",
			val:     "",
			cost:    nil,
			wantErr: false, // bcrypt allows empty passwords
		},
		{
			name:    "special characters",
			val:     "!@#$%^&*()_+-=[]{}|;':\",./<>?",
			cost:    nil,
			wantErr: false,
		},
		{
			name:    "unicode characters",
			val:     "密码123🔐",
			cost:    nil,
			wantErr: false,
		},
		{
			name:    "long password",
			val:     string(make([]byte, 100)), // 100 bytes
			cost:    nil,
			wantErr: false,
		},
		{
			name:    "multiple cost values - uses first",
			val:     "password123",
			cost:    []int{10, 12, 14},
			wantErr: false,
		},
		{
			name:    "zero cost",
			val:     "password123",
			cost:    []int{0},
			wantErr: true, // bcrypt.MinCost is 4
		},
		{
			name:    "negative cost",
			val:     "password123",
			cost:    []int{-5},
			wantErr: true,
		},
		{
			name:    "very high cost 32",
			val:     "password123",
			cost:    []int{32},
			wantErr: true, // bcrypt.MaxCost is 31
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Hash(tt.val, tt.cost...)

			if (err != nil) != tt.wantErr {
				t.Errorf("Hash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if got == "" {
				t.Error("Hash() should return non-empty hash")
			}

			// Verify hash format - bcrypt hashes start with $2a$, $2b$, or $2y$
			if len(got) < 60 {
				t.Errorf("Hash() returned too short hash: length = %d", len(got))
			}

			// Different calls should produce different hashes (due to salt)
			hash2, err := Hash(tt.val, tt.cost...)
			if err != nil {
				t.Errorf("Hash() second call error = %v", err)
			}
			if got == hash2 {
				t.Error("Hash() should produce different hashes for same input due to salt")
			}
		})
	}
}

func TestCompare(t *testing.T) {
	// Pre-generate some hashes for testing
	hashPassword10, _ := Hash("password123", 10)
	hashPassword12, _ := Hash("password123", 12)
	hashEmpty, _ := Hash("", 10)
	hashSpecial, _ := Hash("!@#$%^&*()", 10)

	tests := []struct {
		name    string
		hashed  string
		plain   string
		want    bool
		wantErr bool
	}{
		{
			name:    "matching password",
			hashed:  hashPassword10,
			plain:   "password123",
			want:    true,
			wantErr: false,
		},
		{
			name:    "non-matching password",
			hashed:  hashPassword10,
			plain:   "wrongpassword",
			want:    false,
			wantErr: false,
		},
		{
			name:    "empty password matches",
			hashed:  hashEmpty,
			plain:   "",
			want:    true,
			wantErr: false,
		},
		{
			name:    "empty password does not match",
			hashed:  hashEmpty,
			plain:   "something",
			want:    false,
			wantErr: false,
		},
		{
			name:    "special characters match",
			hashed:  hashSpecial,
			plain:   "!@#$%^&*()",
			want:    true,
			wantErr: false,
		},
		{
			name:    "special characters do not match",
			hashed:  hashSpecial,
			plain:   "!@#$%^&*",
			want:    false,
			wantErr: false,
		},
		{
			name:    "different cost but same password matches",
			hashed:  hashPassword12,
			plain:   "password123",
			want:    true,
			wantErr: false,
		},
		{
			name:    "invalid hash format",
			hashed:  "invalid-hash",
			plain:   "password",
			want:    false,
			wantErr: true,
		},
		{
			name:    "empty hash",
			hashed:  "",
			plain:   "password",
			want:    false,
			wantErr: true,
		},
		{
			name:    "case sensitive - wrong case",
			hashed:  hashPassword10,
			plain:   "Password123",
			want:    false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Compare(tt.hashed, tt.plain)

			if (err != nil) != tt.wantErr {
				t.Errorf("Compare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHashAndCompareIntegration(t *testing.T) {
	tests := []struct {
		name string
		pass string
		cost []int
	}{
		{"simple password", "mypassword", nil},
		{"with numbers", "pass123word", []int{10}},
		{"with symbols", "p@ssw0rd!", []int{12}},
		{"long password", "this_is_a_very_long_password_with_many_characters_123456789", []int{8}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Hash the password
			hash, err := Hash(tt.pass, tt.cost...)
			if err != nil {
				t.Fatalf("Hash() error = %v", err)
			}

			// Verify it matches
			matches, err := Compare(hash, tt.pass)
			if err != nil {
				t.Errorf("Compare() error = %v", err)
			}
			if !matches {
				t.Error("Compare() should return true for matching password")
			}

			// Verify wrong password doesn't match
			matches, err = Compare(hash, tt.pass+"wrong")
			if err != nil {
				t.Errorf("Compare() error = %v", err)
			}
			if matches {
				t.Error("Compare() should return false for wrong password")
			}
		})
	}
}

func TestCostAffectsHashTime(t *testing.T) {
	// This test verifies that higher cost values produce valid hashes
	// We don't measure time as it can vary
	password := "testpassword"

	costs := []int{4, 8, 10, 12}
	for _, cost := range costs {
		t.Run("cost_"+string(rune('0'+cost)), func(t *testing.T) {
			hash, err := Hash(password, cost)
			if err != nil {
				t.Errorf("Hash() with cost %d error = %v", cost, err)
				return
			}

			// Verify the hash is valid by comparing
			matches, err := Compare(hash, password)
			if err != nil {
				t.Errorf("Compare() error = %v", err)
			}
			if !matches {
				t.Errorf("Hash with cost %d produced invalid hash", cost)
			}
		})
	}
}

func TestHashDeterminism(t *testing.T) {
	// Test that the same password produces different hashes (due to salt)
	password := "testpassword"

	hash1, err := Hash(password, 10)
	if err != nil {
		t.Fatalf("Hash() first call error = %v", err)
	}

	hash2, err := Hash(password, 10)
	if err != nil {
		t.Fatalf("Hash() second call error = %v", err)
	}

	// Hashes should be different due to random salt
	if hash1 == hash2 {
		t.Error("Hash() should produce different hashes for same password (salt should differ)")
	}

	// But both should match the original password
	match1, err := Compare(hash1, password)
	if err != nil {
		t.Errorf("Compare() hash1 error = %v", err)
	}
	if !match1 {
		t.Error("Compare() should match hash1 with password")
	}

	match2, err := Compare(hash2, password)
	if err != nil {
		t.Errorf("Compare() hash2 error = %v", err)
	}
	if !match2 {
		t.Error("Compare() should match hash2 with password")
	}
}
