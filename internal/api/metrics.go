package api

import "net/http"

func (s *Server) getInstanceMetricsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := s.readIdParam(r)
		if err != nil {
			s.notFoundResponse(w, r)
			return
		}

		metrics, err := s.models.Metrics.GetInstanceMetrics(int(id))
		if err != nil {
			s.serverErrorResponse(w, r, err)
			return
		}

		err = s.writeJSON(w, http.StatusOK, envelope{"metrics": metrics}, nil)
		if err != nil {
			s.serverErrorResponse(w, r, err)
		}
	}
}
