package pool

// PoolOption is the signature of a functional option that configures
// the worker pool
type PoolOption[T any] func(*WorkerPool[T])

// WithAutoScaleDuration sets the frequency in which the pool should
// allow either new worker spawns or termination of existing workers
// based on the worker pool load.
func WithAutoScaleDuration[T any](max int) PoolOption[T] {
	return func(w *WorkerPool[T]) {
		panic("TODO: not implemented")
	}
}
