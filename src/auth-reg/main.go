package main

import (
	"Auth-Reg/internal/config"
	"Auth-Reg/internal/http-server/handlers/login"
	"Auth-Reg/internal/http-server/handlers/registr"
	smtpMid "Auth-Reg/internal/http-server/middleware/smtp"
	"Auth-Reg/internal/storage/postgres"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	EnvProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("starting auth-reg server", slog.String("env", cfg.Env))

	storage, err := postgres.New(cfg.Storage)
	if err != nil {
		log.Error("failed to init storage")
		os.Exit(1)
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)

	router.Group(func(r chi.Router) {
		r.Use(smtpMid.New(log))
		r.Post("/registration", registr.New(log, storage))
	})

	router.Get("/login", login.New(log, storage))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stoped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case EnvProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}
