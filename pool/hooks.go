package pool

// Hooker defines the interface for callbacks that are fired
// on various stages of the worker pool lifecycle.
type Hooker interface {
	OnPoolStop(graceful bool)
	OnWorkerStart(id int)
	OnWorkerStop(id int)
}
