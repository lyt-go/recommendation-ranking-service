package heatview_test

import (
	"recommendation/internal/heatcache"
	"recommendation/internal/heatstore"
	"recommendation/internal/heatview"
	"sync"
	"testing"
)

func TestCapturedHeatSnapshotIsIsolatedFromWritersAndReaders(t *testing.T) {
	store := heatstore.New()
	cache := &heatcache.Cache{}
	view := heatview.New(store, cache)
	store.Add("item-a", 1)
	view.Capture()
	captured := cache.Get()
	store.Add("item-a", 2)
	if captured["item-a"] != 1 {
		t.Errorf("captured heat changed to %v, want 1", captured["item-a"])
	}
	captured["item-b"] = 99
	if next := cache.Get(); next["item-b"] != 0 {
		t.Errorf("cache was changed through reader snapshot: %v", next)
	}
	shared := store.Snapshot()
	start := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		<-start
		for i := 0; i < 2000; i++ {
			store.Add("item-a", 1)
		}
	}()
	go func() {
		defer wg.Done()
		<-start
		sum := float64(0)
		for i := 0; i < 2000; i++ {
			sum += shared["item-a"]
		}
		_ = sum
	}()
	close(start)
	wg.Wait()
}
