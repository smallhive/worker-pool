package worker_pool

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	workerAmount = 2
)

func TestExportWorker(t *testing.T) {
	ctx := context.Background()
	mu := sync.Mutex{}
	var total int64

	pool := NewExportWorker[int64](workerAmount)
	consumer := func(ctx context.Context, task int64) {
		mu.Lock()
		total += task
		mu.Unlock()
	}

	pool.Consume(ctx, consumer)

	producer := func(ctx context.Context, taskInput TaskChan[int64]) {
		for i := int64(0); i < 100; i++ {
			taskInput <- i
		}
	}

	pool.Produce(ctx, producer)
	assert.Equal(t, int64(4950), total)
}
