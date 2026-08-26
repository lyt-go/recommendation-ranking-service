package rankdispatch_test

import (
	"recommendation/internal/rankdispatch"
	"recommendation/internal/requestscope"
	"testing"
	"time"
)

func TestAsyncRankingKeepsRequestIdentityAfterPoolReuse(t *testing.T) {
	pool := requestscope.NewPool()
	gate := make(chan struct{})
	dispatcher := rankdispatch.New(pool, gate)
	results := make(chan requestscope.Snapshot, 2)
	dispatcher.Dispatch("req-a", "user-a", func(s requestscope.Snapshot) { results <- s })
	dispatcher.Dispatch("req-b", "user-b", func(s requestscope.Snapshot) { results <- s })
	close(gate)
	seen := map[string]string{}
	for i := 0; i < 2; i++ {
		select {
		case s := <-results:
			seen[s.RequestID] = s.UserID
		case <-time.After(time.Second):
			t.Fatal("ranking callback did not finish")
		}
	}
	if seen["req-a"] != "user-a" || seen["req-b"] != "user-b" {
		t.Errorf("async identities = %v, want req-a/user-a and req-b/user-b", seen)
	}
	s := pool.Acquire("req-c", "user-c", nil)
	pool.Release(s)
	clean := pool.Acquire("", "", nil)
	if clean.RequestID != "" || clean.UserID != "" || len(clean.Tags) != 0 {
		t.Errorf("reused request scope retained identity: %+v", clean)
	}
}
