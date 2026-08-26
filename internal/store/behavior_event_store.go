package store

import (
	"recommendation/internal/model"
)

func (s *MemoryStore) CreateBehaviorEvent(b *model.BehaviorEvent) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.behaviorEvents[b.ID] = b
	return nil
}

func (s *MemoryStore) GetBehaviorEvent(id string) (*model.BehaviorEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	b, ok := s.behaviorEvents[id]
	if !ok {
		return nil, ErrNotFound
	}
	return b, nil
}

func (s *MemoryStore) ListBehaviorEvents() []*model.BehaviorEvent {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.BehaviorEvent, 0, len(s.behaviorEvents))
	for _, b := range s.behaviorEvents {
		list = append(list, b)
	}
	return list
}

func (s *MemoryStore) DeleteBehaviorEvent(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.behaviorEvents[id]; !ok {
		return ErrNotFound
	}
	delete(s.behaviorEvents, id)
	return nil
}

func (s *MemoryStore) CreateBehaviorEvents(events []*model.BehaviorEvent) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, b := range events {
		s.behaviorEvents[b.ID] = b
	}
	return nil
}
