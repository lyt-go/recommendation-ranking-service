package rankdispatcher_test

import (
	"context"
	"recommendation/internal/rankdispatcher"
	"recommendation/internal/rankworker"
	"sync/atomic"
	"testing"
	"time"
)

type temporaryError struct{}

func (temporaryError) Error() string   { return "rank backend busy" }
func (temporaryError) Temporary() bool { return true }

type backend struct {
	calls   atomic.Int64
	started chan struct{}
}

func (b *backend) Call(ctx context.Context) error {
	b.calls.Add(1)
	select {
	case b.started <- struct{}{}:
	default:
	}
	return temporaryError{}
}
func TestRankingCancellationStopsRetriesAndAllowsPromptShutdown(t *testing.T) {
	b := &backend{started: make(chan struct{}, 1)}
	dispatcher := rankdispatcher.New(rankworker.New(b, 2*time.Millisecond))
	ctx, cancel := context.WithCancel(context.Background())
	result := dispatcher.Submit(ctx)
	<-b.started
	cancel()
	select {
	case err := <-result:
		if err != context.Canceled {
			t.Errorf("cancelled ranking error=%v, want context.Canceled", err)
		}
	case <-time.After(60 * time.Millisecond):
		t.Errorf("cancelled ranking did not return")
	}
	afterCancel := b.calls.Load()
	time.Sleep(15 * time.Millisecond)
	if now := b.calls.Load(); now != afterCancel {
		t.Errorf("backend calls grew after cancellation: %d -> %d", afterCancel, now)
	}
	if err := dispatcher.Shutdown(40 * time.Millisecond); err != nil {
		t.Fatalf("dispatcher shutdown returned %v", err)
	}
}
