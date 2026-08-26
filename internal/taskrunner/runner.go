package taskrunner

import (
	"fmt"
	"recommendation/internal/deliverytask"
	"recommendation/internal/taskcache"
	"sync"
)

// Effect 执行外部投递动作，传入的是幂等键。
type Effect func(string) error

type Runner struct {
	mu       sync.Mutex
	versions map[string]int
	cache    *taskcache.Cache
	effect   Effect
}

// Attempt 代表对某个任务的一次投递尝试。
// deliveryKey 是任务级稳定的幂等键：同一任务的多次重试使用同一个键，
// 这样外部系统可以据此去重；version 只作为内部状态版本，不参与幂等键。
type Attempt struct {
	runner      *Runner
	id          string
	deliveryKey string
	version     int
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
	return &Attempt{runner: r, id: id, deliveryKey: fmt.Sprintf("delivery:%s", id), version: v}
}

// Apply 触发外部投递。同一任务的每次重试都用相同的 deliveryKey，
// 因此即便重试执行了多次，外部系统也能按该键去重，只投递一次。
func (a *Attempt) Apply() error { return a.runner.effect(a.deliveryKey) }

// Finish 以给定状态结束本次尝试。cache.Put 会按版本号/状态次序过滤，
// 确保迟到的旧回调无法覆盖更新的成功状态。
func (a *Attempt) Finish(status string) {
	a.runner.cache.Put(deliverytask.Task{ID: a.id, Version: a.version, Status: status})
}
