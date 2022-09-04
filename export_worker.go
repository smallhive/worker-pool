package worker_pool

import (
	"context"
	"sync"
)

type ExportWorker[T any] struct {
	wg           *sync.WaitGroup
	workerAmount int
	tasks        TaskChan[T]
}

func NewExportWorker[T any](workerAmount int) *ExportWorker[T] {
	return &ExportWorker[T]{
		workerAmount: workerAmount,
		tasks:        make(TaskChan[T], workerAmount),
		wg:           &sync.WaitGroup{},
	}
}

func (w *ExportWorker[T]) Consume(ctx context.Context, consumerFunc ConsumerFunc[T]) {
	w.wg.Add(w.workerAmount)

	for i := 0; i < w.workerAmount; i++ {
		go func() {
			defer w.wg.Done()

			for {
				task, ok := <-w.tasks
				if !ok {
					return
				}

				consumerFunc(ctx, task)
			}
		}()
	}
}

// Produce generates tasks for worker. It is a sync call
func (w *ExportWorker[T]) Produce(ctx context.Context, generator ProducerFunc[T]) {
	generator(ctx, w.tasks)
	close(w.tasks)

	w.wg.Wait()
}
