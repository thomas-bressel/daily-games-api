package pkg

import (
	"log/slog"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/tbressel/daily-games-api/internal/metrics"
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

// statusRecorder wraps ResponseWriter to capture the HTTP status code.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (sr *statusRecorder) WriteHeader(code int) {
	sr.status = code
	sr.ResponseWriter.WriteHeader(code)
}

// LogMiddleware logs the HTTP method, path, and duration of each request,
// and records the latency in the Prometheus histogram.
func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		slog.Info("[REQ]", "method", r.Method, "path", r.URL.Path)
		next.ServeHTTP(rec, r)
		duration := time.Since(start)
		slog.Info("[DONE]", "method", r.Method, "path", r.URL.Path, "duration", duration)
		metrics.HTTPDuration.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(rec.status)).Observe(duration.Seconds())
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

// CORSMiddleware adds the necessary CORS headers to allow browser extensions
// (which have no fixed origin) to call the API. It also handles OPTIONS preflight requests.
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
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
