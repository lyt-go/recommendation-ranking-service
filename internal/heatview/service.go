package heatview

import (
	"recommendation/internal/heatcache"
	"recommendation/internal/heatstore"
)

type Service struct {
	store *heatstore.Store
	cache *heatcache.Cache
}

func New(store *heatstore.Store, cache *heatcache.Cache) *Service {
	return &Service{store: store, cache: cache}
}
func (s *Service) Capture() { s.cache.Store(s.store.Snapshot()) }
