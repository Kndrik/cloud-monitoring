package api

import "net/http"

type Instance struct {
	Name        string `json:"name"`
	Ip          string `json:"ip"`
	RefreshRate int32  `json:"refreshrate"`
}

func (s *Server) getInstancesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := envelope{
			"instances": []Instance{
				{
					Name:        "backend",
					Ip:          "127.0.0.1:5000",
					RefreshRate: 60,
				},
				{
					Name:        "frontend",
					Ip:          "127.0.0.1:5001",
					RefreshRate: 120,
				},
			},
		}

		err := s.writeJSON(w, http.StatusOK, data, nil)
		if err != nil {
			s.serverErrorResponse(w, r, err)
		}
	}
}
