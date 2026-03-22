package handler

import (
	"net/http"
	"strconv"

	"github.com/tbressel/daily-games-api/internal/article"
	"github.com/tbressel/daily-games-api/pkg"
)

// ArticlesHandler holds the dependencies needed to handle article requests.
type ArticlesHandler struct {
	orchestrator *article.Orchestrator
}

// NewArticlesHandler creates a new ArticlesHandler with the given orchestrator.
func NewArticlesHandler(orchestrator *article.Orchestrator) *ArticlesHandler {
	return &ArticlesHandler{orchestrator: orchestrator}
}

// GetArticles handles GET /api/articles.
// It parses optional query parameters for filtering and pagination,
// then delegates to the orchestrator to fetch, cache, and paginate articles.
//
// Query parameters:
//
//	offset   int    - number of articles to skip (default: 0)
//	limit    int    - max articles to return (default: 20, max: 100)
//	source   string - filter by feed source ID (e.g. "amstrad-eu")
//	category string - filter by category (e.g. "retrogaming")
//	refresh  bool   - force bypass of Redis cache (default: false)
func (h *ArticlesHandler) GetArticles(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	offset := parseIntParam(q.Get("offset"), 0)
	limit := parseIntParam(q.Get("limit"), 20)
	if limit > 200 {
		limit = 200
	}

	filters := pkg.ArticleFilters{
		Offset:   offset,
		Limit:    limit,
		Source:   q.Get("source"),
		Category: q.Get("category"),
		Lang:     q.Get("lang"),
		Refresh:  q.Get("refresh") == "true",
	}

	data, err := h.orchestrator.GetArticles(r.Context(), filters)
	if err != nil {
		pkg.WriteError(w, http.StatusInternalServerError, "Failed to fetch articles")
		return
	}

	pkg.WriteSuccess(w, data)
}

// parseIntParam parses a query string value as an integer.
// Returns the fallback value if the string is empty or cannot be parsed.
func parseIntParam(val string, fallback int) int {
	if val == "" {
		return fallback
	}
	n, err := strconv.Atoi(val)
	if err != nil || n < 0 {
		return fallback
	}
	return n
}
