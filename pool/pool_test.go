package pool

import (
	"fmt"
	"sync"
	"testing"
)

// debugger registers arbitrary hooks on the workerpool for
// debuggability aiding.  The pool itself can be highly concurrent
// and difficult to debug rare race/edge conditions.  Registering these
// debug hooks reports events to stdout that may be helpful occassionally.
type debugger struct{}

func (d debugger) OnWorkerStart(id int) {
	fmt.Println("worker started: ", id)
}
func (d debugger) OnWorkerStop(id int) {
	fmt.Println("worker stopped: ", id)
}

func (d debugger) OnPoolStop(graceful bool) {
	fmt.Println("pool was stopped: ", graceful)
}

func TestPoolCanBeClosedConcurrentlyWithoutIssue(t *testing.T) {
	maximumWorkers := 5
	p := New(maximumWorkers, WithHooks[any](debugger{}))
	var wg sync.WaitGroup
	wg.Add(5)
	go func(p *WorkerPool[any]) {
		for range 5 {
			p.Enqueue(func() { wg.Done() })
		}
	}(p)
	wg.Wait()
	p.Stop(true)
}
