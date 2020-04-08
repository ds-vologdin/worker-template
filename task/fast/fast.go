package fast

import (
	"context"
	"log"
	"time"
)

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
