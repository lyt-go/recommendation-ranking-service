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

// Submit 派发一次排名任务到后台。任务上下文同时受请求上下文与
// 关闭上下文控制：任一取消都会中断在途调用并停止重试。因此请求
// 取消能尽快返回取消结果、重试不再增长，关闭也能及时完成。
//
// 返回的 errC 为缓冲通道（容量 1）：即使调用方在收到取消后不再
// 读取结果，worker goroutine 写入也不会阻塞，避免 goroutine 泄漏。
func (d *Dispatcher) Submit(requestCtx context.Context) <-chan error {
	errC := make(chan error, 1)
	// 任务上下文派生自请求上下文，并随关闭上下文取消而取消：
	// 请求取消 → 请求方主动断开；关闭 → Shutdown 取消 d.ctx，
	// cancelShut 使在途任务一并取消。
	taskCtx, taskCancel := context.WithCancel(requestCtx)
	cancelShut := context.AfterFunc(d.ctx, taskCancel)
	d.wg.Add(1)
	go func() {
		defer d.wg.Done()
		defer cancelShut()
		defer taskCancel()
		errC <- d.worker.Run(taskCtx)
	}()
	return errC
}

// Shutdown 取消关闭上下文，并等待所有在途任务退出。等待超过
// timeout 时返回 DeadlineExceeded。由于在途任务在关闭上下文取消
// 时随之取消，其调用与退避会被立刻中断，从而迅速结束。
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
