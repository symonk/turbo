package main

import "github.com/symonk/turbo/pool"

// WorkerPool manages a set of worker goroutines to execute tasks
// submitted concurrently
type WorkerPool[T any] = pool.WorkerPool[T]
