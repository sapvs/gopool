package gopool

import (
	"context"
	"fmt"
	"runtime"
)

type option func(*Pool)

func New(opts ...option) *Pool {
	// Default values
	pool := &Pool{
		numWorkers: uint(runtime.NumCPU()),
		PoolLogger: &NOOPLogger{},
		workChan:   make(chan Work, runtime.NumCPU()),
		resultChan: make(chan Result, runtime.NumCPU()),
		poolCtx:    context.TODO(),
	}

	for _, opt := range opts {
		opt(pool)
	}

	pool.Info("created worker pool", "pool config", pool.String())

	return pool
}

// WithNumWorkers configures number of goroutines in this pool, zero value defaults to runtime.NumCPU
func WithNumWorkers(workercount uint) option {
	return func(p *Pool) {
		if workercount != 0 {
			p.numWorkers = workercount
		}
	}
}

// WithContext context for this pool, may be used for cancellatins etc zero value defaults to nil
func WithContext(ctx context.Context) option {
	return func(p *Pool) {
		if ctx != nil {
			p.poolCtx = ctx
		}
	}
}

// WithWorkBuffer The size of task input queue, zero value defaults to runtime.NumCPU
func WithWorkBuffer(buffer uint) option {
	return func(p *Pool) {
		if buffer != 0 {
			p.workChan = make(chan Work, buffer)
		}
	}
}

// WithResultBuffer The size of output queue, zero value defaults to runtime.NumCPU
func WithResultBuffer(buffer uint) option {
	return func(p *Pool) {
		if buffer != 0 {
			p.resultChan = make(chan Result, buffer)
		}
	}
}

// WithLogger configure PooLogger. zero value defaults to NOOPLogger
func WithLogger(logger PoolLogger) option {
	return func(p *Pool) {
		if logger != nil {
			p.PoolLogger = logger
		}
	}
}
func (p *Pool) String() string {
	return fmt.Sprintf("go-pool {workers=%d, logger=%T, in-flight tasks=%d, in-flight results=%d }",
		p.numWorkers, p.PoolLogger, len(p.workChan), len(p.resultChan))
}
