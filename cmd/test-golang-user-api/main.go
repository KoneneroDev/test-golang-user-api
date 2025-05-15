package main

import (
	"log/slog"
	"os"
	"test_golang_user_api/internal/config"
	"test_golang_user_api/internal/storage/postgres"
)

func main() {
	cfg := config.LoadConfig()

	log := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
	)

	log.Info("starting server ", slog.String("env", cfg.Env))

	storage, err := postgres.New(cfg.Data.Postgres)
	if err != nil {
		log.Error("failed to connect to database", err)
		os.Exit(1)
	}
	log.Info("finished connect to db")

	_ = storage

	//todo: router

	//todo: start
}
