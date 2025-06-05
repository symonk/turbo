package pool

import (
	"runtime"
	"testing"
	"time"
)

func BenchmarkBasic(b *testing.B) {
	// TODO: This is not actually realistic, used as a debugging tool for now.
	pool := NewPool(runtime.NumCPU())
	defer pool.Stop(true)

	// collect heap allocations
	b.ReportAllocs()
	// reset the benchmark to avoid overhead from instantiation
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			done := make(chan struct{})
			pool.Enqueue(func() {
				time.Sleep(100 * time.Millisecond)
				close(done)
			})
			<-done
		}
	})
}
