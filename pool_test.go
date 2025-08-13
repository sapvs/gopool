package gopool

import (
	"context"
	"log/slog"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var cpus int = runtime.NumCPU()

func TestConfigDefaults(t *testing.T) {
	t.Parallel()
	// no Options, values should be default
	p := New()
	assert.Equal(t, uint(cpus), p.numWorkers, "num workers wanted %d; got %d", cpus, p.numWorkers)
	assert.Equal(t, cpus, cap(p.workChan), "work chan cap wanted %d; got %d", cpus, cap(p.workChan))
	assert.Equal(t, cpus, cap(p.resultChan), "result chan cap wanted %d; got %d", cpus, cap(p.resultChan))
}

func TestConfigInvalid(t *testing.T) {
	t.Parallel()

	p := New(WithContext(nil), WithLogger(nil))
	assert.Equal(t, uint(cpus), p.numWorkers, "num workers wanted %d; got %d", cpus, p.numWorkers)
	assert.Equal(t, cpus, cap(p.workChan), "work chan cap wanted %d; got %d", cpus, cap(p.workChan))
	assert.Equal(t, cpus, cap(p.resultChan), "result chan cap wanted %d; got %d", cpus, cap(p.resultChan))

	assert.IsType(t, context.TODO(), p.poolCtx, "context wanted %d; got %d", 10, p.numWorkers)
	assert.IsType(t, &NOOPLogger{}, p.PoolLogger, "Logger wanted %d; got %d", 10, p.numWorkers)

}
func TestConfig(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(t.Context(), 3*time.Second)
	defer cancel()

	p := New(WithNumWorkers(10), WithResultBuffer(20), WithWorkBuffer(30),
		WithContext(ctx), WithLogger(slog.Default()))

	assert.Equal(t, uint(10), p.numWorkers, "num workers wanted %d; got %d", 10, p.numWorkers)
	assert.Equal(t, 30, cap(p.workChan), "work chan cap wanted %d; got %d", 30, cap(p.workChan))
	assert.Equal(t, 20, cap(p.resultChan), "result chan cap wanted %d; got %d", 20, cap(p.resultChan))
	assert.Equal(t, ctx, p.poolCtx, "context")
	assert.Equal(t, slog.Default(), p.PoolLogger, "context")
}

func TestProcessing(t *testing.T) {
	t.Parallel()

	var tasks int64 = 10_000
	var received atomic.Int64
	var resultMut sync.Mutex
	var wg sync.WaitGroup

	p := New(WithNumWorkers(10), WithResultBuffer(20), WithWorkBuffer(30))
	result, err := p.Start()
	assert.NoError(t, err, "error in pool start")

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range result {
			func() {
				resultMut.Lock()
				defer resultMut.Unlock()
				received.Add(1)
			}()
		}
	}()

	for range tasks {
		assert.NoError(t, p.Submit(&TestTask{}), "error submit")
	}

	assert.NoError(t, p.Shutdown(), "error in pool stop")
	wg.Wait()
	assert.Equal(t, tasks, received.Load(), "submitted tasks not equal")
}

type TestTask struct{}

func (w *TestTask) Do() Result {
	return &TestResult{val: time.Now().String()}
}

type TestResult struct {
	val string
}

func (r *TestResult) Get() any {
	return r.val
}

func TestSubmitOnStoppedPool(t *testing.T) {
	p := New()
	var wg sync.WaitGroup
	assert.Error(t, p.Submit(&TestTask{}))

	res, err := p.Start()
	assert.NoError(t, err)
	go func(r chan Result, w *sync.WaitGroup) {
		for range r {
			w.Done()
		}
	}(res, &wg)

	assert.NoError(t, p.Submit(&TestTask{}))
	wg.Add(1)

	assert.Error(t, p.Submit(nil))

	wg.Wait()

	assert.NoError(t, p.Shutdown())
	assert.Error(t, p.Submit(&TestTask{}))

}

func TestShutdownClosedPool(t *testing.T) {
	p := New()
	res, err := p.Start()
	assert.NoError(t, err)
	assert.NotNil(t, res)

	assert.NoError(t, p.Shutdown())
	assert.Error(t, p.Shutdown())

}

func TestStartStartedPool(t *testing.T) {
	p := New()
	res, err := p.Start()
	assert.NoError(t, err)
	assert.NotNil(t, res)
	_, err = p.Start()
	assert.Error(t, err)

}

func TestCancelWithContext(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	p := New(WithContext(ctx))
	_, err := p.Start()
	assert.NoError(t, err)

	cancel()
	assert.NoError(t, p.Shutdown())

	cancel()
}

var pool *Pool

func BenchmarkPoolNew(b *testing.B) {
	for b.Loop() {
		pool = New(WithNumWorkers(10),
			WithContext(context.TODO()),
			WithWorkBuffer(100),
			WithResultBuffer(100),
			WithLogger(&NOOPLogger{}))
		_, err := pool.Start()
		if err != nil {
			b.Fail()
		}
		err = pool.Shutdown()
		if err != nil {
			b.Fail()
		}

	}
}

func BenchmarkPoolNewParallel(b *testing.B) {
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			poolp := New(WithNumWorkers(10),
				WithContext(context.TODO()),
				WithWorkBuffer(100),
				WithResultBuffer(100),
				WithLogger(&NOOPLogger{}))
			_, err := poolp.Start()
			if err != nil {
				b.Fail()
			}
			err = poolp.Shutdown()
			if err != nil {
				b.Fail()
			}
		}
	})
}

func BenchmarkPoolStart(b *testing.B) {
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			b.Log("in next")

		}
	})
}
