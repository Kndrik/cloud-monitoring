package api

import "net/http"

func (s *Server) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/healthcheck", s.HealthcheckHandler())

	mux.HandleFunc("GET /api/v1/instances", s.getInstancesHandler())
	mux.HandleFunc("POST /api/v1/instances", s.addInstanceHandler())
	mux.HandleFunc("DELETE /api/v1/instances/{id}", s.removeInstanceHandler())
}
