package rankworker

import (
	"context"
	"errors"
	"time"
)

type Backend interface{ Call(context.Context) error }
type Worker struct {
	backend Backend
	backoff  time.Duration
}

func New(backend Backend, backoff time.Duration) *Worker {
	return &Worker{backend: backend, backoff: backoff}
}

// Run 反复调用排名服务直到成功、返回非临时错误或上下文取消。
// 请求/关闭取消会经 ctx 传入：Call 使用该上下文，退避也通过 select
// 响应取消，确保取消后立刻停止重试、不再继续发起请求。
func (w *Worker) Run(ctx context.Context) error {
	for {
		// 进入本轮前上下文已取消，直接返回，避免无谓的一次调用。
		if err := ctx.Err(); err != nil {
			return err
		}
		err := w.backend.Call(ctx)
		if err == nil {
			return nil
		}
		// 调用返回错误时若上下文已取消，优先返回取消结果。
		if err := ctx.Err(); err != nil {
			return err
		}
		var temporary interface{ Temporary() bool }
		if !errors.As(err, &temporary) || !temporary.Temporary() {
			return err
		}
		// 退避期间同样响应取消，否则取消后还会再发一轮请求。
		timer := time.NewTimer(w.backoff)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}
}
