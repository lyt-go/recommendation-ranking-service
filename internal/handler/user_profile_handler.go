package handler

import (
	"net/http"

	"recommendation/internal/model"
	"recommendation/pkg/httpx"
)

func (s *Server) registerUserProfileRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/user-profiles", s.createUserProfile)
	mux.HandleFunc("GET /api/user-profiles", s.listUserProfiles)
	mux.HandleFunc("GET /api/user-profiles/{id}", s.getUserProfile)
	mux.HandleFunc("PUT /api/user-profiles/{id}", s.updateUserProfile)
	mux.HandleFunc("DELETE /api/user-profiles/{id}", s.deleteUserProfile)
}

type createUserProfileRequest struct {
	UserID    string   `json:"user_id"`
	Interests []string `json:"interests"`
	Tags      []string `json:"tags"`
	Region    string   `json:"region"`
}

func (s *Server) createUserProfile(w http.ResponseWriter, r *http.Request) {
	var req createUserProfileRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	u, err := s.svc.CreateUserProfile(model.UserProfile{UserID: req.UserID, Interests: req.Interests, Tags: req.Tags, Region: req.Region})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, u)
}

func (s *Server) listUserProfiles(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.UserProfileFilter{
		Region: r.URL.Query().Get("region"),
		Tag:    r.URL.Query().Get("tag"),
	}
	items, total, err := s.svc.ListUserProfiles(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getUserProfile(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	u, err := s.svc.GetUserProfile(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, u)
}

type updateUserProfileRequest struct {
	UserID    string   `json:"user_id"`
	Interests []string `json:"interests"`
	Tags      []string `json:"tags"`
	Region    string   `json:"region"`
}

func (s *Server) updateUserProfile(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateUserProfileRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	u, err := s.svc.UpdateUserProfile(id, model.UserProfile{UserID: req.UserID, Interests: req.Interests, Tags: req.Tags, Region: req.Region})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, u)
}

func (s *Server) deleteUserProfile(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteUserProfile(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}
