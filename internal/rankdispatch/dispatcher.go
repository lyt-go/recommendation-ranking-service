package rankdispatch

import "recommendation/internal/requestscope"

type Dispatcher struct {
	pool *requestscope.Pool
	gate <-chan struct{}
}

func New(pool *requestscope.Pool, gate <-chan struct{}) *Dispatcher {
	return &Dispatcher{pool: pool, gate: gate}
}
func (d *Dispatcher) Dispatch(requestID, userID string, done func(requestscope.Snapshot)) {
	s := d.pool.Acquire(requestID, userID, nil)
	d.pool.Release(s)
	go func() { <-d.gate; done(s.Snapshot()) }()
}
