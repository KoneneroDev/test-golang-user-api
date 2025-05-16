package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"test_golang_user_api/internal/config"
	dr "test_golang_user_api/internal/http_server/handlers/uri/delete"
	"test_golang_user_api/internal/http_server/handlers/uri/get"
	"test_golang_user_api/internal/http_server/handlers/uri/patch"
	"test_golang_user_api/internal/http_server/handlers/uri/save"
	"test_golang_user_api/internal/storage/postgres"
)

func main() {
	cfg := config.LoadConfig()

	log := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
	)
	log.Info("starting server with", slog.String("env", cfg.Env))
	log.Info("starting connect to db")

	storage, err := postgres.New(cfg.Data.Postgres)
	if err != nil {
		log.Error("failed to connect to database", err)
		os.Exit(1)
	}
	log.Info("finished connect to db")

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/user", save.New(log, storage))
	router.Delete("/user/{id}", dr.New(log, storage))
	router.Get("/user/{id}", get.New(log, storage))
	router.Patch("/user/{id}", patch.New(log, storage))

	server := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	log.Info("starting http server on ", slog.String("host", cfg.HTTPServer.Address))

	if err := server.ListenAndServe(); err != nil {
		log.Error("failed to start http server", err)
	}

}
