package handler

import (
	"net/http"

	"recommendation/internal/model"
	"recommendation/pkg/httpx"
)

func (s *Server) registerBehaviorEventRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/behavior-events", s.createBehaviorEvent)
	mux.HandleFunc("POST /api/behavior-events/batch", s.batchCreateBehaviorEvents)
	mux.HandleFunc("GET /api/behavior-events", s.listBehaviorEvents)
	mux.HandleFunc("GET /api/behavior-events/{id}", s.getBehaviorEvent)
	mux.HandleFunc("DELETE /api/behavior-events/{id}", s.deleteBehaviorEvent)
}

type createBehaviorEventRequest struct {
	UserID    string `json:"user_id"`
	ItemID    string `json:"item_id"`
	EventType string `json:"event_type"`
	Weight    int    `json:"weight"`
}

func (s *Server) createBehaviorEvent(w http.ResponseWriter, r *http.Request) {
	var req createBehaviorEventRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	b, err := s.svc.CreateBehaviorEvent(model.BehaviorEvent{UserID: req.UserID, ItemID: req.ItemID, EventType: req.EventType, Weight: req.Weight})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, b)
}

func (s *Server) listBehaviorEvents(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.BehaviorEventFilter{
		UserID:    r.URL.Query().Get("user_id"),
		ItemID:    r.URL.Query().Get("item_id"),
		EventType: r.URL.Query().Get("event_type"),
	}
	items, total, err := s.svc.ListBehaviorEvents(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getBehaviorEvent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	b, err := s.svc.GetBehaviorEvent(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, b)
}

type batchCreateBehaviorEventsRequest struct {
	Events []createBehaviorEventRequest `json:"events"`
}

func (s *Server) batchCreateBehaviorEvents(w http.ResponseWriter, r *http.Request) {
	var req batchCreateBehaviorEventsRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	inputs := make([]model.BehaviorEvent, 0, len(req.Events))
	for _, e := range req.Events {
		inputs = append(inputs, model.BehaviorEvent{UserID: e.UserID, ItemID: e.ItemID, EventType: e.EventType, Weight: e.Weight})
	}
	result, err := s.svc.BatchCreateBehaviorEvents(inputs)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, result)
}

func (s *Server) deleteBehaviorEvent(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteBehaviorEvent(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
