package heatcache

import "sync"

type Cache struct {
	mu     sync.RWMutex
	values map[string]float64
}

func (c *Cache) Store(values map[string]float64) { c.mu.Lock(); c.values = values; c.mu.Unlock() }
func (c *Cache) Get() map[string]float64         { c.mu.RLock(); defer c.mu.RUnlock(); return c.values }
