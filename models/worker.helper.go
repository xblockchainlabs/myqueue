package models

import (
	"errors"
	"time"

	"github.com/xblockchainlabs/myqueue/utils"
)

func SelectJob(name string, limit int) (schds []Schedule, err error) {
	db, err := utils.GetDB()
	if err != nil {
		return
	}
	err = db.Preload("Job").Where("name = ? AND status = ? AND execute_due_at <= ?", name, Open, time.Now().UTC()).Order("id").Limit(limit).Find(&schds).Error
	return
}

func AcquireJob(id uint) (ok bool, err error) {
	db, err := utils.GetDB()
	if err != nil {
		return
	}

	now := time.Now().UTC()
	query := db.Exec("UPDATE schedules SET updated_at = ?, started_at = ?, status = ? WHERE id = ? AND status = ?", now, now, InProcess, id, Open)
	rowsAffected, err := query.RowsAffected, query.Error
	if err != nil {
		ok = false
	} else if rowsAffected == 0 {
		err = errors.New("Cannot acquire the Job")
		ok = false
	} else {
		ok = true
	}
	return
}

func CompleteJob(sid uint) (err error) {
	db, err := utils.GetDB()
	if err != nil {
		return
	}
	tx, connection := NewTx("CompleteJob", db)
	defer tx.Close()
	s := &Schedule{}

	err = connection.Preload("Job").Where("ID = ? AND status = ?", sid, InProcess).First(s).Error
	if err != nil {
		tx.Fail(err)
		return
	}
	s.Status = Success
	s.EndedAt = time.Now().UTC()
	s.Job.Status = CompletedJob
	if err = s.Update(connection); err != nil {
		tx.Fail(err)
	}
	return
}

func FailJob(sid uint, backoff *utils.Backoff) (s *Schedule, err error) {
	db, err := utils.GetDB()
	if err != nil {
		return
	}
	tx, connection := NewTx("CompleteJob", db)
	defer tx.Close()
	s = &Schedule{}

	err = connection.Preload("Job").Where("ID = ? AND status = ?", sid, InProcess).First(s).Error
	if err != nil {
		tx.Fail(err)
		return
	}
	s.Status = Fail
	s.EndedAt = time.Now().UTC()
	nextAttempt := int(s.Job.Attempt + 1)
	if !backoff.AttemptAllowed(nextAttempt) {
		s.Job.Status = FailedJob
	}

	if err = s.Update(connection); err != nil {
		tx.Fail(err)
	}
	return

}
