package api

import "net/http"

func registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/healthcheck", HealthcheckHandler())
}
