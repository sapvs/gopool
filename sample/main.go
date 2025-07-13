package main

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/sapvs/gopool"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	pool := gopool.New(gopool.WithNumWorkers(2),
		gopool.WithContext(ctx),
		gopool.WithLogger(slog.Default()))

	resultChan, err := pool.Start()
	if err != nil {
		return
	}

	go func() {
		for result := range resultChan {
			slog.Info("output reader", "result from pool", result.Result())
		}
	}()

	for i := range 5 {
		task := &ATask{Message: fmt.Sprintf("task number %d", i)}
		err = pool.Submit(task)
		if err != nil {
			return
		}
	}

	cancel() // OR with pool.Shutdown()
	if err != nil {
		slog.Info("", "err", err)
	}

	time.Sleep(2 * time.Second)

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

func (p *AResult) Result() any {
	return p.result
}
