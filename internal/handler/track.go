package handler

import (
	"log/slog"
	"net/http"

	"github.com/tbressel/daily-games-api/internal/cache"
	"github.com/tbressel/daily-games-api/pkg"
)

// TrackHandler handles article interaction tracking requests.
type TrackHandler struct {
	cache *cache.Client
}

// NewTrackHandler creates a new TrackHandler with the given Redis cache client.
func NewTrackHandler(cache *cache.Client) *TrackHandler {
	return &TrackHandler{cache: cache}
}

// GetTrack handles GET /api/track/{articleId}.
// Returns the bookmark and share counters for a given article.
func (h *TrackHandler) GetTrack(w http.ResponseWriter, r *http.Request) {
	articleID := r.PathValue("articleId")
	if articleID == "" {
		pkg.WriteError(w, http.StatusBadRequest, "articleId is required")
		return
	}

	bookmarks, err := h.cache.GetTrack(r.Context(), articleID, "bookmark")
	if err != nil {
		slog.Error("[Track] Redis error", "err", err)
		pkg.WriteError(w, http.StatusInternalServerError, "tracking failed")
		return
	}

	shares, err := h.cache.GetTrack(r.Context(), articleID, "share")
	if err != nil {
		slog.Error("[Track] Redis error", "err", err)
		pkg.WriteError(w, http.StatusInternalServerError, "tracking failed")
		return
	}

	pkg.WriteSuccess(w, map[string]any{
		"articleId": articleID,
		"bookmarks": bookmarks,
		"shares":    shares,
	})
}

// PostTrackBatch handles POST /api/track/batch.
// Returns bookmark and share counters for multiple articles in a single request.
// Expects body: { "ids": ["id1", "id2", ...] }
// Returns: { "status": "ok", "data": { "id1": { "bookmarks": N, "shares": N }, ... } }
func (h *TrackHandler) PostTrackBatch(w http.ResponseWriter, r *http.Request) {
	var body struct {
		IDs []string `json:"ids"`
	}
	if err := pkg.ParseJSON(w, r, &body); err != nil {
		pkg.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if len(body.IDs) == 0 {
		pkg.WriteSuccess(w, map[string]any{})
		return
	}

	counts, err := h.cache.GetBatchTrack(r.Context(), body.IDs)
	if err != nil {
		slog.Error("[Track] Redis batch error", "err", err)
		pkg.WriteError(w, http.StatusInternalServerError, "tracking failed")
		return
	}

	result := make(map[string]any, len(counts))
	for id, c := range counts {
		result[id] = map[string]int64{
			"bookmarks": c["bookmark"],
			"shares":    c["share"],
		}
	}

	pkg.WriteSuccess(w, result)
}

// PostTrack handles POST /api/track.
// Increments the Redis counter for the given article and event type.
// Event must be "share" or "bookmark"  -- any other value returns 400.
func (h *TrackHandler) PostTrack(w http.ResponseWriter, r *http.Request) {
	var body pkg.TrackEvent
	if err := pkg.ParseJSON(w, r, &body); err != nil {
		pkg.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if body.ArticleID == "" {
		pkg.WriteError(w, http.StatusBadRequest, "articleId is required")
		return
	}

	if body.Event != "share" && body.Event != "bookmark" {
		pkg.WriteError(w, http.StatusBadRequest, "event must be 'share' or 'bookmark'")
		return
	}

	count, err := h.cache.IncrTrack(r.Context(), body.ArticleID, body.Event)
	if err != nil {
		slog.Error("[Track] Redis error", "err", err)
		pkg.WriteError(w, http.StatusInternalServerError, "tracking failed")
		return
	}

	pkg.WriteSuccess(w, map[string]any{
		"articleId": body.ArticleID,
		"event":     body.Event,
		"count":     count,
	})
}
