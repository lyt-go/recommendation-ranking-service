package refreshjob_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"recommendation/internal/featureclient"
	"recommendation/internal/refreshjob"
)

type backend struct {
	mu      sync.Mutex
	calls   []string
	started chan struct{}
	release chan struct{}
}

func (b *backend) Load(ctx context.Context, userID string) error {
	b.mu.Lock()
	b.calls = append(b.calls, userID)
	b.mu.Unlock()
	if userID == "first" {
		select {
		case b.started <- struct{}{}:
		default:
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-b.release:
			return nil
		}
	}
	return nil
}

func TestCancelledRefreshStopsAndNextRequestIsIndependent(t *testing.T) {
	b := &backend{started: make(chan struct{}, 1), release: make(chan struct{})}
	c := refreshjob.New(featureclient.New(b), 3, time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- c.Refresh(ctx, "first") }()
	<-b.started
	cancel()
	select {
	case err := <-done:
		if err != context.Canceled {
			t.Errorf("cancelled refresh error = %v, want context.Canceled", err)
		}
	case <-time.After(50 * time.Millisecond):
		t.Errorf("cancelled refresh kept the feature request running")
		close(b.release)
		<-done
	}

	if err := c.Refresh(context.Background(), "second"); err != nil {
		t.Fatalf("next refresh inherited the cancelled request: %v", err)
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if len(b.calls) != 2 || b.calls[1] != "second" {
		t.Fatalf("feature calls = %v, want independent first and second requests", b.calls)
	}
}
