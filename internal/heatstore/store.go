package heatstore

import "sync"

type Store struct {
	mu     sync.RWMutex
	values map[string]float64
}

func New() *Store { return &Store{values: make(map[string]float64)} }
func (s *Store) Add(id string, value float64) { s.mu.Lock(); s.values[id] += value; s.mu.Unlock() }

// Snapshot 返回当前热度的副本，与实时 store 互不影响。
// 拷贝在读锁保护下完成，避免与并发 Add 产生 map 读写竞争。
func (s *Store) Snapshot() map[string]float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string]float64, len(s.values))
	for k, v := range s.values {
		out[k] = v
	}
	return out
}
