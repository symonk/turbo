package pool

import "context"

// Pooler outsides the interface implemented by the pool
// TODO: Roll the Enqueue methods into single methods with a functional option for priority and blocking.
type Pooler[T any] interface {
	// [Task Submission]
	// Enqueue submits the task internally.
	Enqueue(task Task[T])
	// EnqueuePriority submits the task with the given priority
	EnqueuePriority(task Task[T], priority int)
	// EnqueueWait submits a task to the pool and blocks until it has been processed
	EnqueueWait(task Task[T])
	// EnqueueWaitPriority submits a task to the pool with a priority and blocks until it has been processed.
	EnqueueWaitPriority(task Task[T])

	// [Pool Control]
	// Stop the workerpool
	Stop(graceful bool) error
	// Pause the pool from processing any more tasks temporarily.
	Pause(ctx context.Context)
	// Flush blocks until the internal pool queues reach zero
	// or the timeout specified expires.
	Flush(ctx context.Context)
}
