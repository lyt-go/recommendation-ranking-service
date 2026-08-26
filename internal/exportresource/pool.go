package exportresource

import (
	"errors"
	"sync"
)

var ErrExhausted = errors.New("export handles exhausted")

type Pool struct {
	mu            sync.Mutex
	limit, active int
}
type Lease struct {
	pool *Pool
	once sync.Once
}

func NewPool(limit int) *Pool { return &Pool{limit: limit} }
func (p *Pool) Acquire() (*Lease, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.active >= p.limit {
		return nil, ErrExhausted
	}
	p.active++
	return &Lease{pool: p}, nil
}
func (l *Lease) Release() {
	l.once.Do(func() { l.pool.mu.Lock(); l.pool.active--; l.pool.mu.Unlock() })
}
func (p *Pool) Active() int { p.mu.Lock(); defer p.mu.Unlock(); return p.active }
