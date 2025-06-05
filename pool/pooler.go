package pool

import (
	"context"
)

// Pooler outsides the interface implemented by the pool
// TODO: Roll the Enqueue methods into single methods with a functional option for priority and blocking.
type Pooler[T any] interface {
	// [Task Submission]
	// Enqueue submits the task internally.
	Enqueue(task Task[T]) (id string, ok bool)

	// [Pool Control]
	// Stop the workerpool
	Stop(graceful bool) error
	// Pause the pool from processing any more tasks temporarily.
	Pause(ctx context.Context)
	// Flush blocks until the internal pool queues reach zero
	// or the timeout specified expires.
	Flush(ctx context.Context)
}
