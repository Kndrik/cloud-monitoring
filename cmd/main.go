package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"time"

	"github.com/Kndrik/cloud-monitoring/internal/api"
	"github.com/Kndrik/cloud-monitoring/internal/data"
	"github.com/jackc/pgx/v5/pgxpool"
)

type dbConfig struct {
	dsn          string
	maxOpenConns int
	maxIdleTime  time.Duration
}

func main() {
	apiConfig := api.Config{}
	dbConfig := dbConfig{}

	flag.IntVar(&apiConfig.Port, "port", 4000, "API server port")
	flag.StringVar(&apiConfig.Env, "env", "development", "Environment (development|staging|production)")

	flag.StringVar(&dbConfig.dsn, "db-dsn", os.Getenv("CLOUDMONITORING_DSN"), "PostgreSQL DSN")
	flag.IntVar(&dbConfig.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.DurationVar(&dbConfig.maxIdleTime, "db-max-idle-time", 15*time.Minute, "PostgreSQL max connection idle time")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

	logger.Info("connecting to the database")
	pool, err := openDB(dbConfig)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer pool.Close()

	models := data.NewModels(pool)
	apiServer := api.New(logger, &apiConfig, &models)

	if err = apiServer.Start(); err != nil {
		logger.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}

func openDB(cfg dbConfig) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(cfg.dsn)
	if err != nil {
		return nil, err
	}

	config.MaxConns = int32(cfg.maxOpenConns)
	config.MaxConnIdleTime = cfg.maxIdleTime

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = pool.Ping(ctx)
	if err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}
