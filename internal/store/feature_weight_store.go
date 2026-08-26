package store

import (
	"recommendation/internal/model"
)

func (s *MemoryStore) CreateFeatureWeight(f *model.FeatureWeight) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.featureWeights[f.ID] = f
	return nil
}

func (s *MemoryStore) GetFeatureWeight(id string) (*model.FeatureWeight, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	f, ok := s.featureWeights[id]
	if !ok {
		return nil, ErrNotFound
	}
	return f, nil
}

func (s *MemoryStore) ListFeatureWeights() []*model.FeatureWeight {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.FeatureWeight, 0, len(s.featureWeights))
	for _, f := range s.featureWeights {
		list = append(list, f)
	}
	return list
}

func (s *MemoryStore) UpdateFeatureWeight(f *model.FeatureWeight) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.featureWeights[f.ID]; !ok {
		return ErrNotFound
	}
	s.featureWeights[f.ID] = f
	return nil
}

func (s *MemoryStore) DeleteFeatureWeight(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.featureWeights[id]; !ok {
		return ErrNotFound
	}
	delete(s.featureWeights, id)
	return nil
}
