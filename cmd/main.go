package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/Kndrik/cloud-monitoring/internal/api"
)

func main() {
	cfg := api.Config{}

	flag.IntVar(&cfg.Port, "port", 4000, "API server port")
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

	apiServer := api.New(logger, &cfg)

	if err := apiServer.Start(); err != nil {
		logger.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}
