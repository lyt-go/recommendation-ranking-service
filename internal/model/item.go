package model

import (
	"strings"
	"time"
)

const (
	ItemStatusOnline  = "online"
	ItemStatusOffline = "offline"
)

type Item struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Category  string    `json:"category"`
	Tags      []string  `json:"tags"`
	Score     float64   `json:"score"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (i *Item) Validate() error {
	i.Title = strings.TrimSpace(i.Title)
	i.Category = strings.TrimSpace(i.Category)
	if i.Title == "" {
		return NewValidationError("title", "物品标题不能为空")
	}
	if i.Category == "" {
		return NewValidationError("category", "物品分类不能为空")
	}
	if i.Status == "" {
		i.Status = ItemStatusOffline
	}
	if i.Status != ItemStatusOnline && i.Status != ItemStatusOffline {
		return NewValidationError("status", "物品状态不合法")
	}
	return nil
}

var itemTransitions = map[string]map[string]bool{
	ItemStatusOnline:  {ItemStatusOffline: true},
	ItemStatusOffline: {ItemStatusOnline: true},
}

func CanTransitionItem(from, to string) bool {
	if m, ok := itemTransitions[from]; ok {
		return m[to]
	}
	return false
}

type ItemFilter struct {
	Category string
	Status   string
	Tag      string
	Keyword  string
}

func (f ItemFilter) Match(i *Item) bool {
	if f.Category != "" && i.Category != f.Category {
		return false
	}
	if f.Status != "" && i.Status != f.Status {
		return false
	}
	if f.Tag != "" {
		found := false
		for _, t := range i.Tags {
			if t == f.Tag {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(i.Title), k) {
			return false
		}
	}
	return true
}
