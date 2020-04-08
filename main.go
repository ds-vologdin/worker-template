package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/ds-vologdin/worker-template/task"
)

const timeoutHandler = 3 * time.Second
const taskSize = 13
const batchSize = 5

func handler(ctx context.Context, taskCurrent task.Task) {
	log.Printf("start %v\n", taskCurrent)
	done := make(chan int)
	go taskCurrent.Run(ctx, done)
	select {
	case <-done:
		log.Printf("%s finish\n", taskCurrent)
	case <-ctx.Done():
		log.Printf("%s cancel: %v\n", taskCurrent, ctx.Err())
	}
	log.Printf("stop %v\n", taskCurrent)
}

func runBatchTask(ctx context.Context, tasks []task.Task) {
	var wg sync.WaitGroup
	for _, taskCurrent := range tasks {
		wg.Add(1)
		go func(taskCurrent task.Task) {
			defer wg.Done()
			ctxTimeout, cancel := context.WithTimeout(ctx, timeoutHandler)
			defer cancel()
			handler(ctxTimeout, taskCurrent)
		}(taskCurrent)
	}
	wg.Wait()
}

func main() {
	log.Println("start worker")
	ctx := context.Background()
	tasks := task.CreateTasks(taskSize)
	for i := 0; i < len(tasks); i = i + batchSize {
		end := i + batchSize
		if end > len(tasks) {
			end = len(tasks)
		}
		runBatchTask(ctx, tasks[i:end])
		log.Println("-----------------------------")
		time.Sleep(1 * time.Second)
	}
}
