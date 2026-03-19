package handler

import (
	"net/http"

	"github.com/tbressel/daily-games-api/internal/feed"
	"github.com/tbressel/daily-games-api/pkg"
)

// GetFeeds handles GET /api/feeds.
// It returns all active RSS feed sources with their metadata.
func GetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds := feed.GetActive()
	pkg.WriteSuccess(w, feeds)
}
