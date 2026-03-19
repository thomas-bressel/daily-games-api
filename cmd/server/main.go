package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tbressel/daily-games-api/config"
	"github.com/tbressel/daily-games-api/internal/article"
	"github.com/tbressel/daily-games-api/internal/cache"
	"github.com/tbressel/daily-games-api/internal/handler"
	"github.com/tbressel/daily-games-api/internal/router"
	"github.com/tbressel/daily-games-api/internal/rss"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Load configuration from environment variables
	cfg := config.Load()

	// Connect to Redis cache
	redisClient, err := cache.New(cfg)
	if err != nil {
		logger.Error("Failed to connect to Redis", "error", err)
		os.Exit(1)
	}
	defer redisClient.Close()
	logger.Info("Redis connected", "addr", cfg.RedisAddr)

	// Initialise the RSS parser
	rssParser := rss.New(cfg.FetchTimeoutSeconds, cfg.MaxArticlesPerFeed)

	// Initialise the article orchestrator
	orchestrator := article.New(rssParser, redisClient)

	// Initialise the articles HTTP handler
	articlesHandler := handler.NewArticlesHandler(orchestrator)

	// Build the router with all routes and middlewares
	// Configure the HTTP server with protective timeouts (anti-Slowloris)
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router.Create(articlesHandler),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// This context is cancelled automatically on SIGINT / SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start the server in a goroutine so we can listen for shutdown signals
	go func() {
		logger.Info("Server started", "addr", ":"+cfg.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Server error", "error", err)
			os.Exit(1)
		}
	}()

	// Block until a shutdown signal is received
	<-ctx.Done()
	logger.Info("Shutdown signal received, stopping...")

	// Allow up to 5 seconds for in-flight requests to complete
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("Graceful shutdown failed", "error", err)
	}

	logger.Info("Server stopped cleanly")
}
