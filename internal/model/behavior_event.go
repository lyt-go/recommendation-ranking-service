package model

import (
	"strings"
	"time"
)

const (
	EventTypeView    = "view"
	EventTypeLike    = "like"
	EventTypeShare   = "share"
	EventTypeCollect = "collect"
)

type BehaviorEvent struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	ItemID     string    `json:"item_id"`
	EventType  string    `json:"event_type"`
	Weight     int       `json:"weight"`
	OccurredAt time.Time `json:"occurred_at"`
}

func (b *BehaviorEvent) Validate() error {
	b.UserID = strings.TrimSpace(b.UserID)
	b.ItemID = strings.TrimSpace(b.ItemID)
	b.EventType = strings.TrimSpace(b.EventType)
	if b.UserID == "" {
		return NewValidationError("user_id", "用户ID不能为空")
	}
	if b.ItemID == "" {
		return NewValidationError("item_id", "物品ID不能为空")
	}
	if b.EventType == "" {
		return NewValidationError("event_type", "事件类型不能为空")
	}
	if b.EventType != EventTypeView && b.EventType != EventTypeLike && b.EventType != EventTypeShare && b.EventType != EventTypeCollect {
		return NewValidationError("event_type", "事件类型不合法")
	}
	if b.Weight < 0 {
		return NewValidationError("weight", "权重不能为负数")
	}
	return nil
}

type BehaviorEventFilter struct {
	UserID    string
	ItemID    string
	EventType string
}

func (f BehaviorEventFilter) Match(b *BehaviorEvent) bool {
	if f.UserID != "" && b.UserID != f.UserID {
		return false
	}
	if f.ItemID != "" && b.ItemID != f.ItemID {
		return false
	}
	if f.EventType != "" && b.EventType != f.EventType {
		return false
	}
	return true
}
