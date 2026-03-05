package time

import (
	"testing"
	"time"
)

func TestBod(t *testing.T) {
	tests := []struct {
		name string
		t    time.Time
		want time.Time
	}{
		{
			name: "midnight stays midnight",
			t:    time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			want: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "noon goes to midnight",
			t:    time.Date(2024, 1, 15, 12, 30, 45, 123456789, time.UTC),
			want: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "end of day goes to beginning",
			t:    time.Date(2024, 1, 15, 23, 59, 59, 999999999, time.UTC),
			want: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "middle of day",
			t:    time.Date(2024, 6, 30, 15, 30, 0, 0, time.UTC),
			want: time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "leap year",
			t:    time.Date(2024, 2, 29, 12, 0, 0, 0, time.UTC),
			want: time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "different timezone - positive offset",
			t:    time.Date(2024, 1, 15, 12, 0, 0, 0, time.FixedZone("EST+5", 5*3600)),
			want: time.Date(2024, 1, 15, 0, 0, 0, 0, time.FixedZone("EST+5", 5*3600)),
		},
		{
			name: "different timezone - negative offset",
			t:    time.Date(2024, 1, 15, 12, 0, 0, 0, time.FixedZone("PST-8", -8*3600)),
			want: time.Date(2024, 1, 15, 0, 0, 0, 0, time.FixedZone("PST-8", -8*3600)),
		},
		{
			name: "year boundary",
			t:    time.Date(2024, 12, 31, 18, 30, 0, 0, time.UTC),
			want: time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "month boundary",
			t:    time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
			want: time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "with nanoseconds",
			t:    time.Date(2024, 3, 15, 8, 15, 30, 500000000, time.UTC),
			want: time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Bod(tt.t)
			if !got.Equal(tt.want) {
				t.Errorf("Bod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEod(t *testing.T) {
	tests := []struct {
		name string
		t    time.Time
		want time.Time
	}{
		{
			name: "midnight goes to end of day",
			t:    time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			want: time.Date(2024, 1, 15, 23, 59, 59, 999999999, time.UTC),
		},
		{
			name: "noon stays end of day",
			t:    time.Date(2024, 1, 15, 12, 30, 45, 123456789, time.UTC),
			want: time.Date(2024, 1, 15, 23, 59, 59, 999999999, time.UTC),
		},
		{
			name: "end of day stays end of day",
			t:    time.Date(2024, 1, 15, 23, 59, 59, 999999999, time.UTC),
			want: time.Date(2024, 1, 15, 23, 59, 59, 999999999, time.UTC),
		},
		{
			name: "middle of day",
			t:    time.Date(2024, 6, 30, 15, 30, 0, 0, time.UTC),
			want: time.Date(2024, 6, 30, 23, 59, 59, 999999999, time.UTC),
		},
		{
			name: "leap year",
			t:    time.Date(2024, 2, 29, 12, 0, 0, 0, time.UTC),
			want: time.Date(2024, 2, 29, 23, 59, 59, 999999999, time.UTC),
		},
		{
			name: "different timezone - positive offset",
			t:    time.Date(2024, 1, 15, 12, 0, 0, 0, time.FixedZone("EST+5", 5*3600)),
			want: time.Date(2024, 1, 15, 23, 59, 59, 999999999, time.FixedZone("EST+5", 5*3600)),
		},
		{
			name: "different timezone - negative offset",
			t:    time.Date(2024, 1, 15, 12, 0, 0, 0, time.FixedZone("PST-8", -8*3600)),
			want: time.Date(2024, 1, 15, 23, 59, 59, 999999999, time.FixedZone("PST-8", -8*3600)),
		},
		{
			name: "year boundary",
			t:    time.Date(2024, 12, 31, 18, 30, 0, 0, time.UTC),
			want: time.Date(2024, 12, 31, 23, 59, 59, 999999999, time.UTC),
		},
		{
			name: "month boundary",
			t:    time.Date(2024, 1, 31, 23, 59, 59, 0, time.UTC),
			want: time.Date(2024, 1, 31, 23, 59, 59, 999999999, time.UTC),
		},
		{
			name: "with nanoseconds",
			t:    time.Date(2024, 3, 15, 8, 15, 30, 500000000, time.UTC),
			want: time.Date(2024, 3, 15, 23, 59, 59, 999999999, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Eod(tt.t)
			if !got.Equal(tt.want) {
				t.Errorf("Eod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBodEodConsistency(t *testing.T) {
	// Test that Bod and Eod are consistent for various times
	testTimes := []time.Time{
		time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 6, 30, 12, 30, 45, 123456789, time.UTC),
		time.Date(2024, 12, 31, 23, 59, 59, 999999999, time.UTC),
		time.Now(),
	}

	for _, tt := range testTimes {
		bod := Bod(tt)
		eod := Eod(tt)

		// Bod should have hour, minute, second, nanosecond all zero
		if bod.Hour() != 0 || bod.Minute() != 0 || bod.Second() != 0 || bod.Nanosecond() != 0 {
			t.Errorf("Bod() should have time 00:00:00.000000000, got %02d:%02d:%02d.%09d",
				bod.Hour(), bod.Minute(), bod.Second(), bod.Nanosecond())
		}

		// Eod should have hour 23, minute 59, second 59, nanosecond 999999999
		if eod.Hour() != 23 || eod.Minute() != 59 || eod.Second() != 59 || eod.Nanosecond() != 999999999 {
			t.Errorf("Eod() should have time 23:59:59.999999999, got %02d:%02d:%02d.%09d",
				eod.Hour(), eod.Minute(), eod.Second(), eod.Nanosecond())
		}

		// Both should have same year, month, day, and location
		if bod.Year() != eod.Year() || bod.Month() != eod.Month() || bod.Day() != eod.Day() {
			t.Errorf("Bod() and Eod() should have same date: Bod=%v, Eod=%v", bod, eod)
		}

		if bod.Location() != eod.Location() {
			t.Errorf("Bod() and Eod() should have same location: Bod=%v, Eod=%v", bod.Location(), eod.Location())
		}

		// Bod should be before Eod
		if !bod.Before(eod) && !bod.Equal(eod) {
			t.Errorf("Bod() should be before or equal to Eod(): Bod=%v, Eod=%v", bod, eod)
		}
	}
}
