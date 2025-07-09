# gopool
Worker pool implementation in go where created goroutines are reused for processing supplied tasks.

_Note: This is not a sample/learning/educational implementation. This is a full fledged library to intended for use in other cases._

## Usage

```golang
pool := gopool.New(gopool.WithNumWorkers(2))

resultChan, err := pool.Start()
	
go func() {
	for result := range resultChan {
		slog.Info("output reader", "result from pool", result.Result()) // process results 
	}
  // loop/ goroutine exits when resultChan is closed.
}()

for i := range 5 { // Submit 5 tasks
  err = pool.Submit(&sample.ATask{Message: fmt.Sprintf("task number %d", i)})
}

err = pool.Shutdown() // Shutdown pool, in flight tasks are processed.
```
