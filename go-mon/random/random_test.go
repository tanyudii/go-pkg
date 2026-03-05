package random

import (
	"strings"
	"testing"
)

func TestStringRand(t *testing.T) {
	tests := []struct {
		name       string
		size       int
		wantErr    bool
		checkLen   bool
		checkAlpha bool
	}{
		{
			name:       "valid size 5",
			size:       5,
			wantErr:    false,
			checkLen:   true,
			checkAlpha: true,
		},
		{
			name:       "valid size 10",
			size:       10,
			wantErr:    false,
			checkLen:   true,
			checkAlpha: true,
		},
		{
			name:       "valid size 32",
			size:       32,
			wantErr:    false,
			checkLen:   true,
			checkAlpha: true,
		},
		{
			name:    "too large",
			size:    33,
			wantErr: true,
		},
		{
			name:    "negative size - panics",
			size:    -1,
			wantErr: true,
		},
		{
			name:       "zero size",
			size:       0,
			wantErr:    false,
			checkLen:   true,
			checkAlpha: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Handle panic case for negative size
			if tt.name == "negative size - panics" {
				defer func() {
					if r := recover(); r == nil {
						t.Error("StringRand(-1) should panic, but it didn't")
					}
				}()
			}

			result, err := StringRand(tt.size)

			if (err != nil) != tt.wantErr {
				t.Errorf("StringRand(%d) error = %v; wantErr %v", tt.size, err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if tt.checkLen && len(result) != tt.size {
				t.Errorf("StringRand(%d) length = %d; want %d", tt.size, len(result), tt.size)
			}

			if tt.checkAlpha {
				alpha := "abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789"
				for _, c := range result {
					if !strings.ContainsRune(alpha, c) {
						t.Errorf("StringRand(%d) contains invalid character %c", tt.size, c)
					}
				}
			}
		})
	}
}

func TestStringRandUniqueness(t *testing.T) {
	// Test that multiple calls produce different results (high probability)
	results := make(map[string]bool)
	for i := 0; i < 100; i++ {
		result, err := StringRand(10)
		if err != nil {
			t.Fatalf("StringRand(10) unexpected error: %v", err)
		}
		results[result] = true
	}

	// With 100 random 10-character strings, we should have close to 100 unique values
	// (the probability of collision is extremely low)
	if len(results) < 95 {
		t.Errorf("StringRand produced only %d unique values out of 100; expected > 95", len(results))
	}
}

func TestNumberRand(t *testing.T) {
	tests := []struct {
		name    string
		size    int
		wantErr bool
	}{
		{"valid size 1", 1, false},
		{"valid size 5", 5, false},
		{"valid size 10", 10, false},
		{"valid size 18", 18, false}, // 18 is max that doesn't overflow int
		{"large size 19", 19, false}, // Doesn't error, but we can't validate range due to overflow
		{"too large", 33, true},
		{"negative size - panics", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Handle panic case for negative size
			if tt.name == "negative size - panics" {
				defer func() {
					if r := recover(); r == nil {
						t.Error("NumberRand(-1) should panic, but it didn't")
					}
				}()
			}

			result, err := NumberRand(tt.size)

			if (err != nil) != tt.wantErr {
				t.Errorf("NumberRand(%d) error = %v; wantErr %v", tt.size, err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if result <= 0 {
				t.Errorf("NumberRand(%d) = %d; want > 0", tt.size, result)
			}

			// Check that the result has the correct number of digits (skip for large sizes due to overflow)
			if tt.size <= 18 {
				min := pow10(tt.size - 1)
				max := pow10(tt.size)
				if int64(result) < min || int64(result) >= max {
					t.Errorf("NumberRand(%d) = %d; want between %d and %d", tt.size, result, min, max)
				}
			}
		})
	}
}

func TestNumberRandUniqueness(t *testing.T) {
	// Test that multiple calls produce different results
	results := make(map[int]bool)
	for i := 0; i < 50; i++ {
		result, err := NumberRand(5)
		if err != nil {
			t.Fatalf("NumberRand(5) unexpected error: %v", err)
		}
		results[result] = true
	}

	// With 50 random 5-digit numbers, we should have close to 50 unique values
	if len(results) < 40 {
		t.Errorf("NumberRand produced only %d unique values out of 50; expected > 40", len(results))
	}
}

// Helper function for power of 10 (using int64 to avoid overflow)
func pow10(n int) int64 {
	result := int64(1)
	for i := 0; i < n; i++ {
		result *= 10
	}
	return result
}

func TestErrRandomTooMuch(t *testing.T) {
	if ErrRandomTooMuch == nil {
		t.Error("ErrRandomTooMuch should not be nil")
	}

	expected := "[ERROR]: The length of randomness is too much"
	if ErrRandomTooMuch.Error() != expected {
		t.Errorf("ErrRandomTooMuch.Error() = %s; want %s", ErrRandomTooMuch.Error(), expected)
	}
}
