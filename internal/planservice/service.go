package planservice

import (
	"recommendation/internal/planbuilder"
	"recommendation/internal/plancache"
)

type Service struct{ cache *plancache.Cache }

func New(c *plancache.Cache) *Service { return &Service{cache: c} }
func (s *Service) Generate(key, stage string) (*planbuilder.Plan, error) {
	p := &planbuilder.Plan{Key: key, Features: make(map[string]float64)}
	s.cache.Put(p)
	if err := planbuilder.Populate(p, stage); err != nil {
		return p, err
	}
	return p, nil
}
func (s *Service) StageName(key string) (string, bool) {
	p, ok := s.cache.Get(key)
	if !ok {
		return "", false
	}
	return p.Stage.Name, true
}
