package api

import "net/http"

func HealthcheckHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := envelope{
			"status":    "available",
			"instances": 12,
		}

		err := writeJSON(w, http.StatusOK, status, nil)
		if err != nil {
			http.Error(w, "failed to encode JSON", http.StatusInternalServerError)
		}
	}
}
