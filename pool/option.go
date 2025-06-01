package pool

// PoolOption is the signature of a functional option that configures
// the worker pool
type PoolOption[T any] func(*WorkerPool[T])

func WithMaxWorkers[T any](max int) PoolOption[T] {
	return func(w *WorkerPool[T]) {
		w.maxWorkers = max
	}
}
