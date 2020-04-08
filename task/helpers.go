package task

import (
	"fmt"
	"math/rand"

	"github.com/ds-vologdin/worker-template/task/fast"
	"github.com/ds-vologdin/worker-template/task/slow"
)

func CreateTasks(count int) []Task {
	tasks := make([]Task, count)
	for i := 0; i < count; i++ {
		rnd := rand.Intn(100)
		switch {
		case rnd < 50:
			tasks[i] = fast.TaskFast{Name: fmt.Sprintf("fast task %d", i)}
		default:
			tasks[i] = slow.TaskSlow{Name: fmt.Sprintf("slow task %d", i)}
		}
	}
	return tasks
}
