package refreshjob

import (
	"context"
	"sync"
	"time"

	"recommendation/internal/featureclient"
)

type temporary interface{ Temporary() bool }

type Coordinator struct {
	client      *featureclient.Client
	maxAttempts int
	backoff     time.Duration
	ctxOnce     sync.Once
	ctx         context.Context
}

func New(client *featureclient.Client, maxAttempts int, backoff time.Duration) *Coordinator {
	return &Coordinator{client: client, maxAttempts: maxAttempts, backoff: backoff}
}

func (c *Coordinator) Refresh(ctx context.Context, userID string) error {
	c.ctxOnce.Do(func() { c.ctx = ctx })
	for attempt := 0; attempt < c.maxAttempts; attempt++ {
		err := c.client.Load(c.ctx, userID)
		if err == nil {
			return nil
		}
		t, ok := err.(temporary)
		if !ok || !t.Temporary() {
			return err
		}
		time.Sleep(c.backoff)
	}
	return context.DeadlineExceeded
}
