package model

import (
	"strings"
	"time"
)

const (
	StrategyTypeCollaborative = "collaborative"
	StrategyTypeContent       = "content"
	StrategyTypeHot           = "hot"
)

const (
	StrategyStatusEnabled  = "enabled"
	StrategyStatusDisabled = "disabled"
)

type Strategy struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Params    string    `json:"params"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *Strategy) Validate() error {
	s.Name = strings.TrimSpace(s.Name)
	s.Type = strings.TrimSpace(s.Type)
	if s.Name == "" {
		return NewValidationError("name", "策略名称不能为空")
	}
	if s.Type == "" {
		return NewValidationError("type", "策略类型不能为空")
	}
	if s.Type != StrategyTypeCollaborative && s.Type != StrategyTypeContent && s.Type != StrategyTypeHot {
		return NewValidationError("type", "策略类型不合法")
	}
	if s.Status == "" {
		s.Status = StrategyStatusDisabled
	}
	if s.Status != StrategyStatusEnabled && s.Status != StrategyStatusDisabled {
		return NewValidationError("status", "策略状态不合法")
	}
	return nil
}

var strategyTransitions = map[string]map[string]bool{
	StrategyStatusEnabled:  {StrategyStatusDisabled: true},
	StrategyStatusDisabled: {StrategyStatusEnabled: true},
}

func CanTransitionStrategy(from, to string) bool {
	if m, ok := strategyTransitions[from]; ok {
		return m[to]
	}
	return false
}

type StrategyFilter struct {
	Name   string
	Type   string
	Status string
}

func (f StrategyFilter) Match(s *Strategy) bool {
	if f.Name != "" && !strings.Contains(strings.ToLower(s.Name), strings.ToLower(strings.TrimSpace(f.Name))) {
		return false
	}
	if f.Type != "" && s.Type != f.Type {
		return false
	}
	if f.Status != "" && s.Status != f.Status {
		return false
	}
	return true
}
