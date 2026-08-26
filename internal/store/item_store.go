package store

import (
	"recommendation/internal/model"
)

func (s *MemoryStore) CreateItem(i *model.Item) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[i.ID] = i
	return nil
}

func (s *MemoryStore) GetItem(id string) (*model.Item, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	i, ok := s.items[id]
	if !ok {
		return nil, ErrNotFound
	}
	return i, nil
}

func (s *MemoryStore) ListItems() []*model.Item {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Item, 0, len(s.items))
	for _, i := range s.items {
		list = append(list, i)
	}
	return list
}

func (s *MemoryStore) UpdateItem(i *model.Item) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[i.ID]; !ok {
		return ErrNotFound
	}
	s.items[i.ID] = i
	return nil
}

func (s *MemoryStore) DeleteItem(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[id]; !ok {
		return ErrNotFound
	}
	delete(s.items, id)
	return nil
}
