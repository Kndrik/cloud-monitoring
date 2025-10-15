package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type Config struct {
	Port int
	Env  string
	Db   struct {
		Dsn          string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  time.Duration
	}
}

type Server struct {
	logger *slog.Logger
	config *Config
}

func New(logger *slog.Logger, config *Config) *Server {
	return &Server{
		logger: logger.With("package", "api"),
		config: config,
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
