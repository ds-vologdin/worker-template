package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/ds-vologdin/worker-template/task"
)

const timeoutHandler = 2 * time.Second
const taskSize = 113
const batchSize = 15

func handler(ctx context.Context, taskCurrent task.Task) {
	log.Printf("Start %v\n", taskCurrent)
	err := taskCurrent.Run(ctx)
	if err != nil {
		// mark task as Failed
		log.Printf("Task %v error: %v", taskCurrent, err)
	} else {
		// mark task as Ok
		// mark next task as Pending
		log.Printf("Task %v done", taskCurrent)
	}
	log.Printf("Stop %v\n", taskCurrent)
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
	// Здесь должен быть бесконечный цикл с чтением заданий из БД.
	// Читаем таски со статусом Pending и меткой auto.
	// При чтении помечаем выбранные задачи в БД статусом Processed.
	// Если мы в базе нет задач на обработку, спим заданное время (1 минуту)
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
