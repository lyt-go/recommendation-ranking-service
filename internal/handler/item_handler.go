package handler

import (
	"net/http"

	"recommendation/internal/model"
	"recommendation/pkg/httpx"
)

func (s *Server) registerItemRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/items", s.createItem)
	mux.HandleFunc("GET /api/items", s.listItems)
	mux.HandleFunc("GET /api/items/{id}", s.getItem)
	mux.HandleFunc("PUT /api/items/{id}", s.updateItem)
	mux.HandleFunc("PUT /api/items/{id}/status", s.updateItemStatus)
	mux.HandleFunc("DELETE /api/items/{id}", s.deleteItem)
}

type createItemRequest struct {
	Title    string   `json:"title"`
	Category string   `json:"category"`
	Tags     []string `json:"tags"`
	Score    float64  `json:"score"`
	Status   string   `json:"status"`
}

func (s *Server) createItem(w http.ResponseWriter, r *http.Request) {
	var req createItemRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	i, err := s.svc.CreateItem(model.Item{Title: req.Title, Category: req.Category, Tags: req.Tags, Score: req.Score, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, i)
}

func (s *Server) listItems(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ItemFilter{
		Category: r.URL.Query().Get("category"),
		Status:   r.URL.Query().Get("status"),
		Tag:      r.URL.Query().Get("tag"),
		Keyword:  r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListItems(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getItem(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	i, err := s.svc.GetItem(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, i)
}

type updateItemRequest struct {
	Title    string   `json:"title"`
	Category string   `json:"category"`
	Tags     []string `json:"tags"`
	Score    float64  `json:"score"`
}

func (s *Server) updateItem(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateItemRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	i, err := s.svc.UpdateItem(id, model.Item{Title: req.Title, Category: req.Category, Tags: req.Tags, Score: req.Score})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, i)
}

type updateItemStatusRequest struct {
	Status string `json:"status"`
}

func (s *Server) updateItemStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateItemStatusRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	i, err := s.svc.UpdateItemStatus(id, req.Status)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, i)
}

func (s *Server) deleteItem(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteItem(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
