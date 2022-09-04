package worker_pool

import (
	"context"
	"sync"
)

type ImportWorker[T any] struct {
	wg         *sync.WaitGroup
	chanLength int
	tasks      TaskChan[T]
}

func NewImportWorker[T any](chanLength int) *ImportWorker[T] {
	return &ImportWorker[T]{
		chanLength: chanLength,
		tasks:      make(TaskChan[T], chanLength),
		wg:         &sync.WaitGroup{},
	}
}

func (w *ImportWorker[T]) Consume(ctx context.Context, workerFunc ConsumerFunc[T]) {
	go func() {
		w.wg.Wait()
		close(w.tasks)
	}()

	for {
		task, ok := <-w.tasks
		if !ok {
			return
		}

		workerFunc(ctx, task)
	}
}

func (w *ImportWorker[T]) Produce(ctx context.Context, generator ProducerFunc[T]) {
	w.wg.Add(1)

	go func() {
		defer w.wg.Done()

		generator(ctx, w.tasks)
	}()
}
