/*
Package pool provides a fast, goroutine-based worker pool
for executing tasks concurrently with minimal overhead.
*/
package pool

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

/*
TODO: Don't like the worker frequency check naming
TODO: Basic task submission uses a blocking channel, is this the approach we want?
TODO: Exposes some 'stats' similar to how runtime exposes stats
TODO: Think about the internal data structures, how tasks can move across
TODO: Think about a 'priority' capability with Enqueue()

*/

var (
	// defaultWorkerIdleFrequency is the frequency of the ticker for determining
	// if spawned workers are potentially idle and should receive a
	// shutdown signal
	defaultWorkerIdleFrequency = 10 * time.Second
)

var (
	// ErrWorkerPoolClosed is returned by Enqueue() variants when the workerpool has
	// been terminated and is (potentially) in a teardown state.  Submitting new tasks
	// is not permitted, the pool is likely waiting for existing tasks to finalize
	// before terminating.
	ErrWorkerPoolClosed = errors.New("cannot enqueue tasks after pool was stopped")
)

type WorkerPool[T any] struct {
	maxWorkers               int
	stopped                  chan struct{}
	finalized                chan struct{}
	tasks                    chan Task[T]
	workerQueue              chan Task[T]
	stopper                  sync.Once
	workerIdleCheckFrequency time.Duration
	hooks                    Hooker
	closed                   atomic.Bool
}

// New creates a new worker pool and returns it
func New[T any](maxWorkers int, options ...PoolOption[T]) *WorkerPool[T] {
	w := &WorkerPool[T]{}
	w.workerIdleCheckFrequency = defaultWorkerIdleFrequency
	for _, option := range options {
		option(w)
	}
	w.maxWorkers = max(1, maxWorkers)
	w.stopped = make(chan struct{})
	w.finalized = make(chan struct{})
	w.tasks = make(chan Task[T])
	w.workerQueue = make(chan Task[T])
	go w.dispatch()
	return w
}

// dispatch is an asynchronous event loop processing tasks
// as they enter the worker pool queues.
func (w *WorkerPool[T]) dispatch() {
	defer close(w.finalized)
	var nextWorkerId int
	var currentWorkers int
	var workerWg sync.WaitGroup

	idleChecker := time.NewTicker(w.workerIdleCheckFrequency)

eventloop:
	for {
		select {
		case <-w.stopped:
			// Client has requested a shutdown.
			w.closed.Store(true)
			break eventloop
		case task, ok := <-w.tasks:
			if !ok {
				// the pool has been stopped, do not enqueue any more tasks.
				break eventloop
			}
			if currentWorkers < w.maxWorkers {
				nextWorkerId++
				w.startWorker(nextWorkerId, task, &workerWg)
				currentWorkers++
			}
			w.workerQueue <- task
		case <-idleChecker.C:
			// TODO: Determine if workers are idle and require a potential shutdown
		}
	}
	// TODO: Consider how we gracefully exit and clean up if requested.
	for range currentWorkers {
		w.stopWorker(&workerWg)
	}

	workerWg.Wait()
}

// Stop shuts down the worker pool, optionally waiting for all tasks
// in flight to be finished.  After Stop() has been called, enqueing
// new tasks into the pool will return an error.
//
// Stop is thread safe and is a blocking call until the worker pool
// has completely ceased.  Subsequent calls to the workerpool after
// Stop has been called will exit immediately.
func (w *WorkerPool[T]) Stop(graceful bool) {
	w.stopper.Do(func() {
		close(w.stopped)
	})

	<-w.finalized
	if w.hooks != nil {
		w.hooks.OnPoolStop(graceful)
	}
}

// Enqueue submits a task onto the worker pool
//
// nil tasks are not allowed from callers as they are used internally
// by the pool to signal a worker shutdown.
func (w *WorkerPool[T]) Enqueue(t Task[T]) (string, bool) {
	id := uuid.New().String()
	if t != nil {
		select {
		case <-w.stopped:
			return "", false
		default:
			w.tasks <- t
			return id, true
		}
	}
	return "", false
}

// startWorker starts a new worker goroutine
func (w *WorkerPool[T]) startWorker(id int, t Task[T], wg *sync.WaitGroup) {
	wg.Add(1)
	go worker(t, wg, w.workerQueue)
	if w.hooks != nil {
		w.hooks.OnWorkerStart(id)
	}
}

// stopWorker terminates an existing worker goroutine
func (w *WorkerPool[T]) stopWorker(wg *sync.WaitGroup) {
	w.workerQueue <- nil
}

// worker is a generic worker function that receives tasks
// from the input channel.
//
// If a nil tasks is sent to a worker, the worker will
// terminate
func worker[T any](t Task[T], wg *sync.WaitGroup, input <-chan Task[T]) {
	defer wg.Done()
	for t != nil {
		t := <-input
		if t != nil {
			t()
			continue
		}
		break
	}
}
