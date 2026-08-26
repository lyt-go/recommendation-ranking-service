package service

import (
	"sort"
	"time"

	"recommendation/internal/model"
	"recommendation/pkg/idgen"
)

func (s *Service) CreateUserProfile(input model.UserProfile) (*model.UserProfile, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	now := time.Now()
	u := &model.UserProfile{
		ID:        idgen.Hex(),
		UserID:    input.UserID,
		Interests: input.Interests,
		Tags:      input.Tags,
		Region:    input.Region,
		UpdatedAt: now,
	}
	if err := s.store.CreateUserProfile(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Service) GetUserProfile(id string) (*model.UserProfile, error) {
	return s.store.GetUserProfile(id)
}

func (s *Service) GetUserProfileByUserID(userID string) (*model.UserProfile, error) {
	return s.store.GetUserProfileByUserID(userID)
}

func (s *Service) ListUserProfiles(filter model.UserProfileFilter, page, size int) ([]*model.UserProfile, int, error) {
	all := s.store.ListUserProfiles()
	matched := make([]*model.UserProfile, 0, len(all))
	for _, u := range all {
		if filter.Match(u) {
			matched = append(matched, u)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].UpdatedAt.After(matched[j].UpdatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.UserProfile{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateUserProfile(id string, input model.UserProfile) (*model.UserProfile, error) {
	u, err := s.store.GetUserProfile(id)
	if err != nil {
		return nil, err
	}
	if input.UserID != "" {
		u.UserID = input.UserID
	}
	if input.Interests != nil {
		u.Interests = input.Interests
	}
	if input.Tags != nil {
		u.Tags = input.Tags
	}
	if input.Region != "" {
		u.Region = input.Region
	}
	u.UpdatedAt = time.Now()
	if err := u.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateUserProfile(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Service) DeleteUserProfile(id string) error {
	return s.store.DeleteUserProfile(id)
}
