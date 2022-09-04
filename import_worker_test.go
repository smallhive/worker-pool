package worker_pool

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	chanBufferSize = 10
)

func TestImportWorker(t *testing.T) {
	ctx := context.Background()
	var total int64

	pool := NewImportWorker[int64](chanBufferSize)
	producer := func(ctx context.Context, taskInput TaskChan[int64]) {
		for i := int64(0); i < 100; i++ {
			taskInput <- i
		}
	}

	pool.Produce(ctx, producer)
	consumer := func(ctx context.Context, task int64) {
		total += task
	}

	pool.Consume(ctx, consumer)
	assert.Equal(t, int64(4950), total)
}

func TestImportWorker2Generators(t *testing.T) {
	ctx := context.Background()
	var total int64

	pool := NewImportWorker[int64](workerAmount)
	pool.Produce(ctx, func(ctx context.Context, taskInput TaskChan[int64]) {
		for i := int64(0); i < 50; i++ {
			taskInput <- i
		}
	})

	pool.Produce(ctx, func(ctx context.Context, taskInput TaskChan[int64]) {
		for i := int64(50); i < 100; i++ {
			taskInput <- i
		}
	})

	consumer := func(ctx context.Context, task int64) {
		total += task
	}

	pool.Consume(ctx, consumer)
	assert.Equal(t, int64(4950), total)
}
