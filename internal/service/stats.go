package service

import (
	"sort"
)

type StatsOverview struct {
	UserCount                int            `json:"user_count"`
	ItemCount                int            `json:"item_count"`
	BehaviorEventTypeDist    map[string]int `json:"behavior_event_type_dist"`
	HotItems                 []HotItem      `json:"hot_items"`
	RecommendationStatusDist map[string]int `json:"recommendation_status_dist"`
}

type HotItem struct {
	ItemID string  `json:"item_id"`
	Heat   float64 `json:"heat"`
}

func (s *Service) GetStatsOverview() (*StatsOverview, error) {
	users := s.store.ListUserProfiles()
	items := s.store.ListItems()
	events := s.store.ListBehaviorEvents()
	recommendations := s.store.ListRecommendations()

	eventDist := make(map[string]int)
	for _, e := range events {
		eventDist[e.EventType]++
	}

	heat := s.AggregateItemHeat()
	type itemHeat struct {
		itemID string
		heat   float64
	}
	heats := make([]itemHeat, 0, len(heat))
	for id, h := range heat {
		heats = append(heats, itemHeat{itemID: id, heat: h})
	}
	sort.Slice(heats, func(i, j int) bool {
		return heats[i].heat > heats[j].heat
	})

	topN := 10
	if s.cfg != nil && s.cfg.DefaultTopN > 0 {
		topN = s.cfg.DefaultTopN
	}
	hotItems := make([]HotItem, 0, topN)
	for i := 0; i < len(heats) && i < topN; i++ {
		hotItems = append(hotItems, HotItem{ItemID: heats[i].itemID, Heat: heats[i].heat})
	}

	recDist := make(map[string]int)
	for _, r := range recommendations {
		recDist[r.Status]++
	}

	return &StatsOverview{
		UserCount:                len(users),
		ItemCount:                len(items),
		BehaviorEventTypeDist:    eventDist,
		HotItems:                 hotItems,
		RecommendationStatusDist: recDist,
	}, nil
}
