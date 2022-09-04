# worker-pool

## Motivations

Occasionally I need to process data with goroutines. One time I need one source and few workers to process data,
another time I need many sources and one worker to consolidate data in one place. Each time I have to do this from
scratch. Eventually decided to create separated package for these tasks

## Go versions

`worker-pool` requires Go version 1.18+

## Installation

```shell
go get github.com/smallhive/worker-pool
```

## One source and few workers

> This worker has `ExportWorker` name, because it exports tasks to workers (in my mind)

`ExportWorker` has one producer and few consumers which take tasks from common task pool to process

```go
package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/smallhive/worker-pool"
)

var (
	total        int64
	workerAmount = 10
)

func main() {
	ctx := context.Background()
	mu := sync.Mutex{}

	// use go1.18 generics to define type on tasks
	pool := worker_pool.NewExportWorker[int64](workerAmount)

	// consumers do all work
	consumer := func(_ context.Context, task int64) {
		// mutex because we have `workerAmount` goroutines which write to total value
		mu.Lock()
		total += task
		mu.Unlock()
	}

	// run `workerAmount` consumers
	pool.Consume(ctx, consumer)

	// producer generates the tasks
	producer := func(_ context.Context, taskInput worker_pool.TaskChan[int64]) {
		for i := int64(0); i < 100; i++ {
			taskInput <- i
		}
	}

	// sync call, we lock here. After all task generated, it will wait until all consumers complete their jobs
	pool.Produce(ctx, producer)
	fmt.Println(total) // 4950
}
```

## Few sources (example with one) and one worker

> This worker has `ImportWorker` name, because it imports tasks from many sources (in my mind)

`ImportWorker` has few producers which put tasks in one pool and one consumer which read them

```go
package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/smallhive/worker-pool"
)

var (
	total          int64
	chanBufferSize = 2
)

func main() {
	ctx := context.Background()

	// use go1.18 generics to define type on tasks
	// `chanBufferSize` describes buffer size for task channel.
	pool := worker_pool.NewImportWorker[int64](chanBufferSize)

	// create producer which generates tasks
	producer := func(_ context.Context, taskInput worker_pool.TaskChan[int64]) {
		for i := int64(0); i < 100; i++ {
			taskInput <- i
		}
	}

	// run producer which already writes data in tasks chanel
	pool.Produce(ctx, producer)

	// consumer consolidates all tasks data in one place
	consumer := func(ctx context.Context, task int64) {
		total += task
	}

	// sync call, we lock here. Wait until all producers finish their work
	pool.Consume(ctx, consumer)
	fmt.Println(total) // 4950
}
```

## Few sources (example with few) and one worker

> The example similar to previous one except we have two producers

```go
package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/smallhive/worker-pool"
)

var (
	total          int64
	chanBufferSize = 2
)

func main() {
	ctx := context.Background()

	// use go1.18 generics to define type on tasks
	// `chanBufferSize` describes buffer size for task channel.
	pool := worker_pool.NewImportWorker[int64](chanBufferSize)

	// run producer one
	pool.Produce(ctx, func(_ context.Context, taskInput worker_pool.TaskChan[int64]) {
		for i := int64(0); i < 50; i++ {
			taskInput <- i
		}
	})

	// run producer two
	pool.Produce(ctx, func(_ context.Context, taskInput worker_pool.TaskChan[int64]) {
		for i := int64(50); i < 100; i++ {
			taskInput <- i
		}
	})

	// consumer consolidates all tasks data in one place
	consumer := func(ctx context.Context, task int64) {
		total += task
	}

	// sync call, we lock here. Wait until all producers finish their work
	pool.Consume(ctx, consumer)
	fmt.Println(total) // 4950
}
```

## Tests

```shell
make test
```