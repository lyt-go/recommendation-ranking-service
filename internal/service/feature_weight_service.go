package service

import (
	"sort"
	"time"

	"recommendation/internal/model"
	"recommendation/pkg/idgen"
)

func (s *Service) CreateFeatureWeight(input model.FeatureWeight) (*model.FeatureWeight, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	now := time.Now()
	f := &model.FeatureWeight{
		ID:         idgen.Hex(),
		Feature:    input.Feature,
		StrategyID: input.StrategyID,
		Weight:     input.Weight,
		Enabled:    input.Enabled,
		UpdatedAt:  now,
	}
	if err := s.store.CreateFeatureWeight(f); err != nil {
		return nil, err
	}
	return f, nil
}

func (s *Service) GetFeatureWeight(id string) (*model.FeatureWeight, error) {
	return s.store.GetFeatureWeight(id)
}

func (s *Service) ListFeatureWeights(filter model.FeatureWeightFilter, page, size int) ([]*model.FeatureWeight, int, error) {
	all := s.store.ListFeatureWeights()
	matched := make([]*model.FeatureWeight, 0, len(all))
	for _, f := range all {
		if filter.Match(f) {
			matched = append(matched, f)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].UpdatedAt.After(matched[j].UpdatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.FeatureWeight{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateFeatureWeight(id string, input model.FeatureWeight) (*model.FeatureWeight, error) {
	f, err := s.store.GetFeatureWeight(id)
	if err != nil {
		return nil, err
	}
	if input.Feature != "" {
		f.Feature = input.Feature
	}
	if input.StrategyID != "" {
		f.StrategyID = input.StrategyID
	}
	if input.Weight != 0 {
		f.Weight = input.Weight
	}
	f.Enabled = input.Enabled
	f.UpdatedAt = time.Now()
	if err := f.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateFeatureWeight(f); err != nil {
		return nil, err
	}
	return f, nil
}

func (s *Service) DeleteFeatureWeight(id string) error {
	return s.store.DeleteFeatureWeight(id)
}
