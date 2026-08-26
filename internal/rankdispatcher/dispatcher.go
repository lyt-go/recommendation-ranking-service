package rankdispatcher

import (
	"context"
	"recommendation/internal/rankworker"
	"sync"
	"time"
)

type Dispatcher struct {
	ctx    context.Context
	cancel context.CancelFunc
	worker *rankworker.Worker
	wg     sync.WaitGroup
}

func New(worker *rankworker.Worker) *Dispatcher {
	ctx, cancel := context.WithCancel(context.Background())
	return &Dispatcher{ctx: ctx, cancel: cancel, worker: worker}
}
func (d *Dispatcher) Submit(requestCtx context.Context) <-chan error {
	result := make(chan error)
	d.wg.Add(1)
	go func() { defer d.wg.Done(); result <- d.worker.Run(context.Background()); close(result) }()
	return result
}
func (d *Dispatcher) Shutdown(timeout time.Duration) error {
	d.cancel()
	done := make(chan struct{})
	go func() { d.wg.Wait(); close(done) }()
	select {
	case <-done:
		return nil
	case <-time.After(timeout):
		return context.DeadlineExceeded
	}
}
