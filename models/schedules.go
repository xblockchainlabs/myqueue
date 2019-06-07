package models

import (
	"database/sql/driver"
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

type ScheduleStatus uint

const (
	Open      ScheduleStatus = 0
	InProcess ScheduleStatus = 1
	Success   ScheduleStatus = 2
	Fail      ScheduleStatus = 3
)

func (s *ScheduleStatus) Scan(value interface{}) error { *s = ScheduleStatus(value.(int64)); return nil }

func (s ScheduleStatus) Value() (driver.Value, error) { return int64(s), nil }

type Schedule struct {
	gorm.Model
	Name         string         `gorm:"column:name;index:idx_name;not null"`
	ExecuteDueAt time.Time      `gorm:"column:execute_due_at;DEFAULT:current_timestamp"`
	StartedAt    time.Time      `gorm:"column:started_at;DEFAULT:NULL"`
	EndedAt      time.Time      `gorm:"column:end_at;DEFAULT:NULL"`
	Status       ScheduleStatus `sql:"default:0"`
	Attempt      int8           `gorm:"unique_index:idx_jid_atmpt"`
	JobID        uint           `gorm:"unique_index:idx_jid_atmpt"`
	Job          Job
}

func (s *Schedule) IsEmpty() bool {
	return len(s.Name) == 0 || s.JobID == 0
}

func (s *Schedule) Save(db *gorm.DB) (err error) {
	err = db.Create(s).Error
	return
}

func (s *Schedule) Update(db *gorm.DB) (err error) {
	err = db.Save(s).Error
	return
}

func (s *Schedule) BeforeSave() (err error) {
	if s.ExecuteDueAt.IsZero() {
		err = errors.New("`Execution Due Timestamp` cannot be zero value")
	}
	return
}
