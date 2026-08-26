package handler

import (
	"net/http"

	"recommendation/internal/model"
	"recommendation/pkg/httpx"
)

func (s *Server) registerRecommendationRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/recommendations", s.createRecommendation)
	mux.HandleFunc("GET /api/recommendations", s.listRecommendations)
	mux.HandleFunc("GET /api/recommendations/{id}", s.getRecommendation)
	mux.HandleFunc("PUT /api/recommendations/{id}", s.updateRecommendation)
	mux.HandleFunc("PUT /api/recommendations/{id}/status", s.updateRecommendationStatus)
	mux.HandleFunc("POST /api/recommendations/generate", s.generateRecommendation)
	mux.HandleFunc("DELETE /api/recommendations/{id}", s.deleteRecommendation)
}

type createRecommendationRequest struct {
	UserID     string    `json:"user_id"`
	StrategyID string    `json:"strategy_id"`
	ItemIDs    []string  `json:"item_ids"`
	Scores     []float64 `json:"scores"`
	Status     string    `json:"status"`
}

func (s *Server) createRecommendation(w http.ResponseWriter, r *http.Request) {
	var req createRecommendationRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rec, err := s.svc.CreateRecommendation(model.Recommendation{UserID: req.UserID, StrategyID: req.StrategyID, ItemIDs: req.ItemIDs, Scores: req.Scores, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, rec)
}

func (s *Server) listRecommendations(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.RecommendationFilter{
		UserID:     r.URL.Query().Get("user_id"),
		StrategyID: r.URL.Query().Get("strategy_id"),
		Status:     r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListRecommendations(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getRecommendation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	rec, err := s.svc.GetRecommendation(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rec)
}

type updateRecommendationRequest struct {
	StrategyID string    `json:"strategy_id"`
	ItemIDs    []string  `json:"item_ids"`
	Scores     []float64 `json:"scores"`
}

func (s *Server) updateRecommendation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateRecommendationRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rec, err := s.svc.UpdateRecommendation(id, model.Recommendation{StrategyID: req.StrategyID, ItemIDs: req.ItemIDs, Scores: req.Scores})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rec)
}

type updateRecommendationStatusRequest struct {
	Status string `json:"status"`
}

func (s *Server) updateRecommendationStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateRecommendationStatusRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rec, err := s.svc.UpdateRecommendationStatus(id, req.Status)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rec)
}

type generateRecommendationRequest struct {
	UserID     string `json:"user_id"`
	StrategyID string `json:"strategy_id"`
	TopN       int    `json:"top_n"`
}

func (s *Server) generateRecommendation(w http.ResponseWriter, r *http.Request) {
	var req generateRecommendationRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rec, err := s.svc.GenerateRecommendation(req.UserID, req.StrategyID, req.TopN)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, rec)
}

func (s *Server) deleteRecommendation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteRecommendation(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
