package store

import (
	"recommendation/internal/model"
)

func (s *MemoryStore) CreateUserProfile(u *model.UserProfile) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.userProfiles {
		if exist.UserID == u.UserID {
			return ErrConflict
		}
	}
	s.userProfiles[u.ID] = u
	return nil
}

func (s *MemoryStore) GetUserProfile(id string) (*model.UserProfile, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.userProfiles[id]
	if !ok {
		return nil, ErrNotFound
	}
	return u, nil
}

func (s *MemoryStore) GetUserProfileByUserID(userID string) (*model.UserProfile, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, u := range s.userProfiles {
		if u.UserID == userID {
			return u, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListUserProfiles() []*model.UserProfile {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.UserProfile, 0, len(s.userProfiles))
	for _, u := range s.userProfiles {
		list = append(list, u)
	}
	return list
}

func (s *MemoryStore) UpdateUserProfile(u *model.UserProfile) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.userProfiles[u.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.userProfiles {
		if exist.ID != u.ID && exist.UserID == u.UserID {
			return ErrConflict
		}
	}
	s.userProfiles[u.ID] = u
	return nil
}

func (s *MemoryStore) DeleteUserProfile(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.userProfiles[id]; !ok {
		return ErrNotFound
	}
	delete(s.userProfiles, id)
	return nil
}
