/*
*
package main exports a public API that clients should use.
*
*/
package main

import "github.com/symonk/turbo/pool"

// WorkerPool manages a set of worker goroutines to execute tasks
// submitted concurrently
type WorkerPool = pool.WorkerPool
