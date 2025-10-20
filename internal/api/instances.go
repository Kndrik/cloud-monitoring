package api

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/Kndrik/cloud-monitoring/internal/data"
	"github.com/Kndrik/cloud-monitoring/internal/validator"
)

var instances []data.Instance = []data.Instance{}

func (s *Server) getInstancesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		instances, err := s.models.Instances.GetAll()
		if err != nil {
			s.serverErrorResponse(w, r, err)
			return
		}

		err = s.writeJSON(w, http.StatusOK, envelope{"instances": instances}, nil)
		if err != nil {
			s.serverErrorResponse(w, r, err)
		}
	}
}

func (s *Server) addInstanceHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input struct {
			Name        string `json:"name"`
			Ip          string `json:"ip"`
			RefreshRate int    `json:"refresh_rate"`
		}

		err := s.readJSON(w, r, &input)
		if err != nil {
			s.badRequestResponse(w, r, err)
			return
		}

		newInstance := &data.Instance{
			Name:        input.Name,
			Ip:          input.Ip,
			RefreshRate: time.Duration(input.RefreshRate) * time.Second,
		}

		v := validator.New()
		if data.ValidateInstance(v, newInstance); !v.Valid() {
			s.failedValidationResponse(w, r, v.Errors)
			return
		}

		err = s.models.Instances.Insert(newInstance)
		if err != nil {
			s.serverErrorResponse(w, r, err)
			return
		}

		err = s.writeJSON(w, http.StatusCreated, envelope{"instance": newInstance}, nil)
		if err != nil {
			s.serverErrorResponse(w, r, err)
		}
	}
}

func (s *Server) removeInstanceHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pathId := r.PathValue("id")
		if pathId == "" {
			s.badRequestResponse(w, r, errors.New("id cannot be empty"))
		}

		id, err := strconv.ParseInt(pathId, 10, 0)
		if err != nil {
			s.badRequestResponse(w, r, err)
		}

		err = s.models.Instances.Delete(id)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				s.notFoundResponse(w, r)
			default:
				s.serverErrorResponse(w, r, err)
			}
			return
		}

		err = s.writeJSON(w, http.StatusOK, envelope{"message": "instance successfully removed"}, nil)
		if err != nil {
			s.serverErrorResponse(w, r, err)
		}
	}
}
