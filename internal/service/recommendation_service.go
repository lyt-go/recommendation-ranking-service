package service

import (
	"sort"
	"time"

	"recommendation/internal/model"
	"recommendation/internal/store"
	"recommendation/pkg/idgen"
)

func (s *Service) CreateRecommendation(input model.Recommendation) (*model.Recommendation, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetStrategy(input.StrategyID); err != nil {
		if err == store.ErrNotFound {
			return nil, model.NewValidationError("strategy_id", "策略不存在")
		}
		return nil, err
	}
	now := time.Now()
	r := &model.Recommendation{
		ID:         idgen.Hex(),
		UserID:     input.UserID,
		StrategyID: input.StrategyID,
		ItemIDs:    input.ItemIDs,
		Scores:     input.Scores,
		Status:     input.Status,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := s.store.CreateRecommendation(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Service) GetRecommendation(id string) (*model.Recommendation, error) {
	return s.store.GetRecommendation(id)
}

func (s *Service) ListRecommendations(filter model.RecommendationFilter, page, size int) ([]*model.Recommendation, int, error) {
	all := s.store.ListRecommendations()
	matched := make([]*model.Recommendation, 0, len(all))
	for _, r := range all {
		if filter.Match(r) {
			matched = append(matched, r)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Recommendation{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateRecommendation(id string, input model.Recommendation) (*model.Recommendation, error) {
	r, err := s.store.GetRecommendation(id)
	if err != nil {
		return nil, err
	}
	if input.StrategyID != "" {
		r.StrategyID = input.StrategyID
	}
	if input.ItemIDs != nil {
		r.ItemIDs = input.ItemIDs
	}
	if input.Scores != nil {
		r.Scores = input.Scores
	}
	r.UpdatedAt = time.Now()
	if err := r.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateRecommendation(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Service) UpdateRecommendationStatus(id string, newStatus string) (*model.Recommendation, error) {
	r, err := s.store.GetRecommendation(id)
	if err != nil {
		return nil, err
	}
	if !model.CanTransitionRecommendation(r.Status, newStatus) {
		return nil, model.NewValidationError("status", "推荐结果状态非法流转")
	}
	r.Status = newStatus
	r.UpdatedAt = time.Now()
	if err := s.store.UpdateRecommendation(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Service) DeleteRecommendation(id string) error {
	return s.store.DeleteRecommendation(id)
}

func (s *Service) GenerateRecommendation(userID string, strategyID string, topN int) (*model.Recommendation, error) {
	if userID == "" {
		return nil, model.NewValidationError("user_id", "用户ID不能为空")
	}
	if strategyID == "" {
		return nil, model.NewValidationError("strategy_id", "策略ID不能为空")
	}
	strategy, err := s.store.GetStrategy(strategyID)
	if err != nil {
		if err == store.ErrNotFound {
			return nil, model.NewValidationError("strategy_id", "策略不存在")
		}
		return nil, err
	}
	if strategy.Status != model.StrategyStatusEnabled {
		return nil, model.NewValidationError("strategy_id", "策略未启用")
	}

	allItems := s.store.ListItems()
	candidates := make([]*model.Item, 0, len(allItems))
	for _, item := range allItems {
		if item.Status == model.ItemStatusOnline {
			candidates = append(candidates, item)
		}
	}

	heat := s.AggregateItemHeat()
	weights := s.store.ListFeatureWeights()
	strategyWeights := make(map[string]float64)
	for _, fw := range weights {
		if fw.StrategyID == strategyID && fw.Enabled {
			strategyWeights[fw.Feature] = fw.Weight
		}
	}

	type scoredItem struct {
		item  *model.Item
		score float64
	}
	scored := make([]scoredItem, 0, len(candidates))
	for _, item := range candidates {
		score := item.Score
		for _, tag := range item.Tags {
			if w, ok := strategyWeights[tag]; ok {
				score += w
			}
		}
		if h, ok := heat[item.ID]; ok {
			score += h
		}
		scored = append(scored, scoredItem{item: item, score: score})
	}

	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	if topN <= 0 {
		if s.cfg != nil && s.cfg.DefaultTopN > 0 {
			topN = s.cfg.DefaultTopN
		} else {
			topN = 10
		}
	}
	if topN > len(scored) {
		topN = len(scored)
	}

	itemIDs := make([]string, topN)
	scores := make([]float64, topN)
	for i := 0; i < topN; i++ {
		itemIDs[i] = scored[i].item.ID
		scores[i] = scored[i].score
	}

	now := time.Now()
	r := &model.Recommendation{
		ID:         idgen.Hex(),
		UserID:     userID,
		StrategyID: strategyID,
		ItemIDs:    itemIDs,
		Scores:     scores,
		Status:     model.RecommendationStatusDraft,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := s.store.CreateRecommendation(r); err != nil {
		return nil, err
	}
	return r, nil
}
