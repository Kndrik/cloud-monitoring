package api

import (
	"errors"
	"net/http"
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
		id, err := s.readIdParam(r)
		if err != nil {
			s.notFoundResponse(w, r)
			return
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

func (s *Server) updateInstanceHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := s.readIdParam(r)
		if err != nil {
			s.notFoundResponse(w, r)
			return
		}

		instance, err := s.models.Instances.Get(id)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				s.notFoundResponse(w, r)
			default:
				s.serverErrorResponse(w, r, err)
			}
			return
		}

		var input struct {
			Name        *string `json:"name"`
			Ip          *string `json:"ip"`
			RefreshRate *int    `json:"refresh_rate"`
		}

		err = s.readJSON(w, r, &input)
		if err != nil {
			s.badRequestResponse(w, r, err)
			return
		}

		if input.Name != nil {
			instance.Name = *input.Name
		}

		if input.Ip != nil {
			instance.Ip = *input.Ip
		}

		if input.RefreshRate != nil {
			instance.RefreshRate = time.Duration(*input.RefreshRate) * time.Second
		}

		v := validator.New()

		if data.ValidateInstance(v, instance); !v.Valid() {
			s.failedValidationResponse(w, r, v.Errors)
			return
		}

		err = s.models.Instances.Update(instance)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrEditConflict):
				s.editConflictResponse(w, r)
			default:
				s.serverErrorResponse(w, r, err)
			}
			return
		}

		err = s.writeJSON(w, http.StatusOK, envelope{"instance": instance}, nil)
		if err != nil {
			s.serverErrorResponse(w, r, err)
		}
	}
}
