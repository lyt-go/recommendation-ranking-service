package store

import (
	"recommendation/internal/model"
)

func (s *MemoryStore) CreateRecommendation(r *model.Recommendation) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.recommendations[r.ID] = r
	return nil
}

func (s *MemoryStore) GetRecommendation(id string) (*model.Recommendation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.recommendations[id]
	if !ok {
		return nil, ErrNotFound
	}
	return r, nil
}

func (s *MemoryStore) ListRecommendations() []*model.Recommendation {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Recommendation, 0, len(s.recommendations))
	for _, r := range s.recommendations {
		list = append(list, r)
	}
	return list
}

func (s *MemoryStore) UpdateRecommendation(r *model.Recommendation) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.recommendations[r.ID]; !ok {
		return ErrNotFound
	}
	s.recommendations[r.ID] = r
	return nil
}

func (s *MemoryStore) DeleteRecommendation(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.recommendations[id]; !ok {
		return ErrNotFound
	}
	delete(s.recommendations, id)
	return nil
}
