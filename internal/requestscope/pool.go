package requestscope

import "sync"

type Scope struct {
	RequestID string
	UserID    string
	Tags      []string
}
type Snapshot struct {
	RequestID string
	UserID    string
	Tags      []string
}
type Pool struct{ pool sync.Pool }

func NewPool() *Pool { p := &Pool{}; p.pool.New = func() any { return &Scope{} }; return p }
func (p *Pool) Acquire(requestID, userID string, tags []string) *Scope {
	s := p.pool.Get().(*Scope)
	s.RequestID = requestID
	s.UserID = userID
	s.Tags = append(s.Tags, tags...)
	return s
}
func (p *Pool) Release(s *Scope) { p.pool.Put(s) }
func (s *Scope) Snapshot() Snapshot {
	return Snapshot{RequestID: s.RequestID, UserID: s.UserID, Tags: s.Tags}
}
