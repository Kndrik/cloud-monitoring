package api

import "net/http"

func (s *Server) HealthcheckHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		count, err := s.models.Instances.Count()
		if err != nil {
			s.serverErrorResponse(w, r, err)
			return
		}

		status := envelope{
			"status":    "available",
			"instances": count,
		}

		err = s.writeJSON(w, http.StatusOK, status, nil)
		if err != nil {
			s.serverErrorResponse(w, r, err)
		}
	}
}
