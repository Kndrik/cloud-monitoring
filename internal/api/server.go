package api

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Kndrik/cloud-monitoring/internal/data"
)

type Config struct {
	Port int
	Env  string
}

type Server struct {
	logger *slog.Logger
	config *Config
	models *data.Models
}

func New(logger *slog.Logger, config *Config, models *data.Models) *Server {
	return &Server{
		logger: logger.With("package", "api"),
		config: config,
		models: models,
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	s.registerRoutes(mux)

	srv := &http.Server{
		Addr:     fmt.Sprintf(":%d", s.config.Port),
		Handler:  mux,
		ErrorLog: slog.NewLogLogger(s.logger.Handler(), slog.LevelError),
	}

	s.logger.Info("starting server", "address", srv.Addr, "env", s.config.Env)

	err := srv.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
