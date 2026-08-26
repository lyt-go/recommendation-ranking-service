package model

import (
	"strings"
	"time"
)

type FeatureWeight struct {
	ID         string    `json:"id"`
	Feature    string    `json:"feature"`
	StrategyID string    `json:"strategy_id"`
	Weight     float64   `json:"weight"`
	Enabled    bool      `json:"enabled"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (f *FeatureWeight) Validate() error {
	f.Feature = strings.TrimSpace(f.Feature)
	f.StrategyID = strings.TrimSpace(f.StrategyID)
	if f.Feature == "" {
		return NewValidationError("feature", "特征名称不能为空")
	}
	if f.StrategyID == "" {
		return NewValidationError("strategy_id", "策略ID不能为空")
	}
	return nil
}

type FeatureWeightFilter struct {
	Feature    string
	StrategyID string
	Enabled    *bool
}

func (f FeatureWeightFilter) Match(fw *FeatureWeight) bool {
	if f.Feature != "" && fw.Feature != f.Feature {
		return false
	}
	if f.StrategyID != "" && fw.StrategyID != f.StrategyID {
		return false
	}
	if f.Enabled != nil && fw.Enabled != *f.Enabled {
		return false
	}
	return true
}
