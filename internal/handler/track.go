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

// PostTrack handles POST /api/track.
// Increments the Redis counter for the given article and event type.
// Event must be "share" or "bookmark" — any other value returns 400.
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
