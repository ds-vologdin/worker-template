package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

const maxHandleTaskSecond = 10
const timeoutHandler = 5 * time.Second
const taskSize = 13
const batchSize = 5

type Task interface {
	Run(ctx context.Context, done chan int)
}

type TaskFast struct {
	Name string
}

func (t TaskFast) Run(ctx context.Context, done chan int) {
	select {
	case <-time.After(10 * time.Microsecond):
		log.Printf("fastOperation '%s' finish\n", t.Name)
		done <- 1
		return
	case <-ctx.Done():
		log.Printf("fastOperation '%s' cancel: %v\n", t.Name, ctx.Err())
		return
	}
}

func (t TaskFast) String() string {
	return t.Name
}

type TaskSlow struct {
	Name string
}

func (t TaskSlow) Run(ctx context.Context, done chan int) {
	second := rand.Intn(maxHandleTaskSecond)
	timeout := time.Duration(second) * time.Second
	select {
	case <-time.After(timeout):
		log.Printf("slowOperation '%s' finish\n", t.Name)
		done <- 1
		return
	case <-ctx.Done():
		log.Printf("slowOperation '%s' cancel: %v\n", t.Name, ctx.Err())
		return
	}
}

func (t TaskSlow) String() string {
	return t.Name
}

func handler(ctx context.Context, task Task) {
	log.Printf("start %v\n", task)
	done := make(chan int)
	go task.Run(ctx, done)
	select {
	case <-done:
		log.Printf("%s finish\n", task)
	case <-ctx.Done():
		log.Printf("%s cancel: %v\n", task, ctx.Err())
	}
	log.Printf("stop %v\n", task)
}

func createTasks(count int) []Task {
	tasks := make([]Task, count)
	for i := 0; i < count; i++ {
		rnd := rand.Intn(100)
		switch {
		case rnd < 50:
			tasks[i] = TaskFast{fmt.Sprintf("fast task %d", i)}
		default:
			tasks[i] = TaskSlow{fmt.Sprintf("slow task %d", i)}
		}
	}
	return tasks
}

func runBatchTask(ctx context.Context, tasks []Task) {
	var wg sync.WaitGroup
	for _, task := range tasks {
		wg.Add(1)
		go func(task Task) {
			defer wg.Done()
			ctxTimeout, cancel := context.WithTimeout(ctx, timeoutHandler)
			defer cancel()
			handler(ctxTimeout, task)
		}(task)
	}
	wg.Wait()
}

func main() {
	log.Println("start worker")
	ctx := context.Background()
	tasks := createTasks(taskSize)
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
