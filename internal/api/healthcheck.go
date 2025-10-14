package api

import "net/http"

func (s *Server) HealthcheckHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := envelope{
			"status":    "available",
			"instances": 12,
		}

		err := s.writeJSON(w, http.StatusOK, status, nil)
		if err != nil {
			s.serverErrorResponse(w, r, err)
		}
	}
}
