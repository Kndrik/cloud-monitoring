package api

import "net/http"

func (s *Server) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/healthcheck", s.HealthcheckHandler())

	mux.HandleFunc("GET /api/v1/instances", s.getInstancesHandler())
}
