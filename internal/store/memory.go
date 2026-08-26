package store

import (
	"sync"

	"recommendation/internal/model"
)

type MemoryStore struct {
	mu              sync.RWMutex
	userProfiles    map[string]*model.UserProfile
	items           map[string]*model.Item
	behaviorEvents  map[string]*model.BehaviorEvent
	strategies      map[string]*model.Strategy
	recommendations map[string]*model.Recommendation
	featureWeights  map[string]*model.FeatureWeight
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		userProfiles:    make(map[string]*model.UserProfile),
		items:           make(map[string]*model.Item),
		behaviorEvents:  make(map[string]*model.BehaviorEvent),
		strategies:      make(map[string]*model.Strategy),
		recommendations: make(map[string]*model.Recommendation),
		featureWeights:  make(map[string]*model.FeatureWeight),
	}
}

var _ Store = (*MemoryStore)(nil)
