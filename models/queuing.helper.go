package models

import (
	"errors"
	"time"

	"github.com/xblockchainlabs/myqueue/utils"
)

func SchedJob(j Job, exeDueAt time.Time) (id uint, err error) {
	db, err := utils.GetDB()
	if err != nil {
		return
	}

	tx, connection := NewTx("SchedJob", db)
	defer tx.Close()
	if err = j.Save(connection); err != nil {
		tx.Fail(err)
		return
	}

	s := &Schedule{}
	s.Name = j.Name
	s.JobID = j.ID
	s.Attempt = j.Attempt
	s.ExecuteDueAt = exeDueAt
	if err = s.Save(connection); err != nil {
		tx.Fail(err)
		return
	}

	id = j.ID
	return
}

func ReschedJob(jID uint, exeDueAt time.Time) (err error) {
	db, err := utils.GetDB()
	if err != nil {
		return
	}
	j := &Job{}
	tx, connection := NewTx("ReschedJob", db)
	defer tx.Close()

	err = connection.Where("ID = ?", jID).First(j).Error
	if err != nil {
		tx.Fail(err)
		return
	}
	var pending uint
	connection.Model(&Schedule{}).Where("job_id = ? AND status < ?", j.ID, Success).Count(&pending)
	if pending > 0 {
		err = errors.New("Previous schedule task is pending")
		tx.Fail(err)
		return
	}

	nextAttempt := j.Attempt + 1
	s := &Schedule{}
	s.Name = j.Name
	s.JobID = j.ID
	s.Attempt = nextAttempt
	s.ExecuteDueAt = exeDueAt
	if err = s.Save(connection); err != nil {
		tx.Fail(err)
		return
	}

	query := connection.Exec("UPDATE jobs SET updated_at = ?, attempt = ? WHERE id = ? AND attempt = ?", time.Now().UTC(), nextAttempt, j.ID, j.Attempt)
	rowsAffected, err := query.RowsAffected, query.Error
	if err != nil {
		tx.Fail(err)
	}
	if rowsAffected == 0 {
		err = errors.New("Already rescheduled or reattempted ")
		tx.Fail(err)
	}
	return
}
