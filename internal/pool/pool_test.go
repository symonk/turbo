package pool

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

func TestPoolCanBeClosedMultipleTimesSafely(t *testing.T) {
	maximumWorkers := 5
	p := NewPool(maximumWorkers, WithHooks(debugger{}))
	var wg sync.WaitGroup
	wg.Add(5)
	go func(p *WorkerPool) {
		for range 5 {
			p.Enqueue(func() { wg.Done() })
		}
	}(p)
	wg.Wait()
	p.Stop(true)
	p.Stop(true)
}

func TestMinTasksIsPositive(t *testing.T) {
	p := NewPool(-1)
	assert.Equal(t, p.maxWorkers, 1)
}

func TestTasksCannotBeEnqueuedWhenClosed(t *testing.T) {
	p := NewPool(1)
	p.Stop(true)
	id, ok := p.Enqueue(func() {})
	assert.Empty(t, id)
	assert.False(t, ok)
}

func TestManyTasks(t *testing.T) {
	p := NewPool(runtime.NumCPU(), WithAutoScaleDuration(100*time.Millisecond), WithHooks(debugger{}))
	var wg sync.WaitGroup
	wg.Add(10_000)
	var done int32
	for range 10_000 {
		p.Enqueue(func() {
			time.Sleep(time.Microsecond)
			atomic.AddInt32(&done, 1)
			defer wg.Done()
		})
	}
	// Wait for all the work to be done, it should auto scale down.
	time.Sleep(3 * time.Second)
	p.Stop(true)
	fmt.Println(done)
	wg.Wait()
}
