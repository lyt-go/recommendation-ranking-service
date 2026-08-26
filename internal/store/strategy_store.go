package store

import (
	"recommendation/internal/model"
)

func (s *MemoryStore) CreateStrategy(st *model.Strategy) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.strategies[st.ID] = st
	return nil
}

func (s *MemoryStore) GetStrategy(id string) (*model.Strategy, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	st, ok := s.strategies[id]
	if !ok {
		return nil, ErrNotFound
	}
	return st, nil
}

func (s *MemoryStore) ListStrategies() []*model.Strategy {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Strategy, 0, len(s.strategies))
	for _, st := range s.strategies {
		list = append(list, st)
	}
	return list
}

func (s *MemoryStore) UpdateStrategy(st *model.Strategy) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.strategies[st.ID]; !ok {
		return ErrNotFound
	}
	s.strategies[st.ID] = st
	return nil
}

func (s *MemoryStore) DeleteStrategy(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.strategies[id]; !ok {
		return ErrNotFound
	}
	delete(s.strategies, id)
	return nil
}
