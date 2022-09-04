package worker_pool

import (
	"context"
)

type TaskChan[T any] chan T

type ConsumerFunc[T any] func(ctx context.Context, task T)

type ProducerFunc[T any] func(ctx context.Context, taskInput TaskChan[T])
