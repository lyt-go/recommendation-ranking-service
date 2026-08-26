package heatcache

import "sync"

type Cache struct {
	mu     sync.RWMutex
	values map[string]float64
}

func New() *Cache { return &Cache{values: make(map[string]float64)} }

// Store 用传入 map 的副本替换缓存，避免调用方后续修改牵连缓存。
func (c *Cache) Store(values map[string]float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.values = make(map[string]float64, len(values))
	for k, v := range values {
		c.values[k] = v
	}
}

// Get 返回缓存的副本，调用方对返回值的修改不会污染缓存。
func (c *Cache) Get() map[string]float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make(map[string]float64, len(c.values))
	for k, v := range c.values {
		out[k] = v
	}
	return out
}
