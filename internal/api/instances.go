package api

import (
	"errors"
	"net/http"
	"strconv"
)

type Instance struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Ip          string `json:"ip"`
	RefreshRate int32  `json:"refreshrate"`
}

var instances []Instance = []Instance{}

func (s *Server) getInstancesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := envelope{
			"instances": instances,
		}

		err := s.writeJSON(w, http.StatusOK, data, nil)
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
			RefreshRate int32  `json:"refreshrate"`
		}

		err := s.readJSON(w, r, &input)
		if err != nil {
			s.badRequestResponse(w, r, err)
			return
		}

		newInstance := Instance{
			Id:          len(instances) + 1,
			Name:        input.Name,
			Ip:          input.Ip,
			RefreshRate: input.RefreshRate,
		}

		instances = append(instances, newInstance)

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

		for i, instance := range instances {
			if instance.Id == int(id) {
				instances = append(instances[:i], instances[i+1:]...)
				break
			}
		}

		err = s.writeJSON(w, http.StatusOK, envelope{"message": "instance successfully removed"}, nil)
		if err != nil {
			s.serverErrorResponse(w, r, err)
		}
	}
}
