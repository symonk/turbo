package pool

import "time"

// PoolOption is the signature of a functional option that configures
// the worker pool
type PoolOption func(*WorkerPool)

// WithAutoScaleDuration sets the frequency in which the pool should
// allow either new worker spawns or termination of existing workers
// based on the worker pool load.
func WithAutoScaleDuration(duration time.Duration) PoolOption {
	return func(w *WorkerPool) {
		w.workerIdleCheckFrequency = duration
	}
}

// WithHooks allows registering arbitrary callbacks for various stages
// in the worker pool life cycle.
func WithHooks(hooks Hooker) PoolOption {
	return func(w *WorkerPool) {
		w.hooks = hooks
	}
}
