package plancache

import (
	"recommendation/internal/planbuilder"
	"sync"
)

type Cache struct {
	mu    sync.RWMutex
	plans map[string]*planbuilder.Plan
}

func New() *Cache                        { return &Cache{plans: make(map[string]*planbuilder.Plan)} }
func (c *Cache) Put(p *planbuilder.Plan) { c.mu.Lock(); c.plans[p.Key] = p; c.mu.Unlock() }
func (c *Cache) Get(k string) (*planbuilder.Plan, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	p, ok := c.plans[k]
	return p, ok
}
