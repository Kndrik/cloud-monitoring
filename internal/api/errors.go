package api

import (
	"net/http"
)

func (s *Server) logError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	s.logger.Error(err.Error(), "method", method, "uri", uri)
}

func (s *Server) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message}

	err := s.writeJSON(w, status, env, nil)
	if err != nil {
		s.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	s.logError(r, err)

	message := "the server encountered a problem and could not process the request"
	s.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (s *Server) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	s.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (s *Server) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	s.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (s *Server) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	s.errorResponse(w, r, http.StatusNotFound, "resource not found")
}

func (s *Server) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to edit the record due to an edit conflict, please try again."
	s.errorResponse(w, r, http.StatusConflict, message)
}
