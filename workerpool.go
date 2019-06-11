package myqueue

import (
	"fmt"
	"sync"
	"time"

	"github.com/xblockchainlabs/myqueue/models"
	"github.com/xblockchainlabs/myqueue/utils"
)

// ProcessorFunc signature that defines the dependency injection to process "Jobs"
type ProcessorFunc func(sched models.Schedule) (result Result)

// ResultProcessorFunc signature that defines the dependency injection to process "Results"
type AllocationFunc func(name string, size int) []models.Schedule

type ResultFunc func(sched models.Schedule, backoff *utils.Backoff, ok bool)

// Result holds the main structure for worker processed job results.
type Result struct {
	Task models.Schedule
	Ok   bool
	Err  error
}

func (r *Result) isEmpty() bool {
	return r.Task.IsEmpty()
}

// Manager generic struct that keeps all the logic to manage the queues
type Pool struct {
	backoff  *utils.Backoff
	procFunc ProcessorFunc
	name     string
	size     int
	tasks    chan models.Schedule
	results  chan Result
	done     chan bool
}

// NewManager returns a new manager structure ready to be used.
func NewPool(name string, backoff *utils.Backoff, size int, procFunc ProcessorFunc) *Pool {
	utils.InfoLog("Creating a new Pool")
	r := &Pool{
		backoff:  backoff,
		name:     name,
		size:     size,
		procFunc: procFunc,
	}
	r.setChannels()

	return r
}

func (m *Pool) setChannels() {
	m.tasks = make(chan models.Schedule, m.size)
	m.results = make(chan Result, m.size)
	return
}

func (m *Pool) Start(allocFunc AllocationFunc, resultFunc ResultFunc) {
	utils.InfoLogf("worker pool starting\n")
	go m.allocate(allocFunc)
	m.done = make(chan bool)
	go m.collect(resultFunc)
	go m.workerPool()
	<-m.done
	m.setChannels()
	go m.Start(allocFunc, resultFunc)
	return
}

func (m *Pool) allocate(alloc AllocationFunc) {
	defer close(m.tasks)
	tasks := alloc(m.name, m.size)
	utils.InfoLogf("Allocating [%d] resources\n", len(tasks))
	for _, t := range tasks {
		m.tasks <- t
	}
	utils.InfoLogf("Done Allocating.")
}

func (m *Pool) work(wg *sync.WaitGroup) {
	defer wg.Done()
	utils.InfoLog("goRoutine work starting\n")
	to := make(chan string, 1)
	go func() {
		time.Sleep(30 * time.Second)
		to <- "timeout"
	}()
	select {
	case <-to:
		fmt.Println("It's timedout")
		m.results <- Result{}
	case t := <-m.tasks:
		if t.IsEmpty() {
			fmt.Println("It's empty")
			m.results <- Result{}

		}
		m.results <- m.procFunc(t)
		utils.InfoLog("goRoutine work done.\n")
	}
}

// workerPool creates or spawns new "work" goRoutines
func (m *Pool) workerPool() {
	defer close(m.results)
	utils.InfoLogf("Worker Pool spawning new goRoutines, total: [%d]", m.size)
	var wg sync.WaitGroup
	for i := 0; i < m.size; i++ {
		wg.Add(1)
		go m.work(&wg)
	}
	wg.Wait()
	utils.InfoLog("all work goroutines done processing")
}

// Collect post processes the channel "Results" and for further processing
func (m *Pool) collect(resultFunc ResultFunc) {
	utils.InfoLog("goRoutine collect starting")
	for r := range m.results {
		if !r.isEmpty() {
			if r.Err != nil {
				utils.ErrorLogf("Job with id: [%d] got an Error: %s", r.Task.ID, r.Err)
			}
			resultFunc(r.Task, m.backoff, r.Ok)
		} else {
			fmt.Println("Result is empty")
		}
	}
	utils.InfoLog("goRoutine collect done, setting channel done as completed")
	m.done <- true
}
