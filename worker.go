package myqueue

import (
	"errors"
	"time"

	"github.com/xblockchainlabs/myqueue/models"
	"github.com/xblockchainlabs/myqueue/utils"
)

type WorkerFunc func([]byte) (done bool, err error)

func Allocator(name string, size int) (sechds []models.Schedule) {
	sechds, err := models.SelectJob(name, size)
	utils.InfoLogf("Schedules %#v\n", sechds)
	if err != nil {
		utils.ErrorLogf("Allocation Error: %s\n", err)
	}
	return
}

func Collector(sched models.Schedule, backoff *utils.Backoff, ok bool) {
	if ok {
		models.CompleteJob(sched.ID)
	} else {
		updatedSched, err := models.FailJob(sched.ID, backoff)
		if err == nil && updatedSched.Job.Status == models.OpenJob {
			job := updatedSched.Job
			job.Attempt += 1
			delay := backoff.GetDelay(int(job.Attempt))
			models.ReschedJob(job.ID, time.Now().UTC().Add(delay))
		}
	}
	return
}

func procClosure(worker WorkerFunc) ProcessorFunc {
	return func(s models.Schedule) (result Result) {
		defer func() {
			if r := recover(); r != nil {
				utils.FatalLog(r)
			}
		}()
		result = Result{s, false, nil}
		ok, _ := models.AcquireJob(s.ID)
		if !ok {
			result = Result{}
			return
		}
		job := s.Job
		done, err := worker(job.Params)
		result.Ok = done
		result.Err = err
		return
	}
}

func Worker(name string, size int, backoff *utils.Backoff, workerFunc WorkerFunc) (pool *Pool, err error) {
	if name == "" {
		err = errors.New("Worker name cannot be blank")
		return
	}
	if size < 1 {
		err = errors.New("Number of workers should be atleast 1")
		return
	}
	if size > 999 {
		err = errors.New("Maximum number of workers can be 999")
		return
	}
	pool = NewPool(name, backoff, size, procClosure(workerFunc))
	return
}

func ZeroBackoff() (b *utils.Backoff) {
	b, err := utils.NewBackoff(utils.NoAttempt, 0, []time.Duration{0 * time.Second}, 0)
	if err != nil {
		panic(err)
	}
	return
}

func ConstanBackoff(delay time.Duration, maxAttempt int) (b *utils.Backoff) {
	b, err := utils.NewBackoff(utils.ConstantBO, 0, []time.Duration{delay}, maxAttempt)
	if err != nil {
		panic(err)
	}
	return
}

func GeometricBackoff(delay time.Duration, mult int, maxAttempt int) (b *utils.Backoff) {
	b, err := utils.NewBackoff(utils.GeometricBO, mult, []time.Duration{delay}, maxAttempt)
	if err != nil {
		panic(err)
	}
	return
}

func ListBackoff(delays []time.Duration) (b *utils.Backoff) {
	b, err := utils.NewBackoff(utils.ConstantBO, 0, delays, len(delays))
	if err != nil {
		panic(err)
	}
	return
}
