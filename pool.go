// gopool is a worker pool implementation.
package gopool

import (
	"context"
	"sync"
	"sync/atomic"
)

type Pool struct {
	PoolLogger
	workChan   chan Work
	resultChan chan Result
	numWorkers uint
	running    atomic.Bool
	wg         sync.WaitGroup
	poolCtx    context.Context
	mu         sync.Mutex
}

// Start Starts the pool workers. Returns the output channel and error if any
func (p *Pool) Start() (chan Result, error) {

	p.mu.Lock()
	defer p.mu.Unlock()

	if p.running.Load() {
		return nil, p.logErr(ERR_START_ON_RUNNING_POOL)
	}
	defer p.running.Store(true)

	for i := range p.numWorkers {
		p.wg.Add(1)
		go p.work(i)
	}

	go p.monitorContext()

	return p.resultChan, nil
}

// Submit submits a Task to the pool, returns error if any
func (p *Pool) Submit(input Work) error {
	if input == nil {
		return p.logErr(ERR_NIL_WORK)
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.running.Load() {
		return p.logErr(ERR_SUBMIT_IN_CLOSED_POOL)
	}

	p.Debug("submitting task to pool")
	p.workChan <- input
	return nil

}

// Shutdown Shuts down the pool workers, currently only graceful shutdown is supported.
// In-flight taks are processed and results returned before pool shutdown
func (p *Pool) Shutdown() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.running.Load() {
		return p.logErr(ERR_STOP_ON_CLOSED_POOL)
	}
	defer p.running.Store(false)
	p.Info("shutting down")
	close(p.workChan)
	p.Info("waiting for workers to finish")
	p.wg.Wait()
	p.Info("all workers finished, closing result channel")
	close(p.resultChan)
	p.Info("Done")

	return nil
}

func (p *Pool) logErr(err error) error {
	p.Warn(err.Error())
	return err
}

func (p *Pool) work(goRoutineID uint) {
	for work := range p.workChan {
		result := work.Do()
		p.resultChan <- result
	}
	p.Info("task channel closed", "routine", goRoutineID)
	p.wg.Done()
}

func (p *Pool) monitorContext() {
	if p.poolCtx == nil || p.poolCtx.Done() == nil {
		p.Warn("nil context not monitoring cancel channel")
		return
	}

	<-p.poolCtx.Done()

	p.Info("context done called, shutting down the pool")
	if err := p.Shutdown(); err != nil {
		p.Info("could not shutdown pool", "err", err)
	}
}
