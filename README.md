# gopool
Worker pool implementation in go where created goroutines are reused for processing supplied tasks.

_Note: This is not a sample/learning/educational implementation. This is a full fledged library to intended for use in other cases._

## Usage

```golang
    pool := gopool.New(gopool.WithNumWorkers(2))

	resultChan, err := pool.Start()
	
	go func() {
		for result := range resultChan {
			// process results 
            slog.Info("output reader", "result from pool", result.Result())
		}
        // loop/ goroutine exits when resultChan is closed.
	}()

    // Submit 5 tasks
	for i := range 5 {
		pool.Submit(&sample.ATask{Message: fmt.Sprintf("task number %d", i)})
	}

    // Shutdown pool, resultChan will be closed after in flight tasks are processed.
    err = pool.Shutdown() 
```
