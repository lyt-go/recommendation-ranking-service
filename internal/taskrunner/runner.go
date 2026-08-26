package taskrunner

import (
	"fmt"
	"recommendation/internal/deliverytask"
	"recommendation/internal/taskcache"
	"sync"
)

type Effect func(string) error
type Runner struct {
	mu       sync.Mutex
	versions map[string]int
	cache    *taskcache.Cache
	effect   Effect
}
type Attempt struct {
	runner  *Runner
	id, key string
	version int
}

func New(cache *taskcache.Cache, effect Effect) *Runner {
	return &Runner{versions: make(map[string]int), cache: cache, effect: effect}
}
func (r *Runner) Begin(id string) *Attempt {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.versions[id]++
	v := r.versions[id]
	r.cache.Put(deliverytask.Task{ID: id, Version: v, Status: deliverytask.StatusRunning})
	return &Attempt{runner: r, id: id, key: fmt.Sprintf("delivery:%s:%d", id, v), version: v}
}
func (a *Attempt) Apply() error { return a.runner.effect(a.key) }
func (a *Attempt) Finish(status string) {
	a.runner.cache.Put(deliverytask.Task{ID: a.id, Version: a.version, Status: status})
}
