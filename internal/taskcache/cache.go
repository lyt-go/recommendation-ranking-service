package taskcache

import (
	"recommendation/internal/deliverytask"
	"sync"
)

type Cache struct {
	mu    sync.RWMutex
	tasks map[string]deliverytask.Task
}

func New() *Cache                        { return &Cache{tasks: make(map[string]deliverytask.Task)} }
func (c *Cache) Put(t deliverytask.Task) { c.mu.Lock(); c.tasks[t.ID] = t; c.mu.Unlock() }
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
