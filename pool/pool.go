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

type WorkerPool struct {
	maxWorkers               int
	stopped                  chan struct{}
	finalized                chan struct{}
	tasks                    chan Task
	workerQueue              chan Task
	stopper                  sync.Once
	workerIdleCheckFrequency time.Duration
	hooks                    Hooker
	closed                   atomic.Bool
}

// NewPool creates a new worker pool and returns it
func NewPool(maxWorkers int, options ...PoolOption) *WorkerPool {
	w := &WorkerPool{}
	w.workerIdleCheckFrequency = defaultWorkerIdleFrequency
	for _, option := range options {
		option(w)
	}
	w.maxWorkers = max(1, maxWorkers)
	w.stopped = make(chan struct{})
	w.finalized = make(chan struct{})
	w.tasks = make(chan Task)
	w.workerQueue = make(chan Task)
	go w.dispatch()
	return w
}

// dispatch is an asynchronous event loop processing tasks
// as they enter the worker pool queues.
func (w *WorkerPool) dispatch() {
	defer close(w.finalized)
	var nextWorkerId int
	var currentWorkers int
	var workerWg sync.WaitGroup

	// Setup a ticker for autoscaling (down) operations.
	autoScaler := time.NewTicker(w.workerIdleCheckFrequency)
	defer autoScaler.Stop()

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
				continue // avoid putting the task on twice.
			}
			w.workerQueue <- task
		case <-autoScaler.C:
			//
		}
	}
	// TODO: Consider how we gracefully exit and clean up if requested.
	for range currentWorkers {
		w.stopWorker()
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
func (w *WorkerPool) Stop(graceful bool) {
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
func (w *WorkerPool) Enqueue(t Task) (string, bool) {
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
// TODO: Consider an implementation that tracks 'idleness' within the worker
// if it gets at ask reset its idle timer, else auto shut itself down.
func (w *WorkerPool) startWorker(id int, t Task, wg *sync.WaitGroup) {
	wg.Add(1)
	go worker(t, wg, w.workerQueue, func() {
		if w.hooks != nil {
			w.hooks.OnWorkerStop(id)
		}
	})
	if w.hooks != nil {
		w.hooks.OnWorkerStart(id)
	}
}

// stopWorker terminates an existing worker goroutine
// by sending in a nil tasks, this causes the decoupled
// `worker` goroutine to pull it off and the channel and exit.
func (w *WorkerPool) stopWorker() {
	w.workerQueue <- nil
}

// worker is a generic worker function that receives tasks
// from the input channel.
//
// If a nil tasks is sent to a worker, the worker will
// terminate
func worker(t Task, wg *sync.WaitGroup, input <-chan Task, callback func()) {
	defer callback()
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
