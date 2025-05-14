package main

import (
	"log/slog"
	"os"
	"test_golang_user_api/internal/config"
)

func main() {
	cfg := config.LoadConfig()

	log := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
	)

	log.Info("starting server ", slog.String("env", cfg.Env))

	//todo: storage

	//todo: router

	//todo: start
}
