package handler

import (
	"net/http"

	"github.com/tbressel/daily-games-api/pkg"
)

// HealthResponse is the payload returned by the health check endpoint.
type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

// GetHealth handles GET /health.
// It returns a simple status response to confirm the API is running.
func GetHealth(w http.ResponseWriter, r *http.Request) {
	pkg.WriteSuccess(w, HealthResponse{
		Status:  "ok",
		Version: "1.0.0",
	})
}
