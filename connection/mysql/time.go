package mysql

import (
	"gorm.io/gorm"
	"time"
)

type CustomTime time.Time

func (t *CustomTime) Scan(v interface{}) error {
	vt, err := time.Parse("2006-01-02 15:04:05", string(v.([]byte)))
	if err != nil {
		return err
	}
	*t = CustomTime(vt)
	return nil
}

func (t *CustomTime) ToTime() *time.Time {
	return (*time.Time)(t)
}

func DeletedAt() gorm.DeletedAt {
	return gorm.DeletedAt{Time: time.Now(), Valid: true}
}

func DeletedAtFromTime(t *time.Time) gorm.DeletedAt {
	if t == nil {
		return gorm.DeletedAt{Valid: false}
	}
	return gorm.DeletedAt{Time: *t, Valid: true}
}

func DeletedAtToTime(t gorm.DeletedAt) *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}
