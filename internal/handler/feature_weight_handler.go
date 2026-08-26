package handler

import (
	"net/http"
	"strconv"

	"recommendation/internal/model"
	"recommendation/pkg/httpx"
)

func (s *Server) registerFeatureWeightRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/feature-weights", s.createFeatureWeight)
	mux.HandleFunc("GET /api/feature-weights", s.listFeatureWeights)
	mux.HandleFunc("GET /api/feature-weights/{id}", s.getFeatureWeight)
	mux.HandleFunc("PUT /api/feature-weights/{id}", s.updateFeatureWeight)
	mux.HandleFunc("DELETE /api/feature-weights/{id}", s.deleteFeatureWeight)
}

type createFeatureWeightRequest struct {
	Feature    string  `json:"feature"`
	StrategyID string  `json:"strategy_id"`
	Weight     float64 `json:"weight"`
	Enabled    bool    `json:"enabled"`
}

func (s *Server) createFeatureWeight(w http.ResponseWriter, r *http.Request) {
	var req createFeatureWeightRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	f, err := s.svc.CreateFeatureWeight(model.FeatureWeight{Feature: req.Feature, StrategyID: req.StrategyID, Weight: req.Weight, Enabled: req.Enabled})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, f)
}

func (s *Server) listFeatureWeights(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.FeatureWeightFilter{
		Feature:    r.URL.Query().Get("feature"),
		StrategyID: r.URL.Query().Get("strategy_id"),
	}
	if v := r.URL.Query().Get("enabled"); v != "" {
		b, _ := strconv.ParseBool(v)
		filter.Enabled = &b
	}
	items, total, err := s.svc.ListFeatureWeights(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getFeatureWeight(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	f, err := s.svc.GetFeatureWeight(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, f)
}

type updateFeatureWeightRequest struct {
	Feature    string  `json:"feature"`
	StrategyID string  `json:"strategy_id"`
	Weight     float64 `json:"weight"`
	Enabled    bool    `json:"enabled"`
}

func (s *Server) updateFeatureWeight(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateFeatureWeightRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	f, err := s.svc.UpdateFeatureWeight(id, model.FeatureWeight{Feature: req.Feature, StrategyID: req.StrategyID, Weight: req.Weight, Enabled: req.Enabled})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, f)
}

func (s *Server) deleteFeatureWeight(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteFeatureWeight(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
