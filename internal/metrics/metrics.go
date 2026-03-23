package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	// HTTPDuration tracks request latency per method and path (histogram).
	HTTPDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds.",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10, 20},
		},
		[]string{"method", "path", "status"},
	)

	// CacheHits counts Redis cache hits per category and lang.
	CacheHits = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total number of Redis cache hits.",
		},
		[]string{"category", "lang"},
	)

	// CacheMisses counts Redis cache misses per category and lang.
	CacheMisses = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total number of Redis cache misses.",
		},
		[]string{"category", "lang"},
	)

	// WarmerDuration tracks how long each full cache warm-up cycle takes.
	WarmerDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "cache_warmer_duration_seconds",
			Help:    "Duration of a full cache warm-up cycle in seconds.",
			Buckets: []float64{1, 5, 10, 20, 30, 60, 120},
		},
	)

	// WarmerErrors counts errors encountered during cache warm-up.
	WarmerErrors = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "cache_warmer_errors_total",
			Help: "Total number of errors during cache warm-up cycles.",
		},
	)
)

func init() {
	prometheus.MustRegister(
		HTTPDuration,
		CacheHits,
		CacheMisses,
		WarmerDuration,
		WarmerErrors,
	)
}
