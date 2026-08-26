package service

import (
	"sort"
	"time"

	"recommendation/internal/model"
	"recommendation/pkg/idgen"
)

func (s *Service) CreateItem(input model.Item) (*model.Item, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	now := time.Now()
	i := &model.Item{
		ID:        idgen.Hex(),
		Title:     input.Title,
		Category:  input.Category,
		Tags:      input.Tags,
		Score:     input.Score,
		Status:    input.Status,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.store.CreateItem(i); err != nil {
		return nil, err
	}
	return i, nil
}

func (s *Service) GetItem(id string) (*model.Item, error) {
	return s.store.GetItem(id)
}

func (s *Service) ListItems(filter model.ItemFilter, page, size int) ([]*model.Item, int, error) {
	all := s.store.ListItems()
	matched := make([]*model.Item, 0, len(all))
	for _, i := range all {
		if filter.Match(i) {
			matched = append(matched, i)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Item{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateItem(id string, input model.Item) (*model.Item, error) {
	i, err := s.store.GetItem(id)
	if err != nil {
		return nil, err
	}
	if input.Title != "" {
		i.Title = input.Title
	}
	if input.Category != "" {
		i.Category = input.Category
	}
	if input.Tags != nil {
		i.Tags = input.Tags
	}
	if input.Score != 0 {
		i.Score = input.Score
	}
	i.UpdatedAt = time.Now()
	if err := i.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateItem(i); err != nil {
		return nil, err
	}
	return i, nil
}

func (s *Service) UpdateItemStatus(id string, newStatus string) (*model.Item, error) {
	i, err := s.store.GetItem(id)
	if err != nil {
		return nil, err
	}
	if !model.CanTransitionItem(i.Status, newStatus) {
		return nil, model.NewValidationError("status", "物品状态非法流转")
	}
	i.Status = newStatus
	i.UpdatedAt = time.Now()
	if err := s.store.UpdateItem(i); err != nil {
		return nil, err
	}
	return i, nil
}

func (s *Service) DeleteItem(id string) error {
	return s.store.DeleteItem(id)
}
