// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"recommendation/internal/model"
)

var (
	ErrNotFound = errors.New("记录不存在")
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法，便于测试时替换实现。
type Store interface {
	// UserProfile
	CreateUserProfile(u *model.UserProfile) error
	GetUserProfile(id string) (*model.UserProfile, error)
	GetUserProfileByUserID(userID string) (*model.UserProfile, error)
	ListUserProfiles() []*model.UserProfile
	UpdateUserProfile(u *model.UserProfile) error
	DeleteUserProfile(id string) error

	// Item
	CreateItem(i *model.Item) error
	GetItem(id string) (*model.Item, error)
	ListItems() []*model.Item
	UpdateItem(i *model.Item) error
	DeleteItem(id string) error

	// BehaviorEvent
	CreateBehaviorEvent(b *model.BehaviorEvent) error
	GetBehaviorEvent(id string) (*model.BehaviorEvent, error)
	ListBehaviorEvents() []*model.BehaviorEvent
	DeleteBehaviorEvent(id string) error
	CreateBehaviorEvents(events []*model.BehaviorEvent) error

	// Strategy
	CreateStrategy(s *model.Strategy) error
	GetStrategy(id string) (*model.Strategy, error)
	ListStrategies() []*model.Strategy
	UpdateStrategy(s *model.Strategy) error
	DeleteStrategy(id string) error

	// Recommendation
	CreateRecommendation(r *model.Recommendation) error
	GetRecommendation(id string) (*model.Recommendation, error)
	ListRecommendations() []*model.Recommendation
	UpdateRecommendation(r *model.Recommendation) error
	DeleteRecommendation(id string) error

	// FeatureWeight
	CreateFeatureWeight(f *model.FeatureWeight) error
	GetFeatureWeight(id string) (*model.FeatureWeight, error)
	ListFeatureWeights() []*model.FeatureWeight
	UpdateFeatureWeight(f *model.FeatureWeight) error
	DeleteFeatureWeight(id string) error
}
