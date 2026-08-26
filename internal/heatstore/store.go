package heatstore

import "sync"

type Store struct {
	mu     sync.RWMutex
	values map[string]float64
}

func New() *Store                             { return &Store{values: make(map[string]float64)} }
func (s *Store) Add(id string, value float64) { s.mu.Lock(); s.values[id] += value; s.mu.Unlock() }
func (s *Store) Snapshot() map[string]float64 { s.mu.RLock(); defer s.mu.RUnlock(); return s.values }
