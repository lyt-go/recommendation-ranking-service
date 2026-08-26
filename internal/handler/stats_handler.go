package handler

import (
	"net/http"

	"recommendation/pkg/httpx"
)

func (s *Server) registerStatsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/stats/overview", s.getStatsOverview)
}

func (s *Server) getStatsOverview(w http.ResponseWriter, r *http.Request) {
	stats, err := s.svc.GetStatsOverview()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, stats)
}
