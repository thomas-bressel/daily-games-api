package pkg

import (
	"log/slog"
	"net/http"
	"sync/atomic"
	"time"
)

// metricsCounter is a thread-safe request counter using atomic operations.
type metricsCounter struct {
	total uint64
}

// Inc atomically increments the counter by 1.
func (m *metricsCounter) Inc() {
	atomic.AddUint64(&m.total, 1)
}

// Value atomically reads and returns the current counter value.
func (m *metricsCounter) Value() uint64 {
	return atomic.LoadUint64(&m.total)
}

// GlobalMetrics is the shared request counter accessible from the metrics endpoint.
var GlobalMetrics metricsCounter

// LogMiddleware logs the HTTP method, path, and duration of each request.
func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		slog.Info("[REQ]", "method", r.Method, "path", r.URL.Path)
		next.ServeHTTP(w, r)
		slog.Info("[DONE]", "method", r.Method, "path", r.URL.Path, "duration", time.Since(start))
	})
}

// RecoverMiddleware catches any panic in downstream handlers,
// logs the error, and returns a 500 Internal Server Error response.
func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				slog.Error("[PANIC]", "err", rec)
				http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// MetricsMiddleware increments the global request counter for each incoming request.
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		GlobalMetrics.Inc()
		next.ServeHTTP(w, r)
	})
}

//ApplyMiddlewares wraps a handler with a chain of middlewares.
// Middlewares are applied in reverse order so the first one in the list
// is the outermost (executed first).
func ApplyMiddlewares(h http.Handler, mws ...func(http.Handler) http.Handler) http.Handler {
	for _, mw := range mws {
		h = mw(h)
	}
	return h
}
