// gopool is a worker pool implementation.
package gopool

import (
	"log/slog"
	"runtime"
	"sync"
	"sync/atomic"
)

type pool struct {
	tasksChan                        chan Task
	resultChan                       chan Result
	taskChanBuffer, resultChanBuffer uint
	workers                          uint8
	l                                PoolLogger
	running                          atomic.Bool
	wg                               sync.WaitGroup
}

type option func(*pool)

func New(opts ...option) *pool {
	// Default values
	pool := &pool{
		workers:          uint8(runtime.NumCPU()),
		l:                slog.Default(),
		taskChanBuffer:   uint(runtime.NumCPU()),
		resultChanBuffer: uint(runtime.NumCPU()),
	}

	for _, opt := range opts {
		opt(pool)
	}
	return pool
}

// WithNumWorkers configures number of goroutines to be available in this pool, default runtime.NumCPU
func WithNumWorkers(workercount uint8) option {
	return func(p *pool) {
		p.workers = workercount
	}
}

// WithInputChannelBuffer The size of task input queue, default runtime.NumCPU
func WithInputChannelBuffer(buffer uint) option {
	return func(p *pool) {
		p.taskChanBuffer = buffer
	}
}

// WithResultChannerBuffer The size of output queue, set this such the pool worker are not blocked returning the processing result default runtime.NumCPU
func WithResultChannerBuffer(buffer uint) option {
	return func(p *pool) {
		p.resultChanBuffer = buffer
	}
}

// WithLogger confure logger for the pool, implement PoolLogger interface for this.
func WithLogger(poolLogger PoolLogger) option {
	return func(p *pool) {
		p.l = poolLogger
	}
}

// Start Starts the pool workers. Returns the output channel and error if any
func (p *pool) Start() (chan Result, error) {
	p.tasksChan = make(chan Task, p.workers)
	p.resultChan = make(chan Result, p.workers)

	if p.running.Load() {
		return nil, p.logAndReturnError(ERR_START_ON_RUNNING_POOL)
	}

	for i := range p.workers {
		p.wg.Add(1)
		go p.work(i)
	}

	p.running.Store(true)
	return p.resultChan, nil
}

// Submit submits a Task to the pool, returns error if any
func (p *pool) Submit(input Task) error {
	if input == nil {
		return p.logAndReturnError(ERR_NIL_TASK)
	}

	if !p.running.Load() {
		return p.logAndReturnError(ERR_SUBMIT_IN_CLOSED_POOL)
	}

	p.l.Debug("submitting task to pool")
	p.tasksChan <- input
	return nil

}

// Shutdown Shuts down the pool workers, currently only graceful shutdown is supported.
// In-flight taks are processed and results returned before pool shutdown
func (p *pool) Shutdown() error {

	if !p.running.Load() {
		return p.logAndReturnError(ERR_STOP_ON_CLOSED_POOL)
	}

	// TODO should need some mutex here?
	p.running.Store(false)

	p.l.Info("shutting down")
	close(p.tasksChan)

	p.l.Info("waiting for workers to finish")
	p.wg.Wait()
	p.l.Info("all workers finished, closing result channel")
	close(p.resultChan)
	p.l.Info("Done")
	return nil
}

func (p *pool) logAndReturnError(err error) error {
	p.l.Warn(err.Error())
	return err
}

func (p *pool) work(goRoutineID uint8) {
	for tasks := range p.tasksChan {
		result := tasks.Do()
		p.resultChan <- result
	}
	p.l.Info("task channel closed", "routine", goRoutineID)
	p.wg.Done()
}
