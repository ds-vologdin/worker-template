package fast

import (
	"context"
	"log"
	"time"
)

type TaskFast struct {
	Name string
}

func (t TaskFast) Run(ctx context.Context) error {
	select {
	case <-time.After(10 * time.Microsecond):
		log.Printf("fastOperation '%s' finish\n", t.Name)
		return nil
	case <-ctx.Done():
		log.Printf("fastOperation '%s' cancel: %v\n", t.Name, ctx.Err())
		return ctx.Err()
	}
}

func (t TaskFast) String() string {
	return t.Name
}
