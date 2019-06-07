package models

import (
	"bytes"
	"database/sql/driver"
	"errors"

	"github.com/jinzhu/gorm"
)

type JSON []byte

func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	s, ok := value.([]byte)
	if !ok {
		return errors.New("Invalid Scan Source")
	}
	*j = append((*j)[0:0], s...)
	return nil
}

func (j JSON) Value() (driver.Value, error) {
	if j.IsNull() {
		return nil, nil
	}
	return string(j), nil
}

func (j JSON) MarshalJSON() ([]byte, error) {
	if j == nil {
		return []byte("null"), nil
	}
	return j, nil
}

func (j *JSON) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("null point exception")
	}
	*j = append((*j)[0:0], data...)
	return nil
}

func (j JSON) IsNull() bool {
	return len(j) == 0 || string(j) == "null"
}

func (j1 JSON) Equals(j2 JSON) bool {
	return bytes.Equal([]byte(j1), []byte(j2))
}

type JobStatus uint

const (
	OpenJob      JobStatus = 0
	CompletedJob JobStatus = 1
	FailedJob    JobStatus = 2
)

func (s *JobStatus) Scan(value interface{}) error { *s = JobStatus(value.(int64)); return nil }
func (s JobStatus) Value() (driver.Value, error)  { return int64(s), nil }

type Job struct {
	gorm.Model
	Name     string    `gorm:"column:name;not null"`
	Params   JSON      `sql:"type:json" json:"object,omitempty"`
	Attempt  int8      `gorm:"column:attempt;default:0"`
	Status   JobStatus `sql:"default:0"`
	ParentID uint
	Parent   *Job
}

func (j *Job) Save(db *gorm.DB) (err error) {
	err = db.Create(j).Error
	return
}
