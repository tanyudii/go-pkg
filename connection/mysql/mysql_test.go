package mysql

import (
	"reflect"
	"testing"
	"time"

	"gorm.io/gorm"
)

func TestCustomTime_Scan(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    CustomTime
		wantErr bool
	}{
		{
			name:    "valid datetime",
			input:   []byte("2024-01-15 10:30:45"),
			want:    CustomTime(time.Date(2024, 1, 15, 10, 30, 45, 0, time.Local)),
			wantErr: false,
		},
		{
			name:    "valid datetime with single digit month",
			input:   []byte("2024-01-05 08:09:10"),
			want:    CustomTime(time.Date(2024, 1, 5, 8, 9, 10, 0, time.Local)),
			wantErr: false,
		},
		{
			name:    "invalid format",
			input:   []byte("2024-01-15T10:30:45Z"),
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   []byte(""),
			wantErr: true,
		},
		{
			name:    "garbage data",
			input:   []byte("not a date"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ct := &CustomTime{}
			err := ct.Scan(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("CustomTime.Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(*ct, tt.want) {
				t.Errorf("CustomTime.Scan() = %v, want %v", *ct, tt.want)
			}
		})
	}
}

func TestCustomTime_ToTime(t *testing.T) {
	tests := []struct {
		name string
		ct   CustomTime
		want *time.Time
	}{
		{
			name: "zero time",
			ct:   CustomTime(time.Time{}),
			want: &time.Time{},
		},
		{
			name: "specific date",
			ct:   CustomTime(time.Date(2024, 1, 15, 10, 30, 45, 0, time.Local)),
			want: func() *time.Time {
				t := time.Date(2024, 1, 15, 10, 30, 45, 0, time.Local)
				return &t
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.ct.ToTime()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CustomTime.ToTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeletedAt(t *testing.T) {
	result := DeletedAt()

	if !result.Valid {
		t.Error("DeletedAt() should have Valid = true")
	}

	// Check that time is within reasonable range (last minute)
	if time.Since(result.Time) > time.Minute {
		t.Error("DeletedAt() time should be recent")
	}
}

func TestDeletedAtFromTime(t *testing.T) {
	tests := []struct {
		name string
		t    *time.Time
		want gorm.DeletedAt
	}{
		{
			name: "nil time",
			t:    nil,
			want: gorm.DeletedAt{Valid: false},
		},
		{
			name: "valid time",
			t: func() *time.Time {
				t := time.Date(2024, 1, 15, 10, 30, 45, 0, time.Local)
				return &t
			}(),
			want: gorm.DeletedAt{
				Time:  time.Date(2024, 1, 15, 10, 30, 45, 0, time.Local),
				Valid: true,
			},
		},
		{
			name: "zero time",
			t:    func() *time.Time { t := time.Time{}; return &t }(),
			want: gorm.DeletedAt{
				Time:  time.Time{},
				Valid: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DeletedAtFromTime(tt.t)
			if got.Valid != tt.want.Valid {
				t.Errorf("DeletedAtFromTime() Valid = %v, want %v", got.Valid, tt.want.Valid)
				return
			}
			if tt.want.Valid && !got.Time.Equal(tt.want.Time) {
				t.Errorf("DeletedAtFromTime() Time = %v, want %v", got.Time, tt.want.Time)
			}
		})
	}
}

func TestDeletedAtToTime(t *testing.T) {
	tests := []struct {
		name string
		t    gorm.DeletedAt
		want *time.Time
	}{
		{
			name: "invalid deleted at",
			t:    gorm.DeletedAt{Valid: false},
			want: nil,
		},
		{
			name: "valid deleted at",
			t: gorm.DeletedAt{
				Time:  time.Date(2024, 1, 15, 10, 30, 45, 0, time.Local),
				Valid: true,
			},
			want: func() *time.Time {
				t := time.Date(2024, 1, 15, 10, 30, 45, 0, time.Local)
				return &t
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DeletedAtToTime(tt.t)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeletedAtToTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
