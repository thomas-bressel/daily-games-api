package pkg

import (
	"encoding/json"
	"net/http"

	"github.com/tbressel/daily-games-api/internal/handler"
)

// CreateRouter builds and returns the main HTTP handler with all routes
// and middlewares applied.
//
// Routes:
//
//	GET /health          - API health check
//	GET /api/articles    - paginated article feed with optional filters
//	GET /api/feeds       - list of active RSS feed sources
//	GET /metrics         - total request count since server start
func CreateRouter(articlesHandler *handler.ArticlesHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handler.GetHealth)
	mux.HandleFunc("GET /api/articles", articlesHandler.GetArticles)
	mux.HandleFunc("GET /api/feeds", handler.GetFeeds)
	mux.HandleFunc("GET /metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"total_requests": GlobalMetrics.Value(),
		})
	})

	return ApplyMiddlewares(mux,
		CORSMiddleware,
		LogMiddleware,
		RecoverMiddleware,
		MetricsMiddleware,
	)
}
