package handler

import (
	"net/http"

	"recommendation/internal/model"
	"recommendation/pkg/httpx"
)

func (s *Server) registerStrategyRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/strategies", s.createStrategy)
	mux.HandleFunc("GET /api/strategies", s.listStrategies)
	mux.HandleFunc("GET /api/strategies/{id}", s.getStrategy)
	mux.HandleFunc("PUT /api/strategies/{id}", s.updateStrategy)
	mux.HandleFunc("PUT /api/strategies/{id}/status", s.updateStrategyStatus)
	mux.HandleFunc("DELETE /api/strategies/{id}", s.deleteStrategy)
}

type createStrategyRequest struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Params string `json:"params"`
	Status string `json:"status"`
}

func (s *Server) createStrategy(w http.ResponseWriter, r *http.Request) {
	var req createStrategyRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	st, err := s.svc.CreateStrategy(model.Strategy{Name: req.Name, Type: req.Type, Params: req.Params, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, st)
}

func (s *Server) listStrategies(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.StrategyFilter{
		Name:   r.URL.Query().Get("name"),
		Type:   r.URL.Query().Get("type"),
		Status: r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListStrategies(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getStrategy(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	st, err := s.svc.GetStrategy(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, st)
}

type updateStrategyRequest struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Params string `json:"params"`
}

func (s *Server) updateStrategy(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateStrategyRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	st, err := s.svc.UpdateStrategy(id, model.Strategy{Name: req.Name, Type: req.Type, Params: req.Params})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, st)
}

type updateStrategyStatusRequest struct {
	Status string `json:"status"`
}

func (s *Server) updateStrategyStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateStrategyStatusRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	st, err := s.svc.UpdateStrategyStatus(id, req.Status)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, st)
}

func (s *Server) deleteStrategy(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteStrategy(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
