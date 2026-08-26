package refreshjob

import (
	"context"
	"time"

	"recommendation/internal/featureclient"
)

type temporary interface{ Temporary() bool }

type Coordinator struct {
	client      *featureclient.Client
	maxAttempts int
	backoff     time.Duration
}

func New(client *featureclient.Client, maxAttempts int, backoff time.Duration) *Coordinator {
	return &Coordinator{client: client, maxAttempts: maxAttempts, backoff: backoff}
}

// Refresh 拉取特征源刷新单个用户画像。ctx 由调用方提供，每次调用各自
// 拥有独立的生命周期：调用方取消后会立刻中断当前的重试/退避，且不会
// 影响后续任何一次刷新（不再复用上一次的 ctx）。
func (c *Coordinator) Refresh(ctx context.Context, userID string) error {
	for attempt := 0; attempt < c.maxAttempts; attempt++ {
		// 每一轮都把调用方的 ctx 透传给特征源，保证取消能即时传播；
		// 失败后若 ctx 已取消，立即退出，绝不挂起。
		err := c.client.Load(ctx, userID)
		if err == nil {
			return nil
		}
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}
		t, ok := err.(temporary)
		if !ok || !t.Temporary() {
			return err
		}
		// 退避期间同样监听 ctx，取消可立即打断等待。
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(c.backoff):
		}
	}
	return context.DeadlineExceeded
}
