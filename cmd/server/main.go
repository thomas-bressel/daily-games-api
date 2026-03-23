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
	"github.com/tbressel/daily-games-api/internal/metrics"
	"github.com/tbressel/daily-games-api/internal/router"
	"github.com/tbressel/daily-games-api/internal/rss"
	"github.com/tbressel/daily-games-api/pkg"
)

// cacheWarmupCombinations lists all category+lang combinations to pre-warm.
// Each entry maps directly to a Redis cache key.
var cacheWarmupCombinations = []pkg.ArticleFilters{
	{Category: "nextgen", Lang: "fr", Limit: 200, Refresh: true},
	{Category: "nextgen", Lang: "en", Limit: 200, Refresh: true},
	{Category: "retrogaming", Lang: "fr", Limit: 200, Refresh: true},
	{Category: "retrogaming", Lang: "en", Limit: 200, Refresh: true},
	{Category: "indie", Lang: "fr", Limit: 200, Refresh: true},
	{Category: "indie", Lang: "en", Limit: 200, Refresh: true},
	{Category: "homebrew", Lang: "en", Limit: 200, Refresh: true},
	{Category: "computing", Lang: "fr", Limit: 200, Refresh: true},
	{Category: "computing", Lang: "en", Limit: 200, Refresh: true},
	{Category: "esport", Lang: "fr", Limit: 200, Refresh: true},
	{Category: "esport", Lang: "en", Limit: 200, Refresh: true},
}

// startCacheWarmer launches a background goroutine that pre-warms the Redis cache
// immediately on startup, then repeats every interval until ctx is cancelled.
func startCacheWarmer(ctx context.Context, o *article.Orchestrator, interval time.Duration) {
	go func() {
		warm := func() {
			start := time.Now()
			for _, filters := range cacheWarmupCombinations {
				if ctx.Err() != nil {
					return
				}
				_, err := o.GetArticles(ctx, filters)
				if err != nil {
					slog.Warn("[CacheWarmer] Failed to warm cache", "category", filters.Category, "lang", filters.Lang, "err", err)
					metrics.WarmerErrors.Inc()
				} else {
					slog.Info("[CacheWarmer] Warmed", "category", filters.Category, "lang", filters.Lang)
				}
			}
			metrics.WarmerDuration.Observe(time.Since(start).Seconds())
		}

		slog.Info("[CacheWarmer] Initial warm-up started")
		warm()
		slog.Info("[CacheWarmer] Initial warm-up done")

		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				slog.Info("[CacheWarmer] Stopped")
				return
			case <-ticker.C:
				slog.Info("[CacheWarmer] Refreshing cache")
				warm()
			}
		}
	}()
}

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

	// Initialise the track HTTP handler
	trackHandler := handler.NewTrackHandler(redisClient)

	// Build the router with all routes and middlewares
	// Configure the HTTP server with protective timeouts (anti-Slowloris)
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router.Create(articlesHandler, trackHandler),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// This context is cancelled automatically on SIGINT / SIGTERM
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start the background cache warmer (runs immediately, then every 10 minutes)
	startCacheWarmer(ctx, orchestrator, 10*time.Minute)

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
