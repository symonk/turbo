package pool

// PoolOption is the signature of a functional option that configures
// the worker pool
type PoolOption[T any] func(*WorkerPool[T])

// TODO: Functional option for auto scaling frequency

// WithAutoScaleDuration sets the frequency in which the pool should
// allow either new worker spawns or termination of existing workers
// based on the worker pool load.
func WithAutoScaleDuration[T any](max int) PoolOption[T] {
	return func(w *WorkerPool[T]) {
		panic("TODO: not implemented")
	}
}

// WithHooks allows registering arbitrary callbacks for various stages
// in the worker pool life cycle.
func WithHooks[T any](hooks Hooker) PoolOption[T] {
	return func(w *WorkerPool[T]) {
		w.hooks = hooks
	}
}
