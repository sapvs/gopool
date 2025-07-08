package main

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/sapvs/gopool"
	"github.com/sapvs/gopool/sample"
)

func main1() {
	ch := make(chan int, 100)
	go func(ch1 chan int) {
		slog.Info("Starting for loop")
		for val := range ch1 {
			fmt.Printf("val: %v\n", val)
		}
		slog.Info("Exitted for loop, channel closed?")
	}(ch)

	for i := range 10 {
		ch <- i
	}

	time.Sleep(1 * time.Second)

	slog.Info("slept")
	for i := range 10 {
		ch <- i
	}
	close(ch)
	slog.Info("closed channel")

	time.Sleep(1 * time.Second)

}

func main() {
	pool := gopool.New(gopool.WithNumWorkers(2))

	resultChan, err := pool.Start()
	if err != nil {
		return
	}

	go func() {
		for result := range resultChan {
			slog.Info("output reader", "result from pool", result.Result())
		}
		slog.Info("output channel closed")
	}()

	for i := range 5 {
		pool.Submit(&sample.ATask{Message: fmt.Sprintf("task number %d", i)})
	}

	// time.Sleep(1 * time.Second)

	slog.Info("shutting now")
	err = pool.Shutdown()
	if err != nil {
		slog.Info("", err)
	}

	time.Sleep(2 * time.Second)

}
