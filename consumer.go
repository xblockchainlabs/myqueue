package myqueue

import (
	"errors"
	"sync"

	"github.com/xblockchainlabs/myqueue/utils"
)

type ConsumerGroup struct {
	once        sync.Once
	started     bool
	workerPools []*Pool
}

func NewCG() *ConsumerGroup {
	cg := &ConsumerGroup{started: false}
	return cg
}

func (cg *ConsumerGroup) AddWorker(name string, size int, backoff *utils.Backoff, workerFunc WorkerFunc) error {
	if cg.started {
		return errors.New("Cannot add new worker")
	}
	wp, err := Worker(name, size, backoff, workerFunc)
	if err != nil {
		return err
	}
	cg.workerPools = append(cg.workerPools, wp)
	return nil
}

func (cg *ConsumerGroup) Start() {
	cg.once.Do(func() {
		for _, wp := range cg.workerPools {
			wp.Start(Allocator, Collector)
		}
		cg.started = true
	})
}
