package taskcache

import (
	"recommendation/internal/deliverytask"
	"sync"
)

type Cache struct {
	mu    sync.RWMutex
	tasks map[string]deliverytask.Task
}

func New() *Cache { return &Cache{tasks: make(map[string]deliverytask.Task)} }

// Put 只接受向前推进的更新：必须版本号更大，或版本号相同但状态次序更高。
// 这样迟到的旧回调（版本更小，或同为 running 但新回调已 succeeded）无法
// 覆盖更新的成功状态。
func (c *Cache) Put(t deliverytask.Task) {
	c.mu.Lock()
	defer c.mu.Unlock()
	prev, ok := c.tasks[t.ID]
	if !ok {
		c.tasks[t.ID] = t
		return
	}
	if t.Version > prev.Version {
		c.tasks[t.ID] = t
		return
	}
	if t.Version == prev.Version && deliverytask.Rank(t.Status) > deliverytask.Rank(prev.Status) {
		c.tasks[t.ID] = t
		return
	}
	// 版本更小，或版本相同但状态更旧：丢弃这条迟到/回退的回调。
}
func (c *Cache) Get(id string) (deliverytask.Task, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	t, ok := c.tasks[id]
	return t, ok
}
func (c *Cache) List() []deliverytask.Task {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]deliverytask.Task, 0, len(c.tasks))
	for _, t := range c.tasks {
		out = append(out, t)
	}
	return out
}
