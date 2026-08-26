package service

import (
	"sort"
	"time"

	"recommendation/internal/model"
	"recommendation/pkg/idgen"
)

func (s *Service) CreateBehaviorEvent(input model.BehaviorEvent) (*model.BehaviorEvent, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if input.Weight == 0 {
		switch input.EventType {
		case model.EventTypeView:
			input.Weight = 1
		case model.EventTypeLike:
			input.Weight = 3
		case model.EventTypeShare:
			input.Weight = 5
		case model.EventTypeCollect:
			input.Weight = 4
		}
	}
	b := &model.BehaviorEvent{
		ID:         idgen.Hex(),
		UserID:     input.UserID,
		ItemID:     input.ItemID,
		EventType:  input.EventType,
		Weight:     input.Weight,
		OccurredAt: time.Now(),
	}
	if err := s.store.CreateBehaviorEvent(b); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *Service) GetBehaviorEvent(id string) (*model.BehaviorEvent, error) {
	return s.store.GetBehaviorEvent(id)
}

func (s *Service) ListBehaviorEvents(filter model.BehaviorEventFilter, page, size int) ([]*model.BehaviorEvent, int, error) {
	all := s.store.ListBehaviorEvents()
	matched := make([]*model.BehaviorEvent, 0, len(all))
	for _, b := range all {
		if filter.Match(b) {
			matched = append(matched, b)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].OccurredAt.After(matched[j].OccurredAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.BehaviorEvent{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) BatchCreateBehaviorEvents(inputs []model.BehaviorEvent) ([]*model.BehaviorEvent, error) {
	now := time.Now()
	events := make([]*model.BehaviorEvent, 0, len(inputs))
	for _, input := range inputs {
		if err := input.Validate(); err != nil {
			return nil, err
		}
		if input.Weight == 0 {
			switch input.EventType {
			case model.EventTypeView:
				input.Weight = 1
			case model.EventTypeLike:
				input.Weight = 3
			case model.EventTypeShare:
				input.Weight = 5
			case model.EventTypeCollect:
				input.Weight = 4
			}
		}
		b := &model.BehaviorEvent{
			ID:         idgen.Hex(),
			UserID:     input.UserID,
			ItemID:     input.ItemID,
			EventType:  input.EventType,
			Weight:     input.Weight,
			OccurredAt: now,
		}
		events = append(events, b)
	}
	if err := s.store.CreateBehaviorEvents(events); err != nil {
		return nil, err
	}
	return events, nil
}

func (s *Service) DeleteBehaviorEvent(id string) error {
	return s.store.DeleteBehaviorEvent(id)
}

func (s *Service) AggregateItemHeat() map[string]float64 {
	all := s.store.ListBehaviorEvents()
	heat := make(map[string]float64)
	for _, b := range all {
		heat[b.ItemID] += float64(b.Weight)
	}
	return heat
}
