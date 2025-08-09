package main

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/sapvs/gopool"
)

func main() {

	var err error
	fmt.Printf("err: %v\n", err)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool := gopool.New(gopool.WithNumWorkers(2),
		gopool.WithContext(ctx),
		gopool.WithLogger(slog.Default()))

	resultChan, _ := pool.Start()

	go func() {
		for result := range resultChan {
			slog.Info("result from pool", "result", result.Get())
		}
	}()

	for i := range 5 {
		task := &ATask{Message: fmt.Sprintf("task number %d", i)}
		err = pool.Submit(task)
		if err != nil {
			return
		}
	}

	pool.Shutdown()
	pool.Submit(&ATask{})
	// time.Sleep(2 * time.Second)

}

type ATask struct {
	Message string
}

func (p *ATask) Do() gopool.Result {
	return &AResult{strings.ToUpper(p.Message)}
}

type AResult struct {
	result string
}

func (p *AResult) Get() any {
	return p.result
}
