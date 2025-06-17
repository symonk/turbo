/*
package turbo exports the public API of the library.
This allows internals to change easily without breaking callers
code as they are solely integrating with turbo via this public
api.
*/
package turbo

import "github.com/symonk/turbo/internal/pool"

// WorkerPool manages a set of worker goroutines to execute tasks
// submitted concurrently
type WorkerPool = pool.WorkerPool

// NewPool instantiates a new pool with options
var NewPool = pool.NewPool

// Option
type PoolOption = pool.PoolOption

// WithAutoScaleDuration is a functional option for setting
// how frequently the worker pool will internally check for
// scalability windows.  That is, when workers can be shutdown
// for being idle.
var WithAutoScaleDuration = pool.WithAutoScaleDuration

// WithHooks is a functional option for setting hooks.  These
// are invoked throughout the worker pool lifecycle.
var WithHooks = pool.WithHooks
