package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"test_task/config"
	"test_task/internal/server/handlers"
	service2 "test_task/internal/service"
	"test_task/internal/storage"
	"test_task/middleware/logger"
	"time"
)

func main() {
	cfg := config.GetConfig()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	clickStorage, err := storage.NewClickStorage(ctx, cfg.ConnectionString)
	if err != nil {
		slog.Error("Failed to initialize storage", "error", err)
		os.Exit(1)
	}

	clickCounter := storage.NewClickCounter(clickStorage)
	defer clickCounter.Shutdown()

	serviceClickCounter := service2.NewClickCounter(clickCounter)
	serviceClickStats := service2.NewClickStats(clickStorage)

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(logger.New(log))
	router.HandleFunc("GET /counter/{bannerID}", handlers.NewClickHandler(serviceClickCounter))
	router.HandleFunc("POST /stats/{bannerID}", handlers.NewStatsHandler(serviceClickStats))
	http.ListenAndServe(cfg.Port, router)
}
