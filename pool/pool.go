/*
Package pool provides a fast, goroutine-based worker pool
for executing tasks concurrently with minimal overhead.
*/
package pool

type WorkerPool[T any] struct {
	maxWorkers int
}

// New creates a new worker pool and returns it
func New[T any](options ...PoolOption[T]) *WorkerPool[T] {
	w := &WorkerPool[T]{}
	for _, option := range options {
		option(w)
	}
	return w
}
