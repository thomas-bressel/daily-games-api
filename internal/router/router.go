package router

import (
	"encoding/json"
	"net/http"

	"github.com/tbressel/daily-games-api/internal/handler"
	"github.com/tbressel/daily-games-api/pkg"
)

// Create builds and returns the main HTTP handler with all routes
// and middlewares applied.
//
// Routes:
//
//	GET  /health          - API health check
//	GET  /api/articles    - paginated article feed with optional filters
//	GET  /api/feeds       - list of active RSS feed sources
//	POST /api/track       - increment share or bookmark counter for an article
//	GET  /metrics         - total request count since server start
func Create(articlesHandler *handler.ArticlesHandler, trackHandler *handler.TrackHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handler.GetHealth)
	mux.HandleFunc("GET /api/articles", articlesHandler.GetArticles)
	mux.HandleFunc("GET /api/feeds", handler.GetFeeds)
	mux.HandleFunc("GET /api/track/{articleId}", trackHandler.GetTrack)
	mux.HandleFunc("POST /api/track", trackHandler.PostTrack)
	mux.HandleFunc("GET /metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"total_requests": pkg.GlobalMetrics.Value(),
		})
	})

	return pkg.ApplyMiddlewares(mux,
		pkg.LogMiddleware,
		pkg.RecoverMiddleware,
		pkg.MetricsMiddleware,
	)
}
