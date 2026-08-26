package model

import (
	"strings"
	"time"
)

const (
	RecommendationStatusDraft     = "draft"
	RecommendationStatusPublished = "published"
	RecommendationStatusExpired   = "expired"
)

type Recommendation struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	StrategyID string    `json:"strategy_id"`
	ItemIDs    []string  `json:"item_ids"`
	Scores     []float64 `json:"scores"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (r *Recommendation) Validate() error {
	r.UserID = strings.TrimSpace(r.UserID)
	r.StrategyID = strings.TrimSpace(r.StrategyID)
	if r.UserID == "" {
		return NewValidationError("user_id", "用户ID不能为空")
	}
	if r.StrategyID == "" {
		return NewValidationError("strategy_id", "策略ID不能为空")
	}
	if len(r.ItemIDs) == 0 {
		return NewValidationError("item_ids", "推荐物品列表不能为空")
	}
	if len(r.ItemIDs) != len(r.Scores) {
		return NewValidationError("scores", "得分列表与物品列表长度不一致")
	}
	if r.Status == "" {
		r.Status = RecommendationStatusDraft
	}
	if r.Status != RecommendationStatusDraft && r.Status != RecommendationStatusPublished && r.Status != RecommendationStatusExpired {
		return NewValidationError("status", "推荐结果状态不合法")
	}
	return nil
}

var recommendationTransitions = map[string]map[string]bool{
	RecommendationStatusDraft:     {RecommendationStatusPublished: true},
	RecommendationStatusPublished: {RecommendationStatusExpired: true},
	RecommendationStatusExpired:   {},
}

func CanTransitionRecommendation(from, to string) bool {
	if m, ok := recommendationTransitions[from]; ok {
		return m[to]
	}
	return false
}

type RecommendationFilter struct {
	UserID     string
	StrategyID string
	Status     string
}

func (f RecommendationFilter) Match(r *Recommendation) bool {
	if f.UserID != "" && r.UserID != f.UserID {
		return false
	}
	if f.StrategyID != "" && r.StrategyID != f.StrategyID {
		return false
	}
	if f.Status != "" && r.Status != f.Status {
		return false
	}
	return true
}
