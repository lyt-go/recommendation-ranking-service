package service

import (
	"sort"
	"time"

	"recommendation/internal/model"
	"recommendation/pkg/idgen"
)

func (s *Service) CreateStrategy(input model.Strategy) (*model.Strategy, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	now := time.Now()
	st := &model.Strategy{
		ID:        idgen.Hex(),
		Name:      input.Name,
		Type:      input.Type,
		Params:    input.Params,
		Status:    input.Status,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.store.CreateStrategy(st); err != nil {
		return nil, err
	}
	return st, nil
}

func (s *Service) GetStrategy(id string) (*model.Strategy, error) {
	return s.store.GetStrategy(id)
}

func (s *Service) ListStrategies(filter model.StrategyFilter, page, size int) ([]*model.Strategy, int, error) {
	all := s.store.ListStrategies()
	matched := make([]*model.Strategy, 0, len(all))
	for _, st := range all {
		if filter.Match(st) {
			matched = append(matched, st)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Strategy{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateStrategy(id string, input model.Strategy) (*model.Strategy, error) {
	st, err := s.store.GetStrategy(id)
	if err != nil {
		return nil, err
	}
	if input.Name != "" {
		st.Name = input.Name
	}
	if input.Type != "" {
		st.Type = input.Type
	}
	if input.Params != "" {
		st.Params = input.Params
	}
	st.UpdatedAt = time.Now()
	if err := st.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateStrategy(st); err != nil {
		return nil, err
	}
	return st, nil
}

func (s *Service) UpdateStrategyStatus(id string, newStatus string) (*model.Strategy, error) {
	st, err := s.store.GetStrategy(id)
	if err != nil {
		return nil, err
	}
	if !model.CanTransitionStrategy(st.Status, newStatus) {
		return nil, model.NewValidationError("status", "策略状态非法流转")
	}
	st.Status = newStatus
	st.UpdatedAt = time.Now()
	if err := s.store.UpdateStrategy(st); err != nil {
		return nil, err
	}
	return st, nil
}

func (s *Service) DeleteStrategy(id string) error {
	return s.store.DeleteStrategy(id)
}
