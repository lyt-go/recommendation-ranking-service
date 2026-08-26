package taskrunner_test

import (
	"recommendation/internal/deliverytask"
	"recommendation/internal/taskcache"
	"recommendation/internal/taskrunner"
	"sync"
	"testing"
)

func TestRetryKeepsOneEffectAndRejectsLateAttemptState(t *testing.T) {
	var mu sync.Mutex
	applied := map[string]bool{}
	effects := 0
	effect := func(key string) error {
		mu.Lock()
		defer mu.Unlock()
		if !applied[key] {
			applied[key] = true
			effects++
		}
		return nil
	}
	cache := taskcache.New()
	runner := taskrunner.New(cache, effect)
	first := runner.Begin("user-feed")
	if err := first.Apply(); err != nil {
		t.Fatalf("first delivery effect returned %v", err)
	}
	retry := runner.Begin("user-feed")
	if err := retry.Apply(); err != nil {
		t.Fatalf("retry delivery effect returned %v", err)
	}
	retry.Finish(deliverytask.StatusSucceeded)
	first.Finish(deliverytask.StatusRunning)
	if effects != 1 {
		t.Errorf("external delivery effects = %d, want 1", effects)
	}
	detail, ok := cache.Get("user-feed")
	if !ok || detail.Status != deliverytask.StatusSucceeded || detail.Version != 2 {
		t.Errorf("detail task = %+v ok=%v, want version 2 succeeded", detail, ok)
	}
	list := cache.List()
	if len(list) != 1 || list[0].Status != deliverytask.StatusSucceeded {
		t.Fatalf("task list = %+v, want one succeeded task", list)
	}
}
