package slow

import (
	"context"
	"log"
	"math/rand"
	"time"
)

const maxHandleTaskSecond = 10

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
