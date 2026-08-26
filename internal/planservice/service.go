package planservice

import (
	"recommendation/internal/planbuilder"
	"recommendation/internal/plancache"
)

type Service struct{ cache *plancache.Cache }

func New(c *plancache.Cache) *Service { return &Service{cache: c} }
func (s *Service) Generate(key, stage string) (*planbuilder.Plan, error) {
	p := &planbuilder.Plan{Key: key, Features: make(map[string]float64)}
	// 先完成构建再入缓存：构建失败的半成品（Stage 未赋值）绝不能对外可见，
	// 否则后续 StageName 读取 nil Stage 会 panic 导致服务崩溃。
	if err := planbuilder.Populate(p, stage); err != nil {
		return nil, err
	}
	s.cache.Put(p)
	return p, nil
}
func (s *Service) StageName(key string) (string, bool) {
	p, ok := s.cache.Get(key)
	if !ok {
		return "", false
	}
	return p.Stage.Name, true
}
