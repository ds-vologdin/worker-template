package task

import "context"

type Task interface {
	Run(ctx context.Context, done chan int)
}
